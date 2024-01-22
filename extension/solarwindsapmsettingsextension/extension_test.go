package solarwindsapmsettingsextension

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
	"testing"
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
