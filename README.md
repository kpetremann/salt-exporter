[![status](https://img.shields.io/badge/status-beta-orange)](https://github.com/kpetremann/salt-exporter)
[![Go](https://img.shields.io/github/go-mod/go-version/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter)
[![CI](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter/blob/main/LICENSE)

> This exporter comes with a TUI to watch event in real time. See [Salt live](cmd/salt-live/README.md)

# Salt Exporter
<img align="right" width="120px" src="https://raw.githubusercontent.com/kpetremann/salt-exporter/main/img/salt-exporter.png" />

## Quickstart

Salt Exporter is a Prometheus export for Salt events. It exposes relevant metrics regarding jobs and results.

Working out of the box: you just need to run the exporter on the same server than the Salt Master.

Notes:

> If you did not setup the external_auth parameter, you need to run the exporter with the same user running the Salt Master. If it is set, any user will do the trick.

> This has only be tested on Linux

## Installation

Just use the binary from [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or, install via source:
- latest release: `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@latest`
- unstable: `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@main`

## Usage

```
./salt-exporter
  -health-functions-filter string
    	Apply filter on functions to monitor, separated by a comma (default "state.highstate")
  -health-minions
    	Enable minion metrics (default true)
  -health-states-filter string
    	Apply filter on states to monitor, separated by a comma (default "highstate")
  -host string
    	listen address
  -ignore-mock
    	ignore mock=True events
  -ignore-test
    	ignore test=True events
  -log-level string
    	log level (debug, info, warn, error, fatal, panic, disabled) (default "info")
  -port int
    	listen port (default 2112)
  -tls
    	enable TLS
  -tls-cert string
    	TLS certificated
  -tls-key string
    	TLS private key
```

## Features

Supported tags:
* `salt/job/<jid>/new`
* `salt/job/<jid>/ret/<*>`
* `salt/run/<jid>/new`
* `salt/run/<jid>/ret/<*>`

It extracts and exposes:
* the execution module, to `function` label
* the states when using state.sls/state.apply/state.highstate, to `state` label
* same info for the runners

## Exposed metrics

### Example

execution modules:
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
```

states (state.sls/apply/highstate) and states module (state.single):
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
```

> Note: `salt_responses_total{minion="local",success="true"}` metrics can be disabled using `-health-minions` flag.

### Minions job status

By default, a Salt highstate will generate a status metric:
```
salt_function_status{function="state.highstate",minion="node1",state="highstate"} 1
```
* `1` means that the last time this couple of function/state were executed, the return was `successful`
* `0` means that the last time this couple of function/state were executed, the return was `failed`

You will find an example of Prometheus alerts that could be used with this metric in the `prometheus_alerts` directory.

The health metrics can be customized by using the `-health-functions-filter` and `-health-states-filter`, example of usage:
```
./salt-exporter -health-functions-filter=test.ping,state.apply -health-states-filter=""
```

This will only generate a metric for the `test.ping` function executed:
```
salt_function_status{function="test.ping",minion="node1",state=""} 1
```

You can disable all the health metrics with this config switch:
```./salt-exporter -health-minions=false```

Note: this also works for scheduled jobs.

### `salt/job/<jid>/new`

It increases:
* `salt_new_job_total`
* `salt_expected_responses_total` by the number of target in the new job

### `salt/job/<jid>/ret/<*>`

Usually, it will increase the `salt_responses_total` (per minion) and `salt_function_responses_total` (per function) counters.

However, if it is of a scheduled job feedback, it increases `salt_scheduled_job_return_total` instead.

#### Why separating `salt_responses_total` and `salt_scheduled_job_return_total`

One of the goal is to be able to calculate the number of missed response, without doing the matching manually between the target of a new job and the minion responses.

This is why scheduled job are in a dedicated metric. Once scheduled, a job is only executing autonomously on Minion side, hence there is no new job request for each scheduled job response. Said differently, if there was no differences made, we would end up with more responses than expected responses.

#### Estimate missing responses

Missing responses can be simply calculated by doing the difference between `salt_expected_responses_total` and `salt_responses_total`.

It can be joined on function label to have details per executed module.

## Upcoming features

* metric regarding IPC connectivity

## Estimated performance

According to some simple benchmark, for a simple event, it takes:
* ~60us for parsing
* ~9us for converting to Prometheus metric

So with a security margin, we can estimate an event should take 100us maximum.

Roughly, the exporter should be able to handle about 10kQps.

For a base of 1000 Salt minions, it should be able to sustain 10 jobs per minion per second, which is a quite high for Salt.

If needed, the exporter can easily scale more up by doing the parsing in dedicated goroutines, the limiting factor being the Prometheus metric update (~9us).
