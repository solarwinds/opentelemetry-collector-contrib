package solarwindsapmsettingsextension

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Key      string `mapstructure:"key"`
	Interval string `mapstructure:"interval"`
}

func (cfg *Config) Validate() error {
	/**
	 * No configuration validation here after discussion.
	 * Extension will check the value and print error/warning message instead
	 */
	return nil
}
