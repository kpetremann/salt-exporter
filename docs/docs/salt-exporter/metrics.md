---
title: Metrics
---

# Exposed metrics

## Examples

??? example "Execution modules"

    ```
    # HELP salt_expected_responses_total Total number of expected minions responses
    # TYPE salt_expected_responses_total counter
    salt_expected_responses_total{function="cmd.run", state=""} 6
    salt_expected_responses_total{function="test.ping", state=""} 6

    # HELP salt_function_responses_total Total number of response per function processed
    # TYPE salt_function_responses_total counter
    salt_function_responses_total{function="cmd.run",state="",success="true"} 6
    salt_function_responses_total{function="test.ping",state="",success="true"} 6

    # HELP salt_new_job_total Total number of new job processed
    # TYPE salt_new_job_total counter
    salt_new_job_total{function="cmd.run",state="",success="false"} 3
    salt_new_job_total{function="test.ping",state="",success="false"} 3

    # HELP salt_responses_total Total number of response job processed
    # TYPE salt_responses_total counter
    salt_responses_total{minion="local",success="true"} 6
    salt_responses_total{minion="node1",success="true"} 6

    # HELP salt_scheduled_job_return_total Total number of scheduled job response
    # TYPE salt_scheduled_job_return_total counter
    salt_scheduled_job_return_total{function="cmd.run",minion="local",state="",success="true"} 2
    ```

??? example "States and state modules"

    States (state.sls/apply/highstate) and state module (state.single):

    ```
    salt_expected_responses_total{function="state.apply",state="highstate"} 1
    salt_expected_responses_total{function="state.highstate",state="highstate"} 2
    salt_expected_responses_total{function="state.sls",state="test"} 1
    salt_expected_responses_total{function="state.single",state="test.nop"} 3

    salt_function_responses_total{function="state.apply",state="highstate",success="true"} 1
    salt_function_responses_total{function="state.highstate",state="highstate",success="true"} 2
    salt_function_responses_total{function="state.sls",state="test",success="true"} 1
    salt_function_responses_total{function="state.single",state="test.nop",success="true"} 3

    salt_function_status{minion="node1",function="state.highstate",state="highstate"} 1

    salt_new_job_total{function="state.apply",state="highstate",success="false"} 1
    salt_new_job_total{function="state.highstate",state="highstate",success="false"} 2
    salt_new_job_total{function="state.sls",state="test",success="false"} 1
    salt_new_job_total{function="state.single",state="test.nop",success="true"} 3

    salt_scheduled_job_return_total{function="state.sls",minion="local",state="test",success="true"} 3
    ```

## Metrics

| Event tag | Metric increase |
|----------------------|--|
| <ul><li>`salt/job/<jid>/new`</li><li>`salt/run/<jid>/new`</li></ul> | <ul><li>`salt_new_job_total`</li><li>`salt_expected_responses_total` by the number of targets in the new job</li></ul>
| <ul><li>`salt/job/<jid>/ret/<*>`</li><li>`salt/run/<jid>/new`</li></ul> | <ul><li>`salt_responses_total`</li><li>`salt_function_responses_total` or `salt_scheduled_job_return_total`<li>`salt_function_status`</li></li></ul>

## Function status

By default, a Salt highstate will generate the following metric:
```
salt_function_status{function="state.highstate",minion="node1",state="highstate"} 1
```

* `1` the last function/state execution was `successful`
* `0` the last function/state execution has `failed`

You will find an example of Prometheus alerts that could be used [here](https://github.com/kpetremann/salt-exporter/blob/main/prometheus_alerts/highstate.yaml).

See the [configuration page](./configuration.md) if you want to watch other functions/states, or if you want to disable this metric.

## Labels

The exporter exposes the label for both classic jobs and runners.

| Prometheus label | Salt information               |
|------------------|--------------------------------|
| `function`       | execution module               |
| `state`          | state and state module         |
| `minion`         | minion which send the response |
| `success`        | status of the job              |

## Estimated missing responses

Simple way:
    ```
    salt_expected_responses_total - on(function) salt_function_responses_total
    ```

More advanced:
    ```
    sum by (instance, function, state) (
        increase(salt_expected_responses_total{function=~"$function", state=~"$state"}[$__rate_interval])
    )
    - sum by (instance, function, state) (
        increase(salt_function_responses_total{function=~"$function", state=~"$state"}[$__rate_interval])
    )
    ```