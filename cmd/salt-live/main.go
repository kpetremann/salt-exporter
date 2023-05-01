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
	maxItems := flag.Int("max-events", 1000, "maximum events to keep in memory")
	versionCmd := flag.Bool("version", false, "print version")
	flag.Parse()

	if *versionCmd {
		printVersion()
		return
	}

	log.Logger = log.Output(nil)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eventChan := make(chan events.SaltEvent, 100)
	eventListener := events.NewEventListener(ctx, eventChan)
	go eventListener.ListenEvents(true)

	p := tea.NewProgram(tui.NewModel(eventChan, *maxItems), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
