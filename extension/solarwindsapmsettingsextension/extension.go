package solarwindsapmsettingsextension

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/solarwindscloud/apm-proto/go/collectorpb"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RawOutputFile   = "/tmp/solarwinds-apm-settings-raw"
	JSONOutputFile  = "/tmp/solarwinds-apm-settings.json"
	MinimumInterval = "5s"
	MaximumInterval = "60s"
)

type solarwindsapmSettingsExtension struct {
	logger *zap.Logger
	config *Config
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	client collectorpb.TraceCollectorClient
}

func newSolarwindsApmSettingsExtension(extensionCfg *Config, logger *zap.Logger) (extension.Extension, error) {
	settingsExtension := &solarwindsapmSettingsExtension{
		config: extensionCfg,
		logger: logger,
	}
	return settingsExtension, nil
}

func validateSolarwindsApmSettingsExtensionConfiguration(extensionCfg *Config, logger *zap.Logger) bool {
	// Endpoint
	if len(extensionCfg.Endpoint) == 0 {
		logger.Error("endpoint must not be empty")
		return false
	}
	endpointArr := strings.Split(extensionCfg.Endpoint, ":")
	if len(endpointArr) != 2 {
		logger.Error("endpoint should be in \"<host>:<port>\" format")
		return false
	}
	if len(endpointArr[0]) == 0 {
		logger.Error("endpoint should be in \"<host>:<port>\" format and \"<host>\" must not be empty")
		return false
	}
	if len(endpointArr[1]) == 0 {
		logger.Error("endpoint should be in \"<host>:<port>\" format and \"<port>\" must not be empty")
		return false
	}
	matched, _ := regexp.MatchString(`apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\-]*.solarwinds.com`, endpointArr[0])
	if !matched {
		logger.Error("endpoint \"<host>\" part should be in \"apm.collector.[a-z]{2,3}-[0-9]{2}.[a-z\\-]*.solarwinds.com\" regex format, see https://documentation.solarwinds.com/en/success_center/observability/content/system_requirements/endpoints.htm for detail")
		return false
	}
	if _, err := strconv.Atoi(endpointArr[1]); err != nil {
		logger.Error("the <port> portion of endpoint has to be an integer")
		return false
	}
	// Key
	if len(extensionCfg.Key) == 0 {
		logger.Error("key must not be empty")
		return false
	}
	keyArr := strings.Split(extensionCfg.Key, ":")
	if len(keyArr) != 2 {
		logger.Error("key should be in \"<token>:<service_name>\" format")
		return false
	}
	if len(keyArr[0]) == 0 {
		logger.Error("key should be in \"<token>:<service_name>\" format and \"<token>\" must not be empty")
		return false
	}
	if len(keyArr[1]) == 0 {
		logger.Error("key should be in \"<token>:<service_name>\" format and \"<service_name>\" must not be empty")
		return false
	}
	/*
	 * Interval
	 * We don't return false here as we always has an interval value
	 */
	if interval, err := time.ParseDuration(extensionCfg.Interval); err != nil {
		logger.Warn("interval has to be a duration string. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\". use default " + DefaultInterval + " instead")
		extensionCfg.Interval = DefaultInterval
	} else {
		minimum, _ := time.ParseDuration(MinimumInterval)
		maximum, _ := time.ParseDuration(MaximumInterval)
		if interval.Seconds() < minimum.Seconds() {
			logger.Warn("Interval " + extensionCfg.Interval + " is less than the minimum supported interval " + MinimumInterval + ". use minimum interval " + MinimumInterval + " instead")
			extensionCfg.Interval = MinimumInterval
		}
		if interval.Seconds() > maximum.Seconds() {
			logger.Warn("Interval " + extensionCfg.Interval + " is greater than the maximum supported interval " + MaximumInterval + ". use maximum interval " + MaximumInterval + " instead")
			extensionCfg.Interval = MaximumInterval
		}
	}
	return true
}

