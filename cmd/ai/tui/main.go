package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/luvx21/coding-go/cmds/ai/utils"
	"github.com/luvx21/coding-go/infra/ai"
)

var (
	curModel  *ai.Model
	program   *tea.Program
	render, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
)

type tui struct {
	textarea  textarea.Model
	markdown  *strings.Builder
	isLoading bool
}

func initialModel() *tui {
	ti := textarea.New()
	ti.Placeholder = "Ask AI..."
	ti.Focus()
	// ti.ShowLineNumbers = false
	ti.SetHeight(6)
	ti.SetWidth(120)

	return &tui{
		textarea:  ti,
		markdown:  &strings.Builder{},
		isLoading: false,
	}
}

func (m *tui) Init() tea.Cmd {
	return textarea.Blink
}

func (m *tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case tea.KeyEsc.String():
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlN.String():
			if m.textarea.Focused() {
				m.textarea.InsertString("\n")
			}
			return m, nil
		case tea.KeyEnter.String():
			query := m.textarea.Value()
			if !m.isLoading && len(query) > 0 {
				m.textarea.Reset()
				m.markdown.Reset()
				m.isLoading = true

				go sendOpenAIRequest(m, query)
				return m, nil
			}
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				m.textarea.Focus()
			}
		}

	case string:
		m.markdown.WriteString(msg)
		return m, nil

	case error:
		m.markdown.WriteString(fmt.Sprintf("Error: %v", msg))
		m.isLoading = false
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *tui) View() string {
	var sb strings.Builder
	sb.WriteString("Ask AI(Ctrl+N:换行 Ctrl+C:退出):\n")
	sb.WriteString(m.textarea.View() + "\n\n")

	if m.isLoading {
		statusLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Render("│ AI正在思考...")
		sb.WriteString(statusLine + "\n")
	} else {
		content := m.markdown.String()
		if len(content) > 0 {
			sb.WriteString(renderMarkdown(content + "\n\n"))
		}
	}
	return sb.String()
}

func sendOpenAIRequest(m *tui, query string) error {
	res, err := curModel.Request(query, true)
	if err != nil {
		fmt.Printf("请求接口失败,模型:%s, 服务商:%s\n", curModel.Id, curModel.Sp.BaseUrl)
		return err
	}
	defer res.Body.Close()

	m.isLoading = false

	// go program.Send("\n### send内容1\n\n\n### " + query)
	// go program.Send("\n### send内容2\n\n\n### " + query)

	program.Send(fmt.Sprintf("## %s(%s) %s\n", query, curModel.Id, time.Now().Format("2006-01-02 15:04:05")))

	scanner := bufio.NewReader(res.Body)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			break
		}

		content := ai.ParseLineContent(line)
		if content != "" {
			program.Send(content)
		}
	}
	program.Send("\n---\n")
	return nil
}

func renderMarkdown(text string) string {
	if text == "" {
		return ""
	}
	out, _ := render.Render(text)
	return out
}

func main() {
	curModel = utils.SelectModel(curModel)
	p := tea.NewProgram(initialModel())
	program = p
	if _, err := p.Run(); err != nil {
		fmt.Printf("程序错误: %v", err)
		os.Exit(1)
	}
}
