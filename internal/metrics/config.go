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

	/*
		New job metrics
	*/

	SaltNewJobTotal struct {
		Enabled bool
	} `mapstructure:"salt_new_job_total"`

	SaltExpectedResponsesTotal struct {
		Enabled bool
	} `mapstructure:"salt_expected_responses_total"`

	/*
		Response metrics
	*/

	SaltFunctionResponsesTotal struct {
		Enabled        bool
		AddMinionLabel bool `mapstructure:"add-minion-label"`
	} `mapstructure:"salt_function_responses_total"`

	SaltScheduledJobReturnTotal struct {
		Enabled        bool
		AddMinionLabel bool `mapstructure:"add-minion-label"`
	} `mapstructure:"salt_scheduled_job_return_total"`

	SaltResponsesTotal struct {
		Enabled bool
	} `mapstructure:"salt_responses_total"`

	SaltFunctionStatus struct {
		Enabled bool
		Filters struct {
			Functions []string
			States    []string
		}
	} `mapstructure:"salt_function_status"`
}
