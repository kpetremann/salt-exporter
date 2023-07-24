[![Latest](https://img.shields.io/github/v/release/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter)
[![CI](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/kpetremann/salt-exporter/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/kpetremann/salt-exporter)](https://github.com/kpetremann/salt-exporter/blob/main/LICENSE)


## Salt Live

> _`salt-run state.event pretty=True` under steroids_

Salt Exporter comes with `Salt Live`. This is a Terminal UI tool to watch events in real time.

<img src="./docs/docs/demo/tui-overview.gif" alt="demo" width="500" />


## Salt Exporter

`Salt Exporter` is a Prometheus exporter for [Saltstack](https://github.com/saltstack/salt) events. It exposes relevant metrics regarding jobs and results.

This exporter is passive. It does not use the Salt API.

It works out of the box: you just need to run the exporter on the same user as the Salt Master.

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

### Deprecation notice

`-health-minions`, `health-functions-filter` and `health-states-filter` are deprecated.
They should be replaced by metrics configuration in the `config.yml` file.

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

### Installation

Just use the binary from [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or, install from source:
- latest published version: `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@latest`
- latest commit (unstable): `go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@main`

### Usage

Simply run:
```./salt-exporter```

The exporter can be configured in different ways, with the following precedence order:
* flags
* environment variables
* configuration file (config.yml)

See the [official documentation](https://kpetremann.github.io/salt-exporter) for more details
