---
title: Configuration
---

# Configuration

!!! info

    The precedence order for the different configuration methods is:

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