func refresh(extension *solarwindsapmSettingsExtension) {
	extension.logger.Info("Time to refresh from " + extension.config.Endpoint)
	if hostname, err := os.Hostname(); err != nil {
		extension.logger.Error("Unable to call os.Hostname() " + err.Error())
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		request := &collectorpb.SettingsRequest{
			ApiKey: extension.config.Key,
			Identity: &collectorpb.HostID{
				Hostname: hostname,
			},
			ClientVersion: "2",
		}
		if response, err := extension.client.GetSettings(ctx, request); err != nil {
			extension.logger.Error("Unable to getSettings from " + extension.config.Endpoint + " " + err.Error())
		} else {
			switch result := response.GetResult(); result {
			case collectorpb.ResultCode_OK:
				if len(response.GetWarning()) != 0 {
					extension.logger.Error(response.GetWarning())
				}
				if bytes, err := proto.Marshal(response); err != nil {
					extension.logger.Error("Unable to marshal response to bytes " + err.Error())
				} else {
					// Output in raw format
					if err := os.WriteFile(RawOutputFile, bytes, 0644); err != nil {
						extension.logger.Error("Unable to write " + RawOutputFile + " " + err.Error())
					} else {
						extension.logger.Info(RawOutputFile + " is refreshed")
					}
				}
				// Output in human-readable format
				var settings []map[string]interface{}
				for _, item := range response.GetSettings() {

					marshalOptions := protojson.MarshalOptions{
						UseEnumNumbers:  true,
						EmitUnpopulated: true,
					}
					if settingBytes, err := marshalOptions.Marshal(item); err != nil {
						extension.logger.Warn("Error to marshal setting JSON[] byte from response.GetSettings() " + err.Error())
					} else {
						setting := make(map[string]interface{})
						if err := json.Unmarshal(settingBytes, &setting); err != nil {
							extension.logger.Warn("Error to unmarshal setting JSON object from setting JSON[]byte " + err.Error())
						} else {
							if value, ok := setting["value"].(string); ok {
								if num, e := strconv.ParseInt(value, 10, 0); e != nil {
									extension.logger.Warn("Unable to parse value " + value + " as number " + e.Error())
								} else {
									setting["value"] = num
								}
							}
							if timestamp, ok := setting["timestamp"].(string); ok {
								if num, e := strconv.ParseInt(timestamp, 10, 0); e != nil {
									extension.logger.Warn("Unable to parse timestamp " + timestamp + " as number " + e.Error())
								} else {
									setting["timestamp"] = num
								}
							}
							if ttl, ok := setting["ttl"].(string); ok {
								if num, e := strconv.ParseInt(ttl, 10, 0); e != nil {
									extension.logger.Warn("Unable to parse ttl " + ttl + " as number " + e.Error())
								} else {
									setting["ttl"] = num
								}
							}
							if _, ok := setting["flags"]; ok {
								setting["flags"] = string(item.Flags)
							}
							if arguments, ok := setting["arguments"].(map[string]interface{}); ok {
								if value, ok := item.Arguments["BucketCapacity"]; ok {
									arguments["BucketCapacity"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["BucketRate"]; ok {
									arguments["BucketRate"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["TriggerRelaxedBucketCapacity"]; ok {
									arguments["TriggerRelaxedBucketCapacity"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["TriggerRelaxedBucketRate"]; ok {
									arguments["TriggerRelaxedBucketRate"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["TriggerStrictBucketCapacity"]; ok {
									arguments["TriggerStrictBucketCapacity"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["TriggerStrictBucketRate"]; ok {
									arguments["TriggerStrictBucketRate"] = math.Float64frombits(binary.LittleEndian.Uint64(value))
								}
								if value, ok := item.Arguments["MetricsFlushInterval"]; ok {
									arguments["MetricsFlushInterval"] = int32(binary.LittleEndian.Uint32(value))
								}
								if value, ok := item.Arguments["MaxTransactions"]; ok {
									arguments["MaxTransactions"] = int32(binary.LittleEndian.Uint32(value))
								}
								if value, ok := item.Arguments["MaxCustomMetrics"]; ok {
									arguments["MaxCustomMetrics"] = int32(binary.LittleEndian.Uint32(value))
								}
								if value, ok := item.Arguments["EventsFlushInterval"]; ok {
									arguments["EventsFlushInterval"] = int32(binary.LittleEndian.Uint32(value))
								}
								if value, ok := item.Arguments["ProfilingInterval"]; ok {
									arguments["ProfilingInterval"] = int32(binary.LittleEndian.Uint32(value))
								}
								// Remove SignatureKey from collector response
								delete(arguments, "SignatureKey")
							}
							settings = append(settings, setting)
						}
					}
				}
				if content, err := json.Marshal(settings); err != nil {
					extension.logger.Warn("Error to marshal setting JSON[] byte from settings " + err.Error())
				} else {
					if err := os.WriteFile(JSONOutputFile, content, 0644); err != nil {
						extension.logger.Error("Unable to write " + JSONOutputFile + " " + err.Error())
					} else {
						extension.logger.Info(JSONOutputFile + " is refreshed")
						extension.logger.Info(string(content))
					}
				}
			case collectorpb.ResultCode_TRY_LATER:
				extension.logger.Warn("GetSettings returned TRY_LATER " + response.GetWarning())
			case collectorpb.ResultCode_INVALID_API_KEY:
				extension.logger.Warn("GetSettings returned INVALID_API_KEY " + response.GetWarning())
			case collectorpb.ResultCode_LIMIT_EXCEEDED:
				extension.logger.Warn("GetSettings returned LIMIT_EXCEEDED " + response.GetWarning())
			case collectorpb.ResultCode_REDIRECT:
				extension.logger.Warn("GetSettings returned REDIRECT " + response.GetWarning())
			default:
				extension.logger.Warn("Unknown ResultCode from GetSettings " + response.GetWarning())
			}
		}
	}
}

func (extension *solarwindsapmSettingsExtension) Start(ctx context.Context, _ component.Host) error {
	extension.logger.Debug("Starting up solarwinds apm settings extension")
	ctx = context.Background()
	ctx, extension.cancel = context.WithCancel(ctx)

	var err error
	extension.conn, err = grpc.Dial(extension.config.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return errors.New("Failed to dial: " + err.Error())
	} else {
		extension.logger.Info("Dailed to " + extension.config.Endpoint)
	}
	extension.client = collectorpb.NewTraceCollectorClient(extension.conn)

	// Refresh immediately if configuration passes validation
	if validateSolarwindsApmSettingsExtensionConfiguration(extension.config, extension.logger) {
		refresh(extension)
	} else {
		extension.logger.Warn("No refresh due to invalid config value")
	}

	// setup lightweight thread to refresh
	interval, _ := time.ParseDuration(extension.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// Refresh at each ticker event
				if validateSolarwindsApmSettingsExtensionConfiguration(extension.config, extension.logger) {
					refresh(extension)
				} else {
					extension.logger.Warn("No refresh due to invalid config value")
				}
			case <-ctx.Done():
				extension.logger.Info("Received ctx.Done() from ticker")
				return
			}
		}
	}()

	return nil
}

func (extension *solarwindsapmSettingsExtension) Shutdown(_ context.Context) error {
	extension.logger.Debug("Shutting down solarwinds apm settings extension")
	err := extension.conn.Close()
	if err != nil {
		return errors.New("Failed to close the gRPC connection to solarwinds APM collector " + err.Error())
	}
	return nil
}
