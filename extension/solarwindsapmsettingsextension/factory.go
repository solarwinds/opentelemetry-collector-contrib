package solarwindsapmsettingsextension

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"

	"github.com/solarwinds/opentelemetry-collector-contrib/extension/solarwindsapmsettingsextension/internal/metadata"
)

const (
	DefaultInterval = "10s"
)

func createDefaultConfig() component.Config {
	return &Config{
		Interval: DefaultInterval,
	}
}

func createExtension(_ context.Context, settings extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return newSolarwindsApmSettingsExtension(cfg.(*Config), settings.Logger)
}

func NewFactory() extension.Factory {
	return extension.NewFactory(
		metadata.Type,
		createDefaultConfig,
		createExtension,
		metadata.ExtensionStability,
	)
}
