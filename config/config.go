package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// The idea behind this, is to load and use the config only one time, only inside the main.go file

// loaded is used to control if the config was loaded only one time
var loaded bool

// MustLoadConfig loads the app configuration and populates it into the Config variable.
// The default configuration file is config/config.yaml
// Environment variables with the prefix "SERVER_" in their names are also read automatically.
func MustLoadConfig() *ApiConfig {
	if loaded {
		panic("api config: is not possible to load config twice")
	}
	loaded = true
	v := viper.New()

	configFile := os.Getenv("API_CONFIG")
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// whit this, wee can have the config side by side with the binary: ./config/config.yaml
		v.AddConfigPath("./config")
		v.AddConfigPath("../../config")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// this will be the prefix for the env vars, following the 12factor: API_PORT=8081
	v.SetEnvPrefix("api")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprint("api config: failed to read the configuration file.", err))
	}
	cfg := &ApiConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Sprint("api config: failed to unmarshal the config file.", err))
	}
	return cfg
}

// DBConfig is a config for the relational db
type DBConfig struct {
	Dialect            string `mapstructure:"dialect"`
	DSN                string `mapstructure:"dsn"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
}

type SessionTokenConfig struct {
	Audience string `mapstructure:"audience"`
	Duration int    `mapstructure:"duration"`
	Issuer   string `mapstructure:"issuer"`
	Secret   string `mapstructure:"secret"`
}

type CorsConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type DocsConfig struct {
	Path string `mapstructure:"path"`
}

type LogConfig struct {
	Level   string           `mapstructure:"level"`
	Outputs LogOutputsConfig `mapstructure:"outputs"`
}

type LogOutputsConfig struct {
	StdoutEnable bool                `mapstructure:"stdout"`
	File         LogOutputFileConfig `mapstructure:"file"`
}

type LogOutputFileConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

type CacheConfig struct {
	DefaultExpire int `mapstructure:"defaultExpire"`
	DefaultPurge  int `mapstructure:"defaultPurge"`
}

type ApiConfig struct {
	DB DBConfig `mapstructure:"db"`
	// debug mode
	Debug bool `mapstructure:"debug"`
	// the server host bind, defaults to localhost
	Host string `mapstructure:"host"`
	// the server port, defaults to 8080
	Port int `mapstructure:"port"`
	// cors config
	Cors CorsConfig `mapstructure:"cors"`
	// SessionToken config used to generate the JWT
	SessionToken SessionTokenConfig `mapstructure:"session_token"`
	// Log config
	Log LogConfig `mapstructure:"log"`
	// Docs config
	Doc DocsConfig `mapstructure:"docs"`
	// Cache config
	Cache CacheConfig `mapstructure:"cache"`
}
