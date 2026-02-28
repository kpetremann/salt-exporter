package tui

import (
	"fmt"
	"time"

	"github.com/kpetremann/salt-exporter/pkg/event"
)

type item struct {
	title       string
	description string
	event       event.SaltEvent
	datetime    string
	sender      string
	state       string
	duration    *time.Duration
	eventJSON   string
	eventYAML   string
}

func (i item) Title() string {
	if i.event.Data.Retcode > 0 {
		return fmt.Sprintf("/!\\ %s", i.event.Tag)
	} else {
		return i.event.Tag
	}
}

func (i item) Description() string {
	out := fmt.Sprintf("%s - %s - %s", i.datetime, i.sender, i.event.Data.Fun)
	if i.state != "" {
		out = fmt.Sprintf("%s %s", out, i.state)
	}
	if i.event.TargetNumber > 0 {
		target := "targets"
		if i.event.TargetNumber == 1 {
			target = "target"
		}
		out = fmt.Sprintf("%s - %d %s", out, i.event.TargetNumber, target)
	}
	if i.duration != nil {
		out = fmt.Sprintf("%s - %.3fs", out, i.duration.Seconds())
	}
	return out
}

func (i item) FilterValue() string {
	return i.title + " " + i.Description() + " " + i.eventJSON
}
