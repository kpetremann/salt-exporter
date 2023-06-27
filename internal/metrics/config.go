package metrics

type MetricsConfig struct {
	// HealtMinions enable/disable the health functions/states metrics
	HealthMinions bool

	// HealthFunctionsFilter permits to limit the number of function exposed
	HealthFunctionsFilters []string

	// HealthFunctionsFilter permits to limit the number of state exposed
	HealthStatesFilters []string

	// Ignore test=True / mock=True events
	IgnoreTest, IgnoreMock bool
}
