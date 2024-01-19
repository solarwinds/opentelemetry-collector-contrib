package solarwindsapmsettingsextension

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Key      string `mapstructure:"key"`
	Interval string `mapstructure:"interval"`
}

func (cfg *Config) Validate() error {
	// Endpoint
	if len(cfg.Endpoint) == 0 {
		return errors.New("endpoint must not be empty")
	}
	endpointArr := strings.Split(cfg.Endpoint, ":")
	if len(endpointArr) != 2 {
		return errors.New("endpoint should be in \"<host>:<port>\" format")
	}
	if len(endpointArr[0]) == 0 {
		return errors.New("endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty")
	}
	if len(endpointArr[1]) == 0 {
		return errors.New("endpoint should be in \"<host>:<port>\" format and \"<port>\" must not be empty")
	}
	matched, _ := regexp.MatchString(`apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\-]*.solarwinds.com`, endpointArr[0])
	if !matched {
		return errors.New("endpoint \"<host>\" part should be in \"apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\\-]*.solarwinds.com\" regex format, see https://documentation.solarwinds.com/en/success_center/observability/content/system_requirements/endpoints.htm for detail")
	}
	if _, err := strconv.Atoi(endpointArr[1]); err != nil {
		return errors.New("the <port> portion of endpoint has to be an integer")
	}
	// Key
	if len(cfg.Key) == 0 {
		return errors.New("key must not be empty")
	}
	keyArr := strings.Split(cfg.Key, ":")
	if len(keyArr) != 2 {
		return errors.New("key should be in \"<token>:<service_name>\" format")
	}
	if len(keyArr[0]) == 0 {
		return errors.New("key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty")
	}
	if len(keyArr[1]) == 0 {
		return errors.New("key should be in \"<token>:<service_name>\" format and \"<service_name>\" must not be empty")
	}
	// Interval
	if len(cfg.Interval) == 0 {
		return errors.New("interval must not be empty")
	}
	if _, err := time.ParseDuration(cfg.Interval); err != nil {
		return errors.New("interval has to be a duration string. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\"")
	}
	return nil
}
