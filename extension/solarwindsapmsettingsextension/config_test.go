package solarwindsapmsettingsextension

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		err  error
	}{
		{
			name: "nothing",
			cfg:  &Config{},
			err:  errors.New("endpoint must not be empty"),
		},
		{
			name: "valid configuration",
			cfg: &Config{
				Endpoint: "apm.collector.na-02.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "10s",
			},
			err: nil,
		},
		{
			name: "endpoint without :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format"),
		},
		{
			name: "endpoint with some :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:a:b",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format"),
		},
		{
			name: "endpoint with invalid point",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:port",
			},
			err: errors.New("the <port> portion of endpoint has to be an integer"),
		},
		{
			name: "bad endpoint",
			cfg: &Config{
				Endpoint: "apm.collector..cloud.solarwinds.com:443",
			},
			err: errors.New("endpoint \"<host>\" part should be in \"apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\\-]*.solarwinds.com\" regex format, see https://documentation.solarwinds.com/en/success_center/observability/content/system_requirements/endpoints.htm for detail"),
		},
		{
			name: "empty endpoint with port",
			cfg: &Config{
				Endpoint: ":433",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty"),
		},
		{
			name: "empty endpoint without port",
			cfg: &Config{
				Endpoint: ":",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty"),
		},
		{
			name: "endpoint without port",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format and \"<port>\" must not be empty"),
		},
		{
			name: "valid endpoint but empty key",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
			},
			err: errors.New("key must not be empty"),
		},
		{
			name: "key is :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      ":",
			},
			err: errors.New("key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty"),
		},
		{
			name: "key is ::",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "::",
			},
			err: errors.New("key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty"),
		},
		{
			name: "key is :name",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      ":name",
			},
			err: errors.New("key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty"),
		},
		{
			name: "key is token:",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:",
			},
			err: errors.New("key should be in \"<token>:<service_name>\" format and \"<service_name>\" must not be empty"),
		},
		{
			name: "empty_interval",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "",
			},
			err: errors.New("interval must not be empty"),
		},
		{
			name: "interval is not a duration string",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "something",
			},
			err: errors.New("interval has to be a duration string. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\""),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
