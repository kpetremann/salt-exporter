package tui

import (
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	teaList "github.com/charmbracelet/bubbles/list"
	teaViewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kpetremann/salt-exporter/pkg/events"
)

const theme = "solarized-dark"

type format int

type Mode int

const (
	Following Mode = iota
	Frozen
)

type model struct {
	eventList      teaList.Model
	itemsBuffer    []teaList.Item
	rawView        teaViewport.Model
	eventChan      <-chan events.SaltEvent
	keys           *keyMap
	sideInfos      string
	terminalWidth  int
	terminalHeight int
	maxItems       int
	outputFormat   format
	currentMode    Mode
	wordWrap       bool
}

func NewModel(eventChan <-chan events.SaltEvent, maxItems int) model {
	var listKeys = defaultKeyMap()

	list := teaList.NewDefaultDelegate()

	selColor := lipgloss.Color("#fcc203")
	list.Styles.SelectedTitle = list.Styles.SelectedTitle.Foreground(selColor).BorderLeftForeground(selColor)
	list.Styles.SelectedDesc = list.Styles.SelectedTitle.Copy()

	eventList := teaList.New([]teaList.Item{}, list, 0, 0)
	eventList.Title = "Events"
	eventList.Styles.Title = listTitleStyle
	eventList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.enableFollow,
			listKeys.toggleWordwrap,
			listKeys.toggleJSONYAML,
		}
	}
	eventList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.enableFollow,
			listKeys.toggleJSONYAML,
		}
	}
	eventList.Filter = WordsFilter
	eventList.KeyMap = bubblesListKeyMap()

	rawView := teaViewport.New(1, 1)
	rawView.KeyMap = teaViewport.KeyMap{}

	return model{
		eventList:   eventList,
		rawView:     rawView,
		keys:        listKeys,
		eventChan:   eventChan,
		currentMode: Following,
		maxItems:    maxItems,
	}
}

func watchEvent(m model) tea.Cmd {
	return func() tea.Msg {
		e := <-m.eventChan
		var sender string = "master"
		if e.Data.Id != "" {
			sender = e.Data.Id
		}
		eventJSON, err := e.RawToJSON(true)
		if err != nil {
			log.Fatalln(err)
		}
		eventYAML, err := e.RawToYAML()
		if err != nil {
			log.Fatalln(err)
		}
		datetime, _ := time.Parse("2006-01-02T15:04:05.999999", e.Data.Timestamp)
		item := item{
			title:       e.Tag,
			description: e.Type,
			datetime:    datetime.Format("2006-01-02 15:04"),
			event:       e,
			sender:      sender,
			state:       e.ExtractState(),
			eventJSON:   string(eventJSON),
			eventYAML:   string(eventYAML),
		}

		return item
	}
}

func (m model) Init() tea.Cmd {
	return watchEvent(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Ensure the mode is Frozen if we are currently navigating
	if m.eventList.Index() > 0 {
		m.currentMode = Frozen
	}

	/*
		Manage events
	*/
	switch msg := msg.(type) {
	case item:
		switch m.currentMode {
		case Following:
			// In follow mode (default), we update both the list and the buffer
			currentList := m.eventList.Items()
			if len(currentList) >= m.maxItems {
				m.eventList.RemoveItem(len(currentList) - 1)
			}
			cmds = append(cmds, m.eventList.InsertItem(0, msg))
			m.itemsBuffer = m.eventList.Items()
		case Frozen:
			// In Frozen mode, we only update the buffer and keep the item list as is
			m.itemsBuffer = append([]teaList.Item{msg}, m.itemsBuffer...)
			if len(m.itemsBuffer) > m.maxItems {
				m.itemsBuffer = m.itemsBuffer[:len(m.itemsBuffer)-1]
			}
		}

		cmds = append(cmds, watchEvent(m))

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.eventList.FilterState() == teaList.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.enableFollow):
			m.currentMode = Following
			m.eventList.ResetSelected()
			cmds = append(cmds, m.eventList.SetItems(m.itemsBuffer))
		case key.Matches(msg, m.keys.toggleWordwrap):
			m.wordWrap = !m.wordWrap
		case key.Matches(msg, m.keys.toggleJSONYAML):
			m.outputFormat = (m.outputFormat + 1) % nbFormat
		}
	}

	/*
		Update embedded components
	*/
	var cmd tea.Cmd
	m.eventList, cmd = m.eventList.Update(msg)
	cmds = append(cmds, cmd)

	m.updateSideInfos()
	m.rawView, cmd = m.rawView.Update(msg)
	cmds = append(cmds, cmd)

	if m.eventList.Index() > 0 {
		m.currentMode = Frozen
	}

	m.updateTitle()
	return m, tea.Batch(cmds...)

}

func (m *model) updateSideInfos() {
	if sel := m.eventList.SelectedItem(); sel != nil {
		switch m.outputFormat {
		case YAML:
			m.sideInfos = sel.(item).eventYAML
			if m.wordWrap {
				m.sideInfos = strings.ReplaceAll(m.sideInfos, "\\n", "  \\\n")
			}
			if info, err := Highlight(m.sideInfos, "yaml", theme); err != nil {
				m.rawView.SetContent(m.sideInfos)
			} else {
				m.rawView.SetContent(info)
			}
		case JSON:
			m.sideInfos = sel.(item).eventJSON
			if m.wordWrap {
				m.sideInfos = strings.ReplaceAll(m.sideInfos, "\\n", "  \\\n")
			}
			if info, err := Highlight(m.sideInfos, "json", theme); err != nil {
				m.rawView.SetContent(m.sideInfos)
			} else {
				m.rawView.SetContent(info)
			}
		}
	}
}

func (m *model) updateTitle() {
	switch m.currentMode {
	case Following:
		m.eventList.Title = "Events"
	case Frozen:
		m.eventList.Title = "Events (frozen)"
	}
}

func (m model) View() string {
	/*
		Top bar
	*/
	topBarStyle.Width(m.terminalWidth)
	topBar := topBarStyle.Render(appTitleStyle.Render("Salt live"))

	var content []string
	contentHeight := m.terminalHeight - lipgloss.Height(topBar)
	contentWidth := m.terminalWidth / 2

	/*
		Left panel
	*/

	leftPanelStyle.Width(contentWidth)
	leftPanelStyle.Height(contentHeight)

	m.eventList.SetSize(
		contentWidth-leftPanelStyle.GetHorizontalFrameSize(),
		contentHeight-leftPanelStyle.GetVerticalFrameSize(),
	)

	content = append(content, leftPanelStyle.Render(m.eventList.View()))

	/*
		Right panel
	*/

	if m.sideInfos != "" {
		rawTitle := rightPanelTitleStyle.Render("Raw details")

		rightPanelStyle.Width(contentWidth)
		rightPanelStyle.Height(contentHeight)

		m.rawView.Width = contentWidth - rightPanelStyle.GetHorizontalFrameSize()
		m.rawView.Height = contentHeight - lipgloss.Height(rawTitle) - rightPanelStyle.GetVerticalFrameSize()

		sideInfos := rightPanelStyle.Render(lipgloss.JoinVertical(0, rawTitle, m.rawView.View()))
		content = append(content, sideInfos)
	}

	/*
		Final rendering
	*/
	return lipgloss.JoinVertical(0, topBar, lipgloss.JoinHorizontal(0, content...))
}
