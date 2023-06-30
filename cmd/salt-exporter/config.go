package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kpetremann/salt-exporter/internal/metrics"
	"github.com/spf13/viper"
)

const defaultLogLevel = "info"
const defaultPort = 2112
const defaultHealthMinion = true
const defaultHealthFunctionsFilter = "state.highstate"
const defaultHealthStatesFilter = "highstate"

var flagConfigMapping = map[string]string{
	"host":                    "listen-address",
	"port":                    "listen-port",
	"tls":                     "tls.enabled",
	"tls-cert":                "tls.certificate",
	"tls-key":                 "tls.key",
	"ignore-test":             "metrics.global.filters.ignore-test",
	"ignore-mock":             "metrics.global.filters.ignore-mock",
	"health-minions":          "metrics.health-minions",
	"health-functions-filter": "metrics.salt_function_status.filters.functions",
	"health-states-filter":    "metrics.salt_function_status.filters.states",
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

	Metrics metrics.Config
}

func parseFlags() {
	// flags
	flag.String("log-level", defaultLogLevel, "log level (debug, info, warn, error, fatal, panic, disabled)")

	flag.String("host", "", "listen address")
	flag.Int("port", defaultPort, "listen port")
	flag.Bool("tls", false, "enable TLS")
	flag.String("tls-cert", "", "TLS certificated")
	flag.String("tls-key", "", "TLS private key")

	flag.Bool("ignore-test", false, "ignore test=True events")
	flag.Bool("ignore-mock", false, "ignore mock=True events")

	flag.Bool("health-minions", defaultHealthMinion, "enable minion metrics")
	flag.String("health-functions-filter", defaultHealthStatesFilter,
		"apply filter on functions to monitor, separated by a comma")
	flag.String("health-states-filter", defaultHealthStatesFilter,
		"apply filter on states to monitor, separated by a comma")
	flag.Parse()
}

func setDefaults() {
	viper.SetDefault("log-level", defaultLogLevel)
	viper.SetDefault("listen-port", defaultPort)
	viper.SetDefault("metrics.health-minions", defaultHealthMinion)
	viper.SetDefault("metrics.salt_function_status.filters.functions", []string{defaultHealthFunctionsFilter})
	viper.SetDefault("metrics.salt_function_status.filters.states", []string{defaultHealthStatesFilter})
}

func getConfig() (Config, error) {
	setDefaults()

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
