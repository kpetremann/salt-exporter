package tui

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	teaList "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	teaViewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k0kubun/pp/v3"
	"github.com/kpetremann/salt-exporter/pkg/event"
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
	sideBlock      teaViewport.Model
	demoText       textinput.Model
	eventChan      <-chan event.SaltEvent
	hardFilter     string
	keys           *keyMap
	sideInfos      string
	sideTitle      string
	terminalWidth  int
	terminalHeight int
	outputFormat   format
	currentMode    Mode
	wordWrap       bool
	demoMode       bool
	demoEnabled    bool
}

func NewModel(eventChan <-chan event.SaltEvent, maxItems int, filter string) model {
	var listKeys = defaultKeyMap()

	list := newDelegate(maxItems)

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
	eventList.SetShowHelp(false)
	eventList.SetShowTitle(false)
	eventList.Filter = WordsFilter
	eventList.KeyMap = bubblesListKeyMap()

	rawView := teaViewport.New(1, 1)
	rawView.KeyMap = teaViewport.KeyMap{}

	m := model{
		eventList:   eventList,
		sideBlock:   rawView,
		keys:        listKeys,
		eventChan:   eventChan,
		hardFilter:  filter,
		currentMode: Following,
	}

	if os.Getenv("SALT_DEMO") == "true" {
		m.demoEnabled = true
		m.demoText = textinput.New()
		m.demoText.Focus()
	}

	return m
}

func watchEvent(m model) tea.Cmd {
	return func() tea.Msg {
		for {
			e := <-m.eventChan

			sender := "master"
			if e.Data.ID != "" {
				sender = e.Data.ID
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
			i := item{
				title:       e.Tag,
				description: e.Type,
				datetime:    datetime.Format("15:04"),
				event:       e,
				sender:      sender,
				state:       e.ExtractState(),
				eventJSON:   string(eventJSON),
				eventYAML:   string(eventYAML),
			}

			// No hard filter set
			if m.hardFilter == "" {
				return i
			}

			// Hard filter set
			if rank := m.eventList.Filter(m.hardFilter, []string{i.FilterValue()}); len(rank) > 0 {
				return i
			}
		}
	}
}

func (m model) Init() tea.Cmd {
	return watchEvent(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.demoMode {
		var cmd tea.Cmd
		m.demoText, cmd = m.demoText.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Ensure the mode is Frozen if we are currently navigating
	if m.eventList.Index() > 0 {
		m.currentMode = Frozen
	}

	/*
		Manage events
	*/
	switch msg := msg.(type) {
	case item:
		cmds = append(cmds, watchEvent(m))

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.eventList.FilterState() == teaList.Filtering {
			break
		}

		if m.demoEnabled && key.Matches(msg, m.keys.demoText) {
			m.demoMode = !m.demoMode
			m.demoText.SetValue("")
		}
		if m.demoMode {
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, m.keys.enableFollow):
			var cmd tea.Cmd
			m.currentMode = Following
			m.eventList, cmd = m.eventList.Update(Following)
			return m, cmd
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
	m.sideBlock, cmd = m.sideBlock.Update(msg)
	cmds = append(cmds, cmd)

	if m.eventList.Index() > 0 {
		var cmd tea.Cmd
		m.currentMode = Frozen
		m.eventList, cmd = m.eventList.Update(Frozen)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateSideInfos() {
	if sel := m.eventList.SelectedItem(); sel != nil {
		switch m.outputFormat {
		case YAML:
			m.sideTitle = "Raw event (YAML)"
			m.sideInfos = sel.(item).eventYAML
			if m.wordWrap {
				m.sideInfos = strings.ReplaceAll(m.sideInfos, "\\n", "  \\\n")
			}
			if info, err := Highlight(m.sideInfos, "yaml", theme); err != nil {
				m.sideBlock.SetContent(m.sideInfos)
			} else {
				m.sideBlock.SetContent(info)
			}
		case JSON:
			m.sideTitle = "Raw event (JSON)"
			m.sideInfos = sel.(item).eventJSON
			if m.wordWrap {
				m.sideInfos = strings.ReplaceAll(m.sideInfos, "\\n", "  \\\n")
			}
			if info, err := Highlight(m.sideInfos, "json", theme); err != nil {
				m.sideBlock.SetContent(m.sideInfos)
			} else {
				m.sideBlock.SetContent(info)
			}
		case PARSED:
			m.sideTitle = "Parsed event (Golang)"
			eventLite := sel.(item).event
			eventLite.RawBody = nil
			m.sideInfos = pp.Sprint(eventLite)
			m.sideBlock.SetContent(m.sideInfos)
		}
	}
}

func (m model) View() string {
	if m.demoMode {
		return lipgloss.Place(m.terminalWidth, m.terminalHeight, lipgloss.Center, lipgloss.Center, m.demoText.View())
	}

	/*
		Bottom
	*/
	helpView := m.eventList.Help.View(m.eventList)

	/*
		Top bar
	*/
	topBarStyle.Width(m.terminalWidth)
	topBar := topBarStyle.Render(appTitleStyle.Render("Salt Live"))

	// Calculate content height for left and right panels
	var content []string
	contentHeight := m.terminalHeight - lipgloss.Height(topBar) - lipgloss.Height(helpView)
	contentWidth := m.terminalWidth / 2

	/*
		Left panel
	*/

	if m.currentMode == Frozen {
		listTitleStyle.Background(lipgloss.Color("#a02725"))
		listTitleStyle.Foreground(lipgloss.Color("#ffffff"))
	} else {
		listTitleStyle.UnsetBackground()
		listTitleStyle.UnsetForeground()
	}
	listTitle := listTitleStyle.Render(m.eventList.Title)

	leftPanelStyle.Width(contentWidth)
	leftPanelStyle.Height(contentHeight)

	m.eventList.SetSize(
		contentWidth-leftPanelStyle.GetHorizontalFrameSize(),
		contentHeight-lipgloss.Height(listTitle)-leftPanelStyle.GetVerticalFrameSize(),
	)

	listWithTitle := lipgloss.JoinVertical(0, listTitle, m.eventList.View())

	content = append(content, leftPanelStyle.Render(listWithTitle))

	/*
		Right panel
	*/

	if m.sideInfos != "" {
		rawTitle := rightPanelTitleStyle.Render(m.sideTitle)

		rightPanelStyle.Width(contentWidth)
		rightPanelStyle.Height(contentHeight)

		m.sideBlock.Width = contentWidth - rightPanelStyle.GetHorizontalFrameSize()
		m.sideBlock.Height = contentHeight - lipgloss.Height(rawTitle) - rightPanelStyle.GetVerticalFrameSize()

		sideInfos := rightPanelStyle.Render(lipgloss.JoinVertical(0, rawTitle, m.sideBlock.View()))
		content = append(content, sideInfos)
	}

	/*
		Final rendering
	*/
	return lipgloss.JoinVertical(0, topBar, lipgloss.JoinHorizontal(0, content...), helpView)
}
