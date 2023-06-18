package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpetremann/salt-exporter/internal/tui"
	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/kpetremann/salt-exporter/pkg/listener"
	"github.com/kpetremann/salt-exporter/pkg/parser"
	"github.com/rs/zerolog/log"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "development"

func printVersion() {
	fmt.Println("Version: ", version)
}

func main() {
	maxItems := flag.Int("max-events", 1000, "maximum events to keep in memory")
	bufferSize := flag.Int("buffer-size", 1000, "buffer size in number of events")
	filter := flag.String("hard-filter", "", "filter when received (filtered out events are discarded forever)")
	versionCmd := flag.Bool("version", false, "print version")
	debug := flag.Bool("debug", false, "enable debug mode (log to debug.log)")
	flag.Parse()

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if *versionCmd {
		printVersion()
		return
	}

	log.Logger = log.Output(nil)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eventChan := make(chan event.SaltEvent, *bufferSize)
	parser := parser.NewEventParser(true)
	eventListener := listener.NewEventListener(ctx, parser, eventChan)
	go eventListener.ListenEvents()

	p := tea.NewProgram(tui.NewModel(eventChan, *maxItems, *filter), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
