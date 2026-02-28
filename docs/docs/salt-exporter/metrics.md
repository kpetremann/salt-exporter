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
| `salt_job_duration_seconds`       | `function`, `state`<br />(opt: `minion`)            | Last duration of a state job in seconds**                                 |
| `salt_health_last_heartbeat`      | `minion` | Last heartbeat from minion in UNIX timestamp
| `salt_health_minions_total`       |           | Total number of registered minions

\* more details in the [Function status](#function-status) section below.

\*\* more details in the [Job duration](#job-duration) section below.



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

## Job duration

`salt_job_duration_seconds` tracks the last known duration of a Salt state job. It is computed by summing the `duration` field of each state step in the return event.

This metric is only available for state functions (`state.sls`, `state.apply`, `state.highstate`, `state.single`) — execution modules do not report per-step durations and will not produce an observation.

Example for a highstate:
``` promql
salt_job_duration_seconds{function="state.highstate",state="highstate"} 1.498
```

With the optional minion label enabled:
``` promql
salt_job_duration_seconds{function="state.highstate",minion="node1",state="highstate"} 1.498
```

!!! warning

    Enabling `add-minion-label` multiplies the number of time series by the number of minions.
    Only enable it in environments with a small and bounded number of minions.

## Minions health

The exporter is supporting "hearbeat"-ing detection from minions which can be used to monitor for non-responding/dead minions. Under the hood it depends on Salt's beacons.
To ensure that all required minions are reported (even if there is no heartbeat from them yet), exporter needs access to the PKI directory of the Salt Master (by default `/etc/salt/pki/master`) where it watches for accepted minion's public keys (located under `/etc/salt/pki/master/minions`).
On startup, all currently accepted minions are added with last heartbeat set to current time. From this point forward, exporter is using __fsnotify__ to detect added or removed minions. This will ensure that once minion is added, it will be monitored for heartbeat and metric will be removed once minion is deleted from Salt master.

To use this functionality you'll need to add [`status` beacon](https://docs.saltproject.io/en/latest/ref/beacons/all/salt.beacons.status.html#:~:text=salt.-,beacons.,presence%20to%20be%20set%20up.) to each minion. It doesn't mater what functions will returned or the period. Exporter will just detect such events (in the format `salt/beacon/<minion id>/status`) and register the timestamp as last heartbeat.

### Detecting dead minions

The most simple way is (e.g. no heartbeat in last hour):
    ``` { .promql .copy }
    (time() - salt_health_last_heartbeat) > 3600
    ```
> __NOTE__: Above is assuming beacon interval is set to < 3600 seconds

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
??? example "Minions heartbeat"

    ```promql
    salt_health_last_heartbeat{minion="local"} 1703053536
    salt_health_last_heartbeat{minion="node1"} 1703052536

    salt_health_minions_total{} 2
    ```