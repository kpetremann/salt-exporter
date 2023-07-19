//go:build e2e

package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	prom "github.com/prometheus/client_model/go"

	"github.com/prometheus/common/expfmt"
)

const exporterURL = "http://127.0.0.1:2112/metrics"

type Metric struct {
	Name   string
	Labels map[string]string
	Value  float64
}

// scrapeMetrics fetches the metrics from the running exporter
func scrapeMetrics(url string) (map[string]*prom.MetricFamily, error) {
	// Fetch the metrics from the running exporter
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body and parse the metrics
	var parser expfmt.TextParser
	parsed, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// getValue returns the value of a metric depending on its type
func getValue(metric *prom.Metric, metrictType prom.MetricType) float64 {
	switch metrictType {
	case prom.MetricType_COUNTER:
		return metric.Counter.GetValue()
	case prom.MetricType_GAUGE:
		return metric.Gauge.GetValue()
	case prom.MetricType_SUMMARY:
		return metric.Summary.GetSampleSum()
	case prom.MetricType_HISTOGRAM:
		return metric.Histogram.GetSampleSum()
	default:
		return 0
	}
}

// getMetrics returns the metrics from the running exporter
func getMetrics(url string) ([]Metric, error) {
	// Fetch the metrics from the running exporter
	parsed, err := scrapeMetrics(url)
	if err != nil {
		return nil, err
	}

	// Parse the metrics
	metrics := []Metric{}
	for _, metric := range parsed {
		for _, m := range metric.Metric {
			labels := map[string]string{}
			for _, l := range m.Label {
				labels[*l.Name] = *l.Value
			}
			metrics = append(metrics, Metric{*metric.Name, labels, getValue(m, *metric.Type)})
		}
	}

	return metrics, nil
}

// getValueForMetric returns the value of a metric with the given name and labels from the parsed metrics
func getValueForMetric(metrics []Metric, name string, labels map[string]string) (float64, error) {
	// Check if the metric exists
	for _, m := range metrics {
		if m.Name == name {
			// Check if the labels match
			match := true
			for k, v := range labels {
				if m.Labels[k] != v {
					match = false
					break
				}
			}
			if match {
				return m.Value, nil
			}
		}
	}

	// Return an error
	return 0, fmt.Errorf("%s with labels %v doesn't exist", name, labels)
}

// TestMetrics tests the metrics exposed by the exporter
func TestMetrics(t *testing.T) {
	// Fetch the metrics from the running exporter
	parsed, err := getMetrics(exporterURL)
	if err != nil {
		t.Fatal(err)
	}

	// Define the expected metrics
	expected := map[string][]Metric{
		"total responses": {
			{"salt_responses_total", map[string]string{"minion": "foo", "success": "false"}, 3},
			{"salt_responses_total", map[string]string{"minion": "foo", "success": "true"}, 3},
		},
		"execution modules": {
			{"salt_new_job_total", map[string]string{"function": "test.exception", "state": ""}, 1},
			{"salt_new_job_total", map[string]string{"function": "test.true", "state": ""}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "test.exception", "state": ""}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "test.true", "state": ""}, 1},
			{"salt_function_responses_total", map[string]string{"function": "test.exception", "state": "", "success": "false"}, 1},
			{"salt_function_responses_total", map[string]string{"function": "test.true", "state": "", "success": "true"}, 1},
		},
		"state modules": {
			{"salt_new_job_total", map[string]string{"function": "state.single", "state": "test.succeed_with_changes"}, 1},
			{"salt_new_job_total", map[string]string{"function": "state.single", "state": "test.fail_with_changes"}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "state.single", "state": "test.succeed_with_changes"}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "state.single", "state": "test.fail_with_changes"}, 1},
			{"salt_function_responses_total", map[string]string{"function": "state.single", "state": "test.succeed_with_changes", "success": "true"}, 1},
			{"salt_function_responses_total", map[string]string{"function": "state.single", "state": "test.fail_with_changes", "success": "false"}, 1},
		},
		"states": {
			{"salt_new_job_total", map[string]string{"function": "state.sls", "state": "test.succeed"}, 1},
			{"salt_new_job_total", map[string]string{"function": "state.sls", "state": "test.fail"}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "state.sls", "state": "test.succeed"}, 1},
			{"salt_expected_responses_total", map[string]string{"function": "state.sls", "state": "test.fail"}, 1},
			{"salt_function_responses_total", map[string]string{"function": "state.sls", "state": "test.succeed", "success": "true"}, 1},
			{"salt_function_responses_total", map[string]string{"function": "state.sls", "state": "test.fail", "success": "false"}, 1},
		},
	}

	// Get values from healthcheck
	healthcheckLabelsSuccess := map[string]string{"function": "status.ping_master", "state": "", "success": "true"}
	healthcheckValueTrue, err := getValueForMetric(parsed, "salt_function_responses_total", healthcheckLabelsSuccess)
	if err != nil {
		healthcheckValueTrue = 0
	}
	healthcheckLabelsFailed := map[string]string{"function": "status.ping_master", "state": "", "success": "false"}
	healthcheckValueFalse, err := getValueForMetric(parsed, "salt_function_responses_total", healthcheckLabelsFailed)
	if err != nil {
		healthcheckValueFalse = 0
	}
	healthcheckFindJobLabels := map[string]string{"function": "saltutil.find_job", "state": "", "success": "true"}
	healthcheckFindJob, err := getValueForMetric(parsed, "salt_function_responses_total", healthcheckFindJobLabels)
	if err != nil {
		healthcheckFindJob = 0
	}

	// Check if the expected metrics are present
	for testName, metrics := range expected {
		for _, e := range metrics {
			value, err := getValueForMetric(parsed, e.Name, e.Labels)

			// Remove events coming from docker healthcheck
			if testName == "total responses" {
				if e.Labels["success"] == "true" {
					value -= healthcheckValueTrue + healthcheckFindJob
				} else if e.Labels["success"] == "false" {
					value -= healthcheckValueFalse
				}
			}

			if err != nil {
				t.Error(err)
			} else if value != e.Value {
				t.Errorf("[%s] %s with labels %v = %f, got %f", testName, e.Name, e.Labels, e.Value, value)
			}
		}
	}
}
