package metrics

type Config struct {
	// HealtMinions enable/disable the health functions/states metrics
	HealthMinions bool `mapstructure:"health-minions"`

	Global struct {
		Filters struct {
			IgnoreTest bool `mapstructure:"ignore-test"`
			IgnoreMock bool `mapstructure:"ignore-mock"`
		}
	}

	SaltFunctionStatus struct {
		Filters struct {
			Functions []string
			States    []string
		}
	} `mapstructure:"salt_function_status"`
}
