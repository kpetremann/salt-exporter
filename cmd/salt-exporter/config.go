package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var flagConfigMapping = map[string]string{
	"host":     "listen-address",
	"port":     "listen-port",
	"tls":      "tls.enabled",
	"tls-cert": "tls.certificate",
	"tls-key":  "tls.key",
}

type Config struct {
	LogLevel string `mapstructure:"log-level"`

	ListenAddress string `mapstructure:"listen-address"`
	ListenPort    int    `mapstructure:"listen-port"`
	TLS           struct {
		Enabled     bool
		Key         string
		Certificate string
	}

	IgnoreTest bool `mapstructure:"ignore-test"`
	IgnoreMock bool `mapstructure:"ignore-mock"`

	HealthMinions         bool     `mapstructure:"health-minions"`
	HealthFunctionsFilter []string `mapstructure:"health-functions-filter"`
	HealthStatesFilter    []string `mapstructure:"health-states-filter"`
}

func parseFlags() {
	// flags
	flag.String("log-level", "info", "log level (debug, info, warn, error, fatal, panic, disabled)")

	flag.String("host", "", "listen address")
	flag.Int("port", 2112, "listen port")
	flag.Bool("tls", false, "enable TLS")
	flag.String("tls-cert", "", "TLS certificated")
	flag.String("tls-key", "", "TLS private key")

	flag.Bool("ignore-test", false, "ignore test=True events")
	flag.Bool("ignore-mock", false, "ignore mock=True events")

	flag.Bool("health-minions", true, "enable minion metrics")
	flag.String("health-functions-filter", "state.highstate",
		"apply filter on functions to monitor, separated by a comma")
	flag.String("health-states-filter", "highstate",
		"apply filter on states to monitor, separated by a comma")
	flag.Parse()
}

func getConfig() (Config, error) {
	// bind flags
	var allFlags []viperFlag
	flag.Visit(func(f *flag.Flag) {
		m := viperFlag{original: *f, alias: flagConfigMapping[f.Name]}
		allFlags = append(allFlags, m)
	})

	fSet := viperFlagSet{
		flags: allFlags,
	}
	if err := viper.BindFlagValues(fSet); err != nil {
		return Config{}, fmt.Errorf("flag binding failure: %w", err)
	}

	// bind configuration file
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return Config{}, fmt.Errorf("invalid config file: %w", err)
	}

	// extract configuration
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

func checkRequirements(cfg Config) error {
	if cfg.TLS.Enabled {
		if cfg.TLS.Certificate == "" {
			return errors.New("TLS Certificate not specified")
		}
		if cfg.TLS.Key == "" {
			return errors.New("TLS Private Key not specified")
		}
	}

	return nil
}

func ReadConfig() (Config, error) {
	var err error

	parseFlags()

	cfg, err := getConfig()
	if err != nil {
		return Config{}, err
	}

	err = checkRequirements(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
