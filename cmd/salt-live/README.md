# Salt live

## Quickstart

Salt live is a Terminal UI tool to watch event in real time.

It is similar to `salt-run state.event pretty=True` runner, but with additional features:

- hard filter: filtered out events are discarded forever
- soft filter: filtered out events are still kept in the buffer
- display some info from the salt-exporter parser
- event details in YAML or JSON
- freeze the refresh list to navigate the events, while still receiving new ones

## Installation

Just use the binary from [Github releases](https://github.com/kpetremann/salt-exporter/releases) page.

Or, install via source:
- latest release: `go install github.com/kpetremann/salt-exporter/cmd/salt-live@latest`
- unstable: `go install github.com/kpetremann/salt-exporter/cmd/salt-live@main`

## Credits

This tool uses the amazing [Bubble tea](https://github.com/charmbracelet/bubbletea) TUI framework and [Bubbles](https://github.com/charmbracelet/bubbles) components.
