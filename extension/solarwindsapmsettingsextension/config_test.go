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
			name: "empty key",
			cfg: &Config{
				Endpoint: "host:12345",
			},
			err: errors.New("key must not be empty"),
		},
		{
			name: "invalid endpoint",
			cfg: &Config{
				Endpoint: "invalid",
				Key:      "token:name",
			},
			err: errors.New("endpoint should be in \"<host>:<port>\" format"),
		},
		{
			name: "invalid endpoint format but port is not an integer",
			cfg: &Config{
				Endpoint: "host:abc",
				Key:      "token:name",
			},
			err: errors.New("the <port> portion of endpoint has to be an integer"),
		},
		{
			name: "invalid key",
			cfg: &Config{
				Endpoint: "host:12345",
				Key:      "invalid",
			},
			err: errors.New("key should be in \"<token>:<service_name>\" format"),
		},
		{
			name: "valid",
			cfg: &Config{
				Endpoint: "host:12345",
				Key:      "token:name",
				Interval: "1s",
			},
			err: nil,
		},
		{
			name: "empty_interval",
			cfg: &Config{
				Endpoint: "host:12345",
				Key:      "token:name",
				Interval: "",
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
