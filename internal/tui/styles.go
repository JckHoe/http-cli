package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	listStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(40)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	responseStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	statusErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	loadingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Blink(true)

	methodGETStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	methodPOSTStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	methodPUTStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true)

	methodDELETEStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	methodPATCHStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true)

	urlStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("117"))

	descriptionStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)
)

func getMethodStyle(method string) lipgloss.Style {
	switch method {
	case "GET":
		return methodGETStyle
	case "POST":
		return methodPOSTStyle
	case "PUT":
		return methodPUTStyle
	case "DELETE":
		return methodDELETEStyle
	case "PATCH":
		return methodPATCHStyle
	default:
		return lipgloss.NewStyle().Bold(true)
	}
}