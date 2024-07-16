module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/cassandraexporter

go 1.21.0

require (
	github.com/gocql/gocql v1.6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.104.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/component v0.104.1-0.20240716040034-2d7dea6a5a2e
	go.opentelemetry.io/collector/config/configopaque v1.11.1-0.20240716040034-2d7dea6a5a2e
	go.opentelemetry.io/collector/confmap v0.104.1-0.20240716040034-2d7dea6a5a2e
	go.opentelemetry.io/collector/exporter v0.104.1-0.20240716040034-2d7dea6a5a2e
	go.opentelemetry.io/collector/pdata v1.11.1-0.20240716040034-2d7dea6a5a2e
	go.opentelemetry.io/otel/metric v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.opentelemetry.io/collector v0.104.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/config/configretry v1.11.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.104.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/consumer v0.104.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/extension v0.104.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/featuregate v1.11.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/collector/internal/globalgates v0.0.0-20240715203212-534768cf7a80 // indirect
	go.opentelemetry.io/collector/receiver v0.104.1-0.20240716040034-2d7dea6a5a2e // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.50.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden => ../../pkg/golden
