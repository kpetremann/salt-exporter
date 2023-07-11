---
title: Quickstart
---

# Salt Exporter
<!-- <img align="right" width="120px" src="https://raw.githubusercontent.com/kpetremann/salt-exporter/main/img/salt-exporter.png" /> -->

## Installation

You can download the binary from the [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or install from source:

* latest published version:
    ``` { .sh .copy }
    go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@latest
    ```

* latest commit (unstable):
    ``` { .sh .copy }
    go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@main
    ```

!!! warning "Deprecation notice"

    The following flags are deprecated:

    * `-health-minions`
    * `-health-functions-filter`
    * `-health-states-filter`

    They should be replaced by metrics configuration in the `config.yml` file.

    The equivalent of:
    ``` shell
    ./salt-exporter -health-minions -health-functions-filter "func1,func2" -health-states-filter "state1,state2"`
    ```

    is:
    ``` { .yaml .copy }
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


## Usage

The exporter runs out of the box:
```./salt-exporter```

!!! note

    You need to run the exporter with the user running the Salt master.

!!! example "Examples of configuration options"

    * All metrics can be either enabled or disabled.
    * You can add a minion label to some metrics (not recommended on large environment as it could lead to cardinality issues).
    * You can filter out `test=true`/`mock=true` events, useful to ignore tests.
    * ... more options can be found in the [configuration page](./configuration.md)
