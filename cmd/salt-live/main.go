package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpetremann/salt-exporter/internal/tui"
	"github.com/kpetremann/salt-exporter/pkg/events"
	"github.com/rs/zerolog/log"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "development"

func printVersion() {
	fmt.Println("Version: ", version)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	maxItems := flag.Int("max-events", 1000, "maximum events to keep in memory")
	bufferSize := flag.Int("buffer-size", 1000, "buffer size in number of events")
	filter := flag.String("hard-filter", "", "filter when received (filtered out events are discarded forever)")
	versionCmd := flag.Bool("version", false, "print version")
	flag.Parse()

	if *versionCmd {
		printVersion()
		return
	}

	log.Logger = log.Output(nil)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eventChan := make(chan events.SaltEvent, *bufferSize)
	eventListener := events.NewEventListener(ctx, eventChan)
	go eventListener.ListenEvents(true)

	p := tea.NewProgram(tui.NewModel(eventChan, *maxItems, *filter), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
