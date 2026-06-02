package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"parthiv.cc/tui/internal/content"
	"parthiv.cc/tui/internal/markdown"
)

type Model struct {
	pages      []content.Page
	filtered   []int
	selected   int
	scroll     int
	width      int
	height     int
	query      string
	searchMode bool
	err        error
}

func New(pages []content.Page) Model {
	m := Model{
		pages:    pages,
		width:    96,
		height:   32,
		filtered: make([]int, len(pages)),
	}
	for i := range pages {
		m.filtered[i] = i
	}
	return m
}

func WithError(err error) Model {
	return Model{err: err, width: 96, height: 24}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.searchMode {
			return m.updateSearch(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.move(1)
		case "k", "up":
			m.move(-1)
		case "g", "home":
			m.selected = 0
			m.scroll = 0
		case "G", "end":
			m.selected = max(0, len(m.filtered)-1)
			m.scroll = 0
		case "pgdown", "ctrl+d":
			m.scroll += max(1, m.contentHeight()/2)
		case "pgup", "ctrl+u":
			m.scroll = max(0, m.scroll-max(1, m.contentHeight()/2))
		case "/":
			m.searchMode = true
			m.query = ""
		case "esc":
			m.query = ""
			m.applyFilter()
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v", m.err)) + "\n"
	}
	if len(m.filtered) == 0 {
		return m.layout(emptyStyle.Render("No pages match your search."))
	}

	return m.layout(m.renderContent())
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc":
		m.searchMode = false
	case "ctrl+c":
		return m, tea.Quit
	case "backspace":
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.query += string(msg.Runes)
		}
	}
	m.applyFilter()
	return m, nil
}

func (m *Model) move(delta int) {
	if len(m.filtered) == 0 {
		return
	}
	m.selected = clamp(m.selected+delta, 0, len(m.filtered)-1)
	m.scroll = 0
}

func (m *Model) applyFilter() {
	query := strings.ToLower(strings.TrimSpace(m.query))
	m.filtered = m.filtered[:0]
	for i, page := range m.pages {
		haystack := strings.ToLower(page.Title + "\n" + page.Description + "\n" + page.Section + "\n" + page.Body)
		if query == "" || strings.Contains(haystack, query) {
			m.filtered = append(m.filtered, i)
		}
	}
	m.selected = clamp(m.selected, 0, max(0, len(m.filtered)-1))
	m.scroll = 0
}

func (m Model) layout(content string) string {
	listWidth := 28
	if m.width < 80 {
		listWidth = 0
	}
	contentWidth := m.width - listWidth - 2
	if listWidth == 0 {
		contentWidth = m.width
	}
	contentWidth = max(24, contentWidth)

	left := ""
	if listWidth > 0 {
		left = listStyle.Width(listWidth).Height(m.contentHeight()).Render(m.renderList(listWidth))
	}

	right := contentStyle.Width(contentWidth).Height(m.contentHeight()).Render(content)
	body := right
	if left != "" {
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, m.footer())
}

func (m Model) renderList(width int) string {
	var b strings.Builder
	for visible, idx := range m.filtered {
		page := m.pages[idx]
		style := itemStyle
		prefix := "  "
		if visible == m.selected {
			style = selectedItemStyle
			prefix = "> "
		}
		line := lipgloss.NewStyle().MaxWidth(width - 2).Render(prefix + page.Title)
		b.WriteString(style.Render(line))
		if visible < len(m.filtered)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Model) renderContent() string {
	page := m.pages[m.filtered[m.selected]]
	contentWidth := m.width - 32
	if m.width < 80 {
		contentWidth = m.width - 2
	}
	rendered := markdown.New(contentWidth).Render(page.Body)
	lines := strings.Split(rendered, "\n")
	height := m.contentHeight()
	if m.scroll > max(0, len(lines)-height) {
		m.scroll = max(0, len(lines)-height)
	}
	end := min(len(lines), m.scroll+height)
	return strings.Join(lines[m.scroll:end], "\n")
}

func (m Model) footer() string {
	mode := ""
	if m.searchMode {
		mode = " /" + m.query
	} else if m.query != "" {
		mode = " filter:" + m.query
	}
	return footerStyle.Width(m.width).Render("j/k move  / search  pgup/pgdn scroll  esc clear  q quit" + mode)
}

func (m Model) contentHeight() int {
	return max(4, m.height-1)
}

func clamp(value, low, high int) int {
	return min(max(value, low), high)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	listStyle         = lipgloss.NewStyle().BorderRight(true).BorderForeground(lipgloss.Color("#545862")).PaddingRight(1)
	contentStyle      = lipgloss.NewStyle().PaddingLeft(2)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef")).Bold(true)
	footerStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#545862")).Background(lipgloss.Color("#282c34"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true)
	emptyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#565c64"))
)
