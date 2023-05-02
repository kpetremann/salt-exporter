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

type listKeyMap struct {
	enableFollow   key.Binding
	toggleJSONYAML key.Binding
	toggleWordwrap key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		enableFollow: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "follow mode"),
		),
		toggleWordwrap: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "toggle JSON word wrap"),
		),
		toggleJSONYAML: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "toggle JSON/YAML"),
		),
	}
}

type model struct {
	eventList      teaList.Model
	itemsBuffer    []teaList.Item
	rawView        teaViewport.Model
	eventChan      <-chan events.SaltEvent
	keys           *listKeyMap
	sideInfos      string
	terminalWidth  int
	terminalHeight int
	maxItems       int
	outputFormat   format
	followMode     bool
	jsonWordwrap   bool
}

func NewModel(eventChan <-chan events.SaltEvent, maxItems int) model {
	var listKeys = newListKeyMap()

	list := teaList.NewDefaultDelegate()

	selColor := lipgloss.Color("#fcc203")
	list.Styles.SelectedTitle = list.Styles.SelectedTitle.Foreground(selColor).BorderLeftForeground(selColor)
	list.Styles.SelectedDesc = list.Styles.SelectedTitle.Copy()

	eventList := teaList.New([]teaList.Item{}, list, 0, 0)
	eventList.Title = "Events"
	eventList.Styles.TitleBar = listTitleStyle
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

	rawView := teaViewport.New(1, 1)
	rawView.KeyMap = teaViewport.KeyMap{}

	return model{
		eventList:  eventList,
		rawView:    rawView,
		keys:       listKeys,
		eventChan:  eventChan,
		followMode: true,
		maxItems:   maxItems,
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
	return tea.Batch(
		watchEvent(m),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.eventList.Index() > 0 {
		m.followMode = false
	}

	if m.followMode {
		m.eventList.Styles.Title = listTitleStyle
		m.eventList.Styles.TitleBar = lipgloss.NewStyle()
		cmds = append(cmds, m.eventList.NewStatusMessage(""))
	} else {
		m.eventList.Styles.TitleBar = listTitleStyle
		m.eventList.Styles.Title = lipgloss.NewStyle()
		cmds = append(cmds, m.eventList.NewStatusMessage(lipgloss.NewStyle().Italic(true).Render("frozen")))
	}

	switch msg := msg.(type) {
	case item:
		m.itemsBuffer = append([]teaList.Item{msg}, m.itemsBuffer...)
		if len(m.itemsBuffer) > m.maxItems {
			m.itemsBuffer = m.itemsBuffer[:len(m.itemsBuffer)-1]
		}

		// When not in follow mode, we freeze the visible list.
		if m.followMode {
			cmds = append(cmds, m.eventList.SetItems(m.itemsBuffer))
		}
		cmds = append(cmds, watchEvent(m))

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.terminalWidth = msg.Width - h*2
		m.terminalHeight = msg.Height - v*2

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.enableFollow):
			m.followMode = true
			m.eventList.ResetSelected()
			return m, nil
		case key.Matches(msg, m.keys.toggleWordwrap):
			m.jsonWordwrap = !m.jsonWordwrap
		case key.Matches(msg, m.keys.toggleJSONYAML):
			m.outputFormat = (m.outputFormat + 1) % nbFormat
		}
	}

	var cmd tea.Cmd
	m.eventList, cmd = m.eventList.Update(msg)
	cmds = append(cmds, cmd)

	if sel := m.eventList.SelectedItem(); sel != nil {
		switch m.outputFormat {
		case YAML:
			m.sideInfos = sel.(item).eventYAML
			if info, err := Highlight(m.sideInfos, "yaml", theme); err != nil {
				m.rawView.SetContent(m.sideInfos)
			} else {
				m.rawView.SetContent(info)
			}
		case JSON:
			m.sideInfos = sel.(item).eventJSON
			if m.jsonWordwrap {
				m.sideInfos = strings.ReplaceAll(m.sideInfos, "\\n", "  \\\n")
			}
			if info, err := Highlight(m.sideInfos, "json", theme); err != nil {
				m.rawView.SetContent(m.sideInfos)
			} else {
				m.rawView.SetContent(info)
			}
		}
	}

	m.rawView, cmd = m.rawView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	topBarStyle.Width(m.terminalWidth)
	topBar := topBarStyle.Render(appTitleStyle.Render("Salt live"))

	var content []string
	contentHeight := m.terminalHeight - lipgloss.Height(topBar)

	m.eventList.SetSize(m.terminalWidth, contentHeight)
	leftPanelStyle.Width(m.terminalWidth / 2)
	leftPanelStyle.Height(contentHeight)

	content = append(content, leftPanelStyle.Render(m.eventList.View()))
	if m.sideInfos != "" {
		rawContent := rightPanelTitleStyle.Render("Raw details")

		m.rawView.Width = m.terminalWidth / 2
		m.rawView.Height = contentHeight - lipgloss.Height(rawContent)

		rightPanelStyle.Width(m.terminalWidth / 2)
		rightPanelStyle.Height(contentHeight)

		sideInfos := rightPanelStyle.Render(lipgloss.JoinVertical(0, rawContent, m.rawView.View()))
		content = append(content, sideInfos)
	}

	app := lipgloss.JoinVertical(0, topBar, lipgloss.JoinHorizontal(0, content...))

	return appStyle.Render(app)
}
