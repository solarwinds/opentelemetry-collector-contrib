package solarwindsapmsettingsextension

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestCreateExtension(t *testing.T) {
	conf := &Config{
		Endpoint: "apm-testcollector.click:443",
		Key:      "valid:unittest",
		Interval: "10s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionWrongEndpoint(t *testing.T) {
	conf := &Config{
		Endpoint: "apm-testcollector.nothing:443",
		Key:      "valid:unittest",
		Interval: "5s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionUnAuthorizedKeyToAPMCollector(t *testing.T) {
	conf := &Config{
		Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
		Key:      "invalid",
		Interval: "60s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionUnAuthorizedKeyWithServiceNameToAPMCollector(t *testing.T) {
	conf := &Config{
		Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
		Key:      "invalid:service_name",
		Interval: "60s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionEmptyKeyWithServiceNameToAPMCollector(t *testing.T) {
	conf := &Config{
		Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
		Key:      ":service_name",
		Interval: "60s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionNoSuchHost(t *testing.T) {
	conf := &Config{
		Endpoint: "apm.collector.na-99.cloud.solarwinds.com:443",
		Key:      "invalid",
		Interval: "60s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionWrongKey(t *testing.T) {
	conf := &Config{
		Endpoint: "apm-testcollector.click:443",
		Key:      "invalid",
		Interval: "60s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionIntervalLessThanMinimum(t *testing.T) {
	conf := &Config{
		Endpoint: "apm-testcollector.click:443",
		Key:      "valid:unittest",
		Interval: "4s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

func TestCreateExtensionIntervalGreaterThanMaximum(t *testing.T) {
	conf := &Config{
		Endpoint: "apm-testcollector.click:443",
		Key:      "valid:unittest",
		Interval: "61s",
	}
	ex := createAnExtension(conf, t)
	ex.Shutdown(context.TODO())
}

// create extension
func createAnExtension(c *Config, t *testing.T) extension.Extension {
	logger, err := zap.NewProduction()
	ex, err := newSolarwindsApmSettingsExtension(c, logger)
	require.NoError(t, err)
	err = ex.Start(context.TODO(), nil)
	require.NoError(t, err)
	return ex
}

func TestValidateSolarwindsApmSettingsExtensionConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		ok   bool

		message string
	}{
		{
			name:    "nothing",
			cfg:     &Config{},
			ok:      false,
			message: "endpoint must not be empty",
		},
		{
			name: "valid configuration",
			cfg: &Config{
				Endpoint: "apm.collector.na-02.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "10s",
			},
			ok:      true,
			message: "",
		},
		{
			name: "endpoint without :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com",
			},
			ok:      false,
			message: "endpoint should be in \"<host>:<port>\" format",
		},
		{
			name: "endpoint with some :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:a:b",
			},
			ok:      false,
			message: "endpoint should be in \"<host>:<port>\" format",
		},
		{
			name: "endpoint with invalid port",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:port",
			},
			ok:      false,
			message: "the <port> portion of endpoint has to be an integer",
		},
		{
			name: "bad endpoint",
			cfg: &Config{
				Endpoint: "apm.collector..cloud.solarwinds.com:443",
			},
			ok:      false,
			message: "endpoint \"<host>\" part should be in \"apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\\-]*.solarwinds.com\" regex format, see https://documentation.solarwinds.com/en/success_center/observability/content/system_requirements/endpoints.htm for detail",
		},
		{
			name: "empty endpoint with port",
			cfg: &Config{
				Endpoint: ":433",
			},
			ok:      false,
			message: "endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty",
		},
		{
			name: "empty endpoint without port",
			cfg: &Config{
				Endpoint: ":",
			},
			ok:      false,
			message: "endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty",
		},
		{
			name: "endpoint without port",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:",
			},
			ok:      false,
			message: "endpoint should be in \"<host>:<port>\" format and \"<port>\" must not be empty",
		},
		{
			name: "valid endpoint but empty key",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
			},
			ok:      false,
			message: "key must not be empty",
		},
		{
			name: "key is :",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      ":",
			},
			ok:      false,
			message: "key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty",
		},
		{
			name: "key is ::",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "::",
			},
			ok:      false,
			message: "key should be in \"<token>:<service_name>\" format",
		},
		{
			name: "key is :name",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      ":name",
			},
			ok:      false,
			message: "key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty",
		},
		{
			name: "key is token:",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:",
			},
			ok:      false,
			message: "key should be in \"<token>:<service_name>\" format and \"<service_name>\" must not be empty",
		},
		{
			name: "empty_interval",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "",
			},
			ok:      true,
			message: "interval has to be a duration string. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\". use default " + DefaultInterval + " instead",
		},
		{
			name: "interval is not a duration string",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "something",
			},
			ok:      true,
			message: "interval has to be a duration string. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\". use default " + DefaultInterval + " instead",
		},
		{
			name: "minimum_interval",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "4s",
			},
			ok:      true,
			message: "Interval 4s is less than the minimum supported interval " + MinimumInterval + ". use minimum interval " + MinimumInterval + " instead",
		},
		{
			name: "maximum_interval",
			cfg: &Config{
				Endpoint: "apm.collector.na-01.cloud.solarwinds.com:443",
				Key:      "token:name",
				Interval: "61s",
			},
			ok:      true,
			message: "Interval 61s is greater than the maximum supported interval " + MaximumInterval + ". use maximum interval " + MaximumInterval + " instead",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			observedLogger := zap.New(observedZapCore)
			require.Equal(t, tc.ok, validateSolarwindsApmSettingsExtensionConfiguration(tc.cfg, observedLogger))
			if len(tc.message) != 0 {
				require.Equal(t, 1, observedLogs.Len())
				require.Equal(t, tc.message, observedLogs.All()[0].Message)
			}
		})
	}
}
