package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cassielabs/cassie-cli/internal/executor"
	"github.com/cassielabs/cassie-cli/internal/parser"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateFileList state = iota
	stateRequestList
	stateResponse
)

type model struct {
	state        state
	files        []string
	fileIndex    int
	httpFile     *parser.HTTPFile
	requests     []parser.HTTPRequest
	requestIndex int
	response     *executor.Response
	loading      bool
	err          error
	width        int
	height       int
	viewport     viewport.Model
	exec         *executor.Executor
	filePath     string
}

type responseMsg struct {
	response *executor.Response
	err      error
}

func initialModel(filePath string, timeout time.Duration) model {
	m := model{
		state:    stateFileList,
		exec:     executor.New(timeout),
		filePath: filePath,
		viewport: viewport.New(80, 20),
	}

	if filePath != "" {
		httpFile, err := parser.ParseFile(filePath)
		if err == nil {
			m.httpFile = httpFile
			m.requests = httpFile.Requests
			m.state = stateRequestList
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	if m.filePath == "" {
		return m.loadFiles()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 8
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			switch m.state {
			case stateFileList:
				if m.fileIndex > 0 {
					m.fileIndex--
				}
			case stateRequestList:
				if m.requestIndex > 0 {
					m.requestIndex--
				}
			case stateResponse:
				m.viewport.LineUp(1)
			}

		case key.Matches(msg, keys.Down):
			switch m.state {
			case stateFileList:
				if m.fileIndex < len(m.files)-1 {
					m.fileIndex++
				}
			case stateRequestList:
				if m.requestIndex < len(m.requests)-1 {
					m.requestIndex++
				}
			case stateResponse:
				m.viewport.LineDown(1)
			}

		case key.Matches(msg, keys.Enter):
			switch m.state {
			case stateFileList:
				if len(m.files) > 0 {
					return m, m.loadFile(m.files[m.fileIndex])
				}
			case stateRequestList:
				if len(m.requests) > 0 {
					m.loading = true
					m.state = stateResponse
					return m, m.executeRequest(m.requests[m.requestIndex])
				}
			}

		case key.Matches(msg, keys.Back):
			switch m.state {
			case stateRequestList:
				if m.filePath == "" {
					m.state = stateFileList
				}
			case stateResponse:
				m.state = stateRequestList
				m.response = nil
			}

		case key.Matches(msg, keys.Refresh):
			if m.state == stateFileList {
				return m, m.loadFiles()
			}
		}

	case []string:
		m.files = msg
		if len(m.files) > 0 && m.fileIndex >= len(m.files) {
			m.fileIndex = len(m.files) - 1
		}

	case *parser.HTTPFile:
		m.httpFile = msg
		m.requests = msg.Requests
		m.state = stateRequestList
		m.requestIndex = 0

	case responseMsg:
		m.loading = false
		m.response = msg.response
		m.err = msg.err
		if m.response != nil {
			m.viewport.SetContent(executor.FormatResponse(m.response))
		}

	case error:
		m.err = msg
		m.loading = false
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.state {
	case stateFileList:
		return m.renderFileList()
	case stateRequestList:
		return m.renderRequestList()
	case stateResponse:
		return m.renderResponse()
	default:
		return ""
	}
}

func (m model) renderFileList() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Select HTTP File") + "\n\n")

	if len(m.files) == 0 {
		b.WriteString("No .http files found in current directory\n")
	} else {
		items := make([]string, len(m.files))
		for i, file := range m.files {
			if i == m.fileIndex {
				items[i] = selectedItemStyle.Render("→ " + file)
			} else {
				items[i] = normalItemStyle.Render("  " + file)
			}
		}
		content := strings.Join(items, "\n")
		b.WriteString(listStyle.Width(m.width - 4).Render(content))
	}

	b.WriteString("\n\n" + helpStyle.Render("↑/↓: navigate • enter: select • q: quit • r: refresh"))
	return b.String()
}

func (m model) renderRequestList() string {
	var b strings.Builder
	
	title := "HTTP Requests"
	if m.httpFile != nil && m.httpFile.Path != "" {
		title = fmt.Sprintf("HTTP Requests - %s", filepath.Base(m.httpFile.Path))
	}
	b.WriteString(titleStyle.Render(title) + "\n\n")

	if len(m.requests) == 0 {
		b.WriteString("No requests found in file\n")
	} else {
		items := make([]string, len(m.requests))
		for i, req := range m.requests {
			method := getMethodStyle(req.Method).Render(req.Method)
			url := urlStyle.Render(req.URL)
			
			line := fmt.Sprintf("%s %s", method, url)
			if req.Name != "" {
				line = fmt.Sprintf("[%s] %s", req.Name, line)
			}
			
			if i == m.requestIndex {
				items[i] = selectedItemStyle.Render("→ ") + line
			} else {
				items[i] = normalItemStyle.Render("  ") + line
			}
		}
		content := strings.Join(items, "\n")
		b.WriteString(listStyle.Width(m.width - 4).Render(content))
	}

	help := "↑/↓: navigate • enter: execute • q: quit"
	if m.filePath == "" {
		help += " • esc: back to files"
	}
	b.WriteString("\n\n" + helpStyle.Render(help))
	return b.String()
}

func (m model) renderResponse() string {
	var b strings.Builder
	
	req := m.requests[m.requestIndex]
	title := fmt.Sprintf("%s %s", req.Method, req.URL)
	if req.Name != "" {
		title = fmt.Sprintf("[%s] %s", req.Name, title)
	}
	b.WriteString(titleStyle.Render(title) + "\n\n")

	if m.loading {
		b.WriteString(loadingStyle.Render("Executing request..."))
	} else if m.err != nil {
		b.WriteString(statusErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	} else if m.response != nil {
		content := m.viewport.View()
		b.WriteString(responseStyle.Width(m.width - 4).Height(m.height - 6).Render(content))
	}

	b.WriteString("\n\n" + helpStyle.Render("↑/↓: scroll • esc: back to requests • q: quit"))
	return b.String()
}

func (m model) loadFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := filepath.Glob("*.http")
		if err != nil {
			return err
		}
		return files
	}
}

func (m model) loadFile(path string) tea.Cmd {
	return func() tea.Msg {
		httpFile, err := parser.ParseFile(path)
		if err != nil {
			return err
		}
		
		for key, value := range httpFile.Variables {
			if envValue := os.Getenv(key); envValue != "" {
				httpFile.Variables[key] = envValue
			} else {
				httpFile.Variables[key] = value
			}
		}
		
		return httpFile
	}
}

func (m model) executeRequest(req parser.HTTPRequest) tea.Cmd {
	return func() tea.Msg {
		if m.httpFile != nil {
			req.ApplyVariables(m.httpFile.Variables)
		}
		
		resp, err := m.exec.Execute(req)
		return responseMsg{
			response: resp,
			err:      err,
		}
	}
}

func Run(filePath string, timeout time.Duration) error {
	p := tea.NewProgram(
		initialModel(filePath, timeout),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}