---
title: Quickstart
---

# Salt Exporter
<!-- <img align="right" width="120px" src="https://raw.githubusercontent.com/kpetremann/salt-exporter/main/img/salt-exporter.png" /> -->

## Installation

You can download the binary from the [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or, install them from source:

* latest release:
    ``` shell
    go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@latest
    ```

* unstable:
    ``` shell
    go install github.com/kpetremann/salt-exporter/cmd/salt-exporter@main
    ```

!!! warning "Deprecation notice"

    The following flags are deprecated:

    * `-health-minions`
    * `-health-functions-filter`
    * `-health-states-filter`

    They should be replaced by configuring metrics in the `config.yml` file.

    The equivalent of:
    ```
    ./salt-exporter -health-minions -health-functions-filter "func1,func2" -health-states-filter "state1,state2"`
    ```

    is:
    ``` yaml
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

The exporter will run out of the box:
```./salt-exporter```

!!! note

    You need to run the exporter with the user running the Salt master.

The feature list can be found [here]("./features.md")

If needed, you can [configure]("./configuration.md") the exporter to better match your needs.
