package metrics

type MetricsConfig struct {
	HealthMinions          bool
	HealthFunctionsFilters []string
	HealthStatesFilters    []string
}
