---
title: Metrics
---

# Exposed metrics

## Metrics

??? info "Supported Salt event tags"

    Each Salt event having a tag in this list will update the metrics:

    * `salt/job/<jid>/new`
    * `salt/job/<jid>/ret/<*>`
    * `salt/run/<jid>/new`
    * `salt/run/<jid>/ret/<*>`

| Metric                            | Labels                                              | Description                                                               |
|-----------------------------------|-----------------------------------------------------|---------------------------------------------------------------------------|
| `salt_new_job_total`              | `function`, `state`                                 | Total number of new jobs                                                  |
| `salt_expected_responses_total`   | `function`, `state`                                 | Counter incremented by the number of targeted minion for each new job     |
| `salt_function_responses_total`   | `function`, `state`, `success`<br />(opt: `minion`) | Total number of job responses by function, state and success<br />        |
| `salt_scheduled_job_return_total` | `function`, `state`, `success`<br />(opt: `minion`) | Counter incremented each time a minion sends a scheduled job result       |
| `salt_responses_total`            | `minion`, `success`                                 | Total number of job responses<br />_including scheduled_job responses_    |
| `salt_function_status`            | `function`, `state`, `minion`                       | Last status of a job execution*                                           |

\* more details in the section below.



## Labels details

The exporter exposes the label for both classic jobs and runners.

| Prometheus label | Salt information               |
|------------------|--------------------------------|
| `function`       | execution module               |
| `state`          | state and state module         |
| `minion`         | minion sending the response    |
| `success`        | job status                     |

## Function status

By default, a Salt highstate generates the following metric:
``` promql
salt_function_status{function="state.highstate",minion="node1",state="highstate"} 1
```

The value can be:

* `1` the last function/state execution was `successful`
* `0` the last function/state execution has `failed`

You can find an example of Prometheus alerts that could be used [here](https://github.com/kpetremann/salt-exporter/blob/main/prometheus_alerts/highstate.yaml).

See the [configuration page](./configuration.md) if you want to watch other functions/states, or if you want to disable this metric.

## How to estimate missing responses

Simple way:
    ``` { .promql .copy }
    salt_expected_responses_total - on(function) salt_function_responses_total
    ```

More advanced:
    ``` { .promql .copy }
    sum by (instance, function, state) (
        increase(salt_expected_responses_total{function=~"$function", state=~"$state"}[$__rate_interval])
    )
    - sum by (instance, function, state) (
        increase(salt_function_responses_total{function=~"$function", state=~"$state"}[$__rate_interval])
    )
    ```

## Examples

??? example "Execution modules"

    ``` promql
    # HELP salt_expected_responses_total Total number of expected minions responses
    # TYPE salt_expected_responses_total counter
    salt_expected_responses_total{function="cmd.run", state=""} 6
    salt_expected_responses_total{function="test.ping", state=""} 6

    # HELP salt_function_responses_total Total number of responses per function processed
    # TYPE salt_function_responses_total counter
    salt_function_responses_total{function="cmd.run",state="",success="true"} 6
    salt_function_responses_total{function="test.ping",state="",success="true"} 6

    # HELP salt_new_job_total Total number of new jobs processed
    # TYPE salt_new_job_total counter
    salt_new_job_total{function="cmd.run",state=""} 3
    salt_new_job_total{function="test.ping",state=""} 3

    # HELP salt_responses_total Total number of responses
    # TYPE salt_responses_total counter
    salt_responses_total{minion="local",success="true"} 6
    salt_responses_total{minion="node1",success="true"} 6

    # HELP salt_scheduled_job_return_total Total number of scheduled job responses
    # TYPE salt_scheduled_job_return_total counter
    salt_scheduled_job_return_total{function="cmd.run",minion="local",state="",success="true"} 2
    ```

??? example "States and state modules"

    States (state.sls/apply/highstate) and state module (state.single):

    ``` promql
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
