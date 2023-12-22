package tui

import (
	"sync"

	teaList "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var safeMu sync.Mutex

func safeUpdate(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		safeMu.Lock()
		defer safeMu.Unlock()
		if cmd != nil {
			return cmd()
		}
		return nil
	}
}

func newDelegate(maxItems int) teaList.DefaultDelegate {
	d := teaList.NewDefaultDelegate()

	buffer := []teaList.Item{}
	currentMode := Following

	d.UpdateFunc = func(msg tea.Msg, m *teaList.Model) tea.Cmd {
		switch msg := msg.(type) {
		case Mode:
			safeMu.Lock()
			defer safeMu.Unlock()

			previousMode := currentMode
			currentMode = msg
			switch msg {
			case Following:
				m.Title = "Events"
				if currentMode != previousMode {
					m.ResetSelected()
					return m.SetItems(buffer)
				}
			case Frozen:
				m.Title = "Events (frozen)"
				return nil
			}

		case item:
			safeMu.Lock()
			defer safeMu.Unlock()

			switch currentMode {
			case Following:
				nb := len(m.Items())
				if nb >= maxItems {
					m.RemoveItem(nb - 1)
				}
				cmd := m.InsertItem(0, msg)
				buffer = m.Items()

				return safeUpdate(cmd)
			case Frozen:
				buffer = append([]teaList.Item{msg}, buffer...)
				if len(buffer) > maxItems {
					buffer = buffer[:len(buffer)-1]
				}
			}
			return nil

		case tea.WindowSizeMsg:
			// Enforce width here to avoid filter overflow
			m.SetWidth(msg.Width/2 - leftPanelStyle.GetHorizontalFrameSize())
			m.Help.Width = msg.Width
			return nil
		}

		return nil
	}

	return d
}
