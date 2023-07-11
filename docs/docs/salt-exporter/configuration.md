---
title: Configuration
---

# Configuration

The salt-exporter can be configured with flags, environments variables and configuration file.

!!! info

    The precedence order for the different methods is:

    * flags
    * environment variables
    * configuration file (config.yml)

## Configuration file

The exporter is looking for `config.yml`.

See below a full example of a configuration file:

```  { .yaml .copy }
log-level: "info"

listen-address: ""
listen-port: 2112

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

### Global parameters

| Parameter      | Default   | Description                                                        |
|----------------|-----------|--------------------------------------------------------------------|
| log-level      | `info`    | log level can be: debug, info, warn, error, fatal, panic, disabled |
| listen-address | `0.0.0.0` | listening address                                                  |
| listen-port    | `2112`    | listening port                                                     |

### TLS settings

All parameters below are in the `tls` section of the configuration.

| Parameter   | Default | Description                                 |
|-------------|---------|---------------------------------------------|
| enabled     | `false` | enables/disables TLS on the metrics webserver |
| key         |         | TLS key for the metrics webserver           |
| certificate |         | TLS certificate for the metrics webserver   |

### Metrics global settings

All parameters below are in the `metrics.global` section of the configuration.

| Parameter           | Default | Description               |
|---------------------|---------|---------------------------|
| filters.ignore-test | `false` | ignores `test=True` events |
| filters.ignore-mock | `false` | ignores `mock=True` events |

### Metrics configuration

All parameters below are in the `metrics` section of the configuration.

| Parameter | Default           | Description |
|-----------|-------------------|-------------------------------------------------------------------|
| `<metrics_name>`.enabled | `true` | enables or disables a metric |
| `<metrics_name>`.add-minion-label<br /><br />Only for:<br /><ul><li>`salt_function_responses_total`</li><li>`salt_scheduled_job_return_total`</li></ul> | `false` | adds minion label<br />_not recommended<br />can lead to cardinality issues_ |
| salt_function_status.filters.function | `state.highstate` | updates the metric only if the event function matches the filter |
| salt_function_status.filters.states | `highstate` | updates the metric only if the event state matches the filter |

## Alternative methods

### Environment variables

All settings available in the configuration file can be set as environment variables, but:

* all variables must be prefixed by `SALT_`
* uppercase only
* `-` in the configuration file becomes a `_`
* `__` is the level separator

For example, the equivalent of this config file:

``` yaml
log-level: "info"
tls:
  enabled: true
metrics:
  global:
    filters:
      ignore-test: true
```

is:

``` shell
SALT_LOG_LEVEL="info"
SALT_TLS__ENABLED=true
SALT_METRICS__GLOBAL__FILTERS__IGNORE_TEST=true
```

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

