package tui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle()

	appTitleStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#255aa0")).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#255aa0")).
			BorderBackground(lipgloss.Color("#255aa0"))

	topBarStyle = lipgloss.NewStyle().Padding(0, 1).
			BorderStyle(lipgloss.InnerHalfBlockBorder()).
			BorderBottom(true).BorderForeground(lipgloss.Color("#255aa0"))

	listTitleStyle = lipgloss.NewStyle().Padding(0, 1).
			Border(lipgloss.RoundedBorder())

	leftPanelStyle = lipgloss.NewStyle().Padding(0, 2)

	rightPanelTitleStyle = lipgloss.NewStyle().Padding(0, 1).MarginBottom(1).
				Border(lipgloss.RoundedBorder())

	rightPanelStyle = lipgloss.NewStyle().
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderForeground(lipgloss.Color("#255aa0"))
)
