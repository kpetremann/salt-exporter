package metrics

import (
	"strconv"

	"github.com/kpetremann/salt-exporter/internal/filters"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Registry struct {
	config Config

	newJobTotal            *prometheus.CounterVec
	expectedResponsesTotal *prometheus.CounterVec

	functionResponsesTotal  *prometheus.CounterVec
	scheduledJobReturnTotal *prometheus.CounterVec

	responseTotal  *prometheus.CounterVec
	functionStatus *prometheus.GaugeVec
}

func NewRegistry(config Config) Registry {
	functionResponsesTotalLabels := []string{"function", "state", "success"}
	if config.SaltFunctionResponsesTotal.AddMinionLabel {
		functionResponsesTotalLabels = append([]string{"minion"}, functionResponsesTotalLabels...)
	}

	scheduledJobReturnTotalLabels := []string{"function", "state", "success"}
	if config.SaltScheduledJobReturnTotal.AddMinionLabel {
		scheduledJobReturnTotalLabels = append([]string{"minion"}, scheduledJobReturnTotalLabels...)
	}

	return Registry{
		config: config,

		newJobTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "salt_new_job_total",
				Help: "Total number of new jobs processed",
			},
			[]string{"function", "state", "success"},
		),

		expectedResponsesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "salt_expected_responses_total",
				Help: "Total number of expected minions responses",
			},
			[]string{"function", "state"},
		),

		functionResponsesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "salt_function_responses_total",
				Help: "Total number of responses per function processed",
			},
			functionResponsesTotalLabels,
		),

		scheduledJobReturnTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "salt_scheduled_job_return_total",
				Help: "Total number of scheduled job responses",
			},
			scheduledJobReturnTotalLabels,
		),

		responseTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "salt_responses_total",
				Help: "Total number of responses",
			},
			[]string{"minion", "success"},
		),

		functionStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "salt_function_status",
				Help: "Last function/state success, 0=Failed, 1=Success",
			},
			[]string{"minion", "function", "state"},
		),
	}
}

func (r Registry) IncreaseNewJobTotal(function, state string, success bool) {
	if r.config.SaltNewJobTotal.Enabled {
		r.newJobTotal.WithLabelValues(function, state, strconv.FormatBool(success)).Inc()
	}
}

func (r Registry) IncreaseExpectedResponsesTotal(function, state string, value float64) {
	if r.config.SaltExpectedResponsesTotal.Enabled {
		r.expectedResponsesTotal.WithLabelValues(function, state).Add(value)
	}
}

func (r Registry) IncreaseFunctionResponsesTotal(function, state, minion string, success bool) {
	labels := []string{function, state, strconv.FormatBool(success)}
	if r.config.SaltFunctionResponsesTotal.AddMinionLabel {
		labels = append([]string{minion}, labels...)
	}

	if r.config.SaltFunctionResponsesTotal.Enabled {
		r.functionResponsesTotal.WithLabelValues(labels...).Inc()
	}
}

func (r Registry) IncreaseScheduledJobReturnTotal(function, state, minion string, success bool) {
	labels := []string{function, state, strconv.FormatBool(success)}
	if r.config.SaltScheduledJobReturnTotal.AddMinionLabel {
		labels = append([]string{minion}, labels...)
	}

	if r.config.SaltScheduledJobReturnTotal.Enabled {
		r.scheduledJobReturnTotal.WithLabelValues(labels...).Inc()
	}
}

func (r Registry) IncreaseResponseTotal(minion string, success bool) {
	if r.config.SaltResponsesTotal.Enabled {
		r.responseTotal.WithLabelValues(minion, strconv.FormatBool(success)).Inc()
	}
}

func (r Registry) SetFunctionStatus(minion, function, state string, success bool) {
	if !r.config.SaltFunctionStatus.Enabled {
		return
	}
	if !filters.Match(function, r.config.SaltFunctionStatus.Filters.Functions) {
		return
	}
	if !filters.Match(state, r.config.SaltFunctionStatus.Filters.States) {
		return
	}

	r.functionStatus.WithLabelValues(minion, function, state).Set(boolToFloat64(success))
}
