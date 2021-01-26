package config

import (
	"time"

	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

func Config() Provider {
	return defaultConfig
}

func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("VOYAGER")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults
	
	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	v.SetDefault("postgres_dsn", "postgres://obsurvive:obsurvive@localhost/obsurvive?sslmode=disable")
	v.SetDefault("scheduler_table_name", "scheduler")
	v.SetDefault("randomize_factor", 0)

	v.SetDefault("bind_address", ":9000")

	v.SetDefault("amqp_dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("amqp_exchange_name", "voyager.probe")
	v.SetDefault("amqp_queue_name", "probes")
	v.SetDefault("amqp_key_name", "probes")

	v.SetDefault("redis_host", "localhost:6379")
	v.SetDefault("redis_password", "")
	v.SetDefault("redis_db", 0)
	v.SetDefault("redis_expire_seconds", 300)


	return v
}
