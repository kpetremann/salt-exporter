[![status](https://img.shields.io/badge/status-beta-orange)](https://github.com/kpetremann/salt-exporter)
[![Go](https://img.shields.io/github/go-mod/go-version/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter)
[![CI](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter/blob/main/LICENSE)

## Salt Live

!!! note ""

    _`salt-run state.event pretty=True` under steroids_

Salt Exporter comes with a Salt Live. This is a Terminal UI tool to watch event in real time.

## Salt Exporter

<img align="right" width="120px" src="https://raw.githubusercontent.com/kpetremann/salt-exporter/main/img/salt-exporter.png" />

`Salt Exporter` is a Prometheus exporter for [Saltstack](https://github.com/saltstack/salt) events. It exposes relevant metrics regarding jobs and results.

This exporter is passive. It does not use the Salt API.

It works out of the box: you just need to run the exporter on the same server as the Salt Master.

```
$ ./salt-exporter
```

```
$ curl -s 127.0.0.1:2112/metrics

salt_expected_responses_total{function="cmd.run", state=""} 6
salt_expected_responses_total{function="state.sls",state="test"} 1

salt_function_responses_total{function="cmd.run",state="",success="true"} 6
salt_function_responses_total{function="state.sls",state="test",success="true"} 1

salt_function_status{minion="node1",function="state.highstate",state="highstate"} 1

salt_new_job_total{function="cmd.run",state="",success="false"} 3
salt_new_job_total{function="state.sls",state="test",success="false"} 1

salt_responses_total{minion="local",success="true"} 6
salt_responses_total{minion="node1",success="true"} 6

salt_scheduled_job_return_total{function="state.sls",minion="local",state="test",success="true"} 2
```

## Deprecation notice

`-health-minions`, `health-functions-filter` and `health-states-filter` are deprecated.
They should be replaced by configuring metrics in the `config.yml` file.

The equivalent of `./salt-exporter -health-minions -health-functions-filter "func1,func2" -health-states-filter "state1,state2"` is:

```yaml
metrics:
  salt_responses_total:
    enabled: true

  salt_function_status:
    enabled: true
    filters:
      functions:
        - "func1"
        - "func2"
      states:
        - "state1"
        - "state2"
```

## Installation

Just use the binary from [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or, install via source:
- latest release: `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@latest`
- unstable: `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@main`

## Usage

### Run

Simply run:
```./salt-exporter```

The exporter can be configured in different ways, with the following precedence order:
* flags
* environment variables
* configuration file (config.yml)

### Flags

```
./salt-exporter -help
  -health-functions-filter string
        [DEPRECATED] apply filter on functions to monitor, separated by a comma (default "highstate")
  -health-states-filter string
        [DEPRECATED] apply filter on states to monitor, separated by a comma (default "highstate")
  -health-minions
        [DEPRECATED] enable minion metrics (default true)
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

### Environment variables

All settings available in the configuration file can be set as environment variables, but:

* all variables must be prefixed by `SALT_`
* uppercase only
* `-` in the configuration file becomes a `_`
* `__` is the level separator

For example, the equivalent of this config file:

```yaml
log-level: "info"
tls:
  enabled: true
metrics:
  global:
    filters:
      ignore-test: true
```

is:

```
SALT_LOG_LEVEL="info"
SALT_TLS__ENABLED=true
SALT_METRICS__GLOBAL__FILTERS__IGNORE_TEST=true
```

### Configuration file

The exporter is looking for `config.yml`.

See below a full example of a configuration file:

```yaml
listen-address: ""
listen-port: 2112

log-level: "info"
tls:
  enabled: true
  key: "/path/to/key"
  certificate: "/path/to/certificate"


metrics:
  global:
    filters:
      ignore-test: false
      ignore-mock: false

  salt_new_job_total:
    enabled: true

  salt_expected_responses_total:
    enabled: true

  salt_function_responses_total:
    enabled: true
    add-minion-label: false  # not recommended in production

  salt_scheduled_job_return_total:
    enabled: true
    add-minion-label: false  # not recommended in production

  salt_responses_total:
    enabled: true

  salt_function_status:
    enabled: true
    filters:
      functions:
        - "state.highstate"
      states:
        - "highstate"
```
