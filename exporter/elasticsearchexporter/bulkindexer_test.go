// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearchexporter

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var defaultRoundTripFunc = func(*http.Request) (*http.Response, error) {
	return &http.Response{
		Body: io.NopCloser(strings.NewReader("{}")),
	}, nil
}

type mockTransport struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.RoundTripFunc == nil {
		return defaultRoundTripFunc(req)
	}
	return t.RoundTripFunc(req)
}

const successResp = `{
  "took": 30,
  "errors": false,
  "items": [
    {
      "create": {
        "_index": "foo",
        "status": 201
      }
    }
  ]
}`

func TestAsyncBulkIndexer_flushOnClose(t *testing.T) {
	cfg := Config{NumWorkers: 1, Flush: FlushSettings{Interval: time.Hour, Bytes: 2 << 30}}
	client, err := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		RoundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
				Body:   io.NopCloser(strings.NewReader(successResp)),
			}, nil
		},
	}})
	require.NoError(t, err)

	bulkIndexer, err := newAsyncBulkIndexer(zap.NewNop(), client, &cfg)
	require.NoError(t, err)
	session, err := bulkIndexer.StartSession(context.Background())
	require.NoError(t, err)

	assert.NoError(t, session.Add(context.Background(), "foo", strings.NewReader(`{"foo": "bar"}`), nil))
	assert.NoError(t, bulkIndexer.Close(context.Background()))
	assert.Equal(t, int64(1), bulkIndexer.stats.docsIndexed.Load())
}

func TestAsyncBulkIndexer_flush(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "flush.bytes",
			config: Config{NumWorkers: 1, Flush: FlushSettings{Interval: time.Hour, Bytes: 1}},
		},
		{
			name:   "flush.interval",
			config: Config{NumWorkers: 1, Flush: FlushSettings{Interval: 50 * time.Millisecond, Bytes: 2 << 30}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, err := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
						Body:   io.NopCloser(strings.NewReader(successResp)),
					}, nil
				},
			}})
			require.NoError(t, err)

			bulkIndexer, err := newAsyncBulkIndexer(zap.NewNop(), client, &tt.config)
			require.NoError(t, err)
			session, err := bulkIndexer.StartSession(context.Background())
			require.NoError(t, err)

			assert.NoError(t, session.Add(context.Background(), "foo", strings.NewReader(`{"foo": "bar"}`), nil))
			// should flush
			time.Sleep(100 * time.Millisecond)
			assert.Equal(t, int64(1), bulkIndexer.stats.docsIndexed.Load())
			assert.NoError(t, bulkIndexer.Close(context.Background()))
		})
	}
}

func TestAsyncBulkIndexer_flush_error(t *testing.T) {
	tests := []struct {
		name          string
		roundTripFunc func(*http.Request) (*http.Response, error)
		wantMessage   string
		wantFields    []zap.Field
	}{
		{
			name: "500",
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
					Body:       io.NopCloser(strings.NewReader("error")),
				}, nil
			},
			wantMessage: "bulk indexer flush error",
		},
		{
			name: "429",
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 429,
					Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
					Body:       io.NopCloser(strings.NewReader("error")),
				}, nil
			},
			wantMessage: "bulk indexer flush error",
		},
		{
			name: "transport error",
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return nil, errors.New("transport error")
			},
			wantMessage: "bulk indexer flush error",
		},
		{
			name: "known version conflict error",
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
					Body: io.NopCloser(strings.NewReader(
						`{"items":[{"create":{"_index":".ds-metrics-generic.otel-default","status":400,"error":{"type":"version_conflict_engine_exception","reason":""}}}]}`)),
				}, nil
			},
			wantMessage: "failed to index document",
			wantFields:  []zap.Field{zap.String("hint", "check the \"Known issues\" section of Elasticsearch Exporter docs")},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := Config{NumWorkers: 1, Flush: FlushSettings{Interval: time.Hour, Bytes: 1}}
			client, err := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
				RoundTripFunc: tt.roundTripFunc,
			}})
			require.NoError(t, err)
			core, observed := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))

			bulkIndexer, err := newAsyncBulkIndexer(zap.New(core), client, &cfg)
			require.NoError(t, err)
			session, err := bulkIndexer.StartSession(context.Background())
			require.NoError(t, err)

			assert.NoError(t, session.Add(context.Background(), "foo", strings.NewReader(`{"foo": "bar"}`), nil))
			// should flush
			time.Sleep(100 * time.Millisecond)
			assert.Equal(t, int64(0), bulkIndexer.stats.docsIndexed.Load())
			assert.NoError(t, bulkIndexer.Close(context.Background()))
			messages := observed.FilterMessage(tt.wantMessage)
			require.Equal(t, 1, messages.Len(), "message not found; observed.All()=%v", observed.All())
			for _, wantField := range tt.wantFields {
				assert.Equal(t, 1, messages.FilterField(wantField).Len(), "message with field not found; observed.All()=%v", observed.All())
			}
		})
	}
}
