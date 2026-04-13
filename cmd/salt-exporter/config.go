package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kpetremann/salt-exporter/internal/metrics"
	"github.com/kpetremann/salt-exporter/pkg/listener"
	"github.com/spf13/viper"
)

const defaultConfigFilename = "config.yml"
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
	IPCFile       string `mapstructure:"ipc-file"`
	PKIDir        string `mapstructure:"pki-dir"`
	TLS           struct {
		Enabled     bool
		Key         string
		Certificate string
	}

	Metrics metrics.Config
}

func parseFlags() (string, bool) {
	// flags
	versionCmd := flag.Bool("version", false, "print version")
	configFile := flag.String("config-file", defaultConfigFilename, "config filepath")
	flag.String("log-level", defaultLogLevel, "log level (debug, info, warn, error, fatal, panic, disabled)")
	flag.String("host", "", "listen address")
	flag.Int("port", defaultPort, "listen port")
	flag.String("ipc-file", listener.DefaultIPCFilepath, "file location of the salt-master event bus")
	flag.Bool("tls", false, "enable TLS")
	flag.String("tls-cert", "", "TLS certificated")
	flag.String("tls-key", "", "TLS private key")

	flag.Bool("ignore-test", false, "ignore test=True events")
	flag.Bool("ignore-mock", false, "ignore mock=True events")

	// deprecated flag
	healthMinions := flag.Bool("health-minions", defaultHealthMinion, "[DEPRECATED] enable minion metrics")
	flag.String("health-functions-filter", defaultHealthStatesFilter,
		"[DEPRECATED] apply filter on functions to monitor, separated by a comma")
	flag.String("health-states-filter", defaultHealthStatesFilter,
		"[DEPRECATED] apply filter on states to monitor, separated by a comma")
	flag.Parse()

	if *versionCmd {
		if version == "unknown" {
			version = fmt.Sprintf("v%s", version)
		}
		fmt.Println("Version:", version)
		fmt.Println("Build date:", date)
		fmt.Println("Commit:", commit)
		os.Exit(0)
	}

	return *configFile, *healthMinions
}

func setDefaults(healthMinions bool) {
	viper.SetDefault("log-level", defaultLogLevel)
	viper.SetDefault("listen-port", defaultPort)
	viper.SetDefault("ipc-file", listener.DefaultIPCFilepath)
	viper.SetDefault("pki-dir", listener.DefaultPKIDirpath)
	viper.SetDefault("metrics.health-minions", defaultHealthMinion)
	viper.SetDefault("metrics.salt_new_job_total.enabled", true)
	viper.SetDefault("metrics.salt_expected_responses_total.enabled", true)
	viper.SetDefault("metrics.salt_function_responses_total.enabled", true)
	viper.SetDefault("metrics.salt_scheduled_job_return_total.enabled", true)
	viper.SetDefault("metrics.salt_function_status.enabled", healthMinions) // TODO: true once health-minions will be removed
	viper.SetDefault("metrics.salt_responses_total.enabled", healthMinions) // TODO: true once health-minions will be removed
	viper.SetDefault("metrics.salt_job_duration_seconds.enabled", true)
	viper.SetDefault("metrics.salt_function_status.filters.functions", []string{defaultHealthFunctionsFilter})
	viper.SetDefault("metrics.salt_function_status.filters.states", []string{defaultHealthStatesFilter})
}

func getConfig(configFileName string, healthMinions bool) (Config, error) {
	setDefaults(healthMinions)

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
	if filepath.IsAbs(configFileName) {
		viper.SetConfigFile(configFileName)
	} else {
		ext := filepath.Ext(configFileName)
		viper.SetConfigName(strings.TrimSuffix(configFileName, ext))
		viper.SetConfigType(strings.TrimPrefix(ext, "."))
		viper.AddConfigPath(".")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
	viper.SetEnvPrefix("SALT")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok { //nolint: errorlint // ConfigFileNotFoundError not implementing error
			return Config{}, fmt.Errorf("invalid config file: %w", err)
		}
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

	configFileName, healthMinions := parseFlags()

	cfg, err := getConfig(configFileName, healthMinions)
	if err != nil {
		return Config{}, err
	}

	err = checkRequirements(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
