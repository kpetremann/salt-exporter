---
title: Quickstart
---

# Salt Live

## Quickstart

`Salt Live` is a Terminal UI tool to watch events in real time.

It includes the following features:

* Hard filter from the CLI: filtered out events are discarded forever.
* Soft filter from the TUI: filtered out events are still kept in the buffer.
* Event details can be displayed in:
    * YAML
    * JSON
    * Golang structure
* The list is frozen when navigating the events.
    * It prevents annoying list updates when checking event details.
    * New events are still received and kept in the buffer.
    * Once the freeze is removed, the events are displayed in real-time.

## Installation

You can download the binary from the [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or install from source:

* latest published version:
    ``` { .sh .copy }
    go install github.com/kpetremann/salt-exporter/cmd/salt-live@latest
    ```
* latest commit (unstable):
    ``` { .sh .copy }
    go install github.com/kpetremann/salt-exporter/cmd/salt-live@main
    ```

## Credits

This tool uses these amazing libraries:

* [Bubble tea](https://github.com/charmbracelet/bubbletea)
* [Bubbles](https://github.com/charmbracelet/bubbles)
