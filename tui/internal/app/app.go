package app

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"parthiv.cc/tui/internal/content"
	"parthiv.cc/tui/internal/markdown"
	"parthiv.cc/tui/internal/site"
)

type Model struct {
	pages      []content.Page
	items      []treeItem
	filtered   []int
	selected   int
	scroll     int
	width      int
	height     int
	query      string
	searchMode bool
	err        error
}

func New(pages []content.Page, entries []site.TreeEntry) Model {
	items := buildTreeItems(pages, entries)
	m := Model{
		pages:    pages,
		items:    items,
		width:    96,
		height:   32,
		filtered: make([]int, 0, len(items)),
	}
	for i := range items {
		m.filtered = append(m.filtered, i)
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
			m.selected = m.firstSelectable()
			m.scroll = 0
		case "G", "end":
			m.selected = m.lastSelectable()
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
	next := m.selected + delta
	for next >= 0 && next < len(m.filtered) {
		if m.isSelectable(next) {
			m.selected = next
			m.scroll = 0
			return
		}
		next += delta
	}
	m.scroll = 0
}

func (m *Model) applyFilter() {
	query := strings.ToLower(strings.TrimSpace(m.query))
	m.filtered = m.filtered[:0]
	for i, item := range m.items {
		if query == "" {
			m.filtered = append(m.filtered, i)
			continue
		}

		if item.kind == itemDir {
			if strings.Contains(strings.ToLower(item.label), query) {
				m.filtered = append(m.filtered, i)
			}
			continue
		}
		if item.kind == itemSymlink {
			haystack := strings.ToLower(item.label + "\n" + item.target)
			if strings.Contains(haystack, query) {
				m.filtered = append(m.filtered, i)
			}
			continue
		}
		if item.kind == itemRoot && strings.Contains(strings.ToLower(item.label), query) {
			m.filtered = append(m.filtered, i)
			continue
		}

		page := m.pages[item.pageIndex]
		haystack := strings.ToLower(page.Title + "\n" + page.Description + "\n" + page.Section + "\n" + page.Body)
		if strings.Contains(haystack, query) || strings.Contains(strings.ToLower(item.label), query) {
			m.filtered = append(m.filtered, i)
		}
	}
	m.selected = clamp(m.selected, 0, max(0, len(m.filtered)-1))
	if len(m.filtered) > 0 && !m.isSelectable(m.selected) {
		m.selected = m.nextSelectable(m.selected, 1)
	}
	m.scroll = 0
}

func (m Model) layout(content string) string {
	listWidth := 34
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
		item := m.items[idx]
		prefix := strings.Repeat("  ", item.depth+1)
		if item.kind == itemRoot {
			prefix = "  "
		}
		line := lipgloss.NewStyle().MaxWidth(width - 2).Render(prefix + renderTreeLabel(item))
		if visible == m.selected && m.isSelectable(visible) {
			line = activeTreeItemStyle.Width(width - 1).Render(line)
		}
		b.WriteString(line)
		if visible < len(m.filtered)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Model) renderContent() string {
	if len(m.filtered) == 0 || !m.isSelectable(m.selected) {
		return emptyStyle.Render("No selectable page.")
	}
	item := m.items[m.filtered[m.selected]]
	if item.kind == itemSymlink {
		return bodyStyle.Render(item.label + " -> " + item.target)
	}
	page := m.pages[item.pageIndex]
	contentWidth := m.width - 38
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

func renderTreeLabel(item treeItem) string {
	switch item.kind {
	case itemRoot:
		return rootStyle.Render(item.label)
	case itemDir:
		return dirStyle.Render(item.label)
	case itemSymlink:
		return symlinkStyle.Render(item.label) + mutedStyle.Render(" -> ") + linkStyle.Render(item.target)
	default:
		return itemStyle.Render(item.label)
	}
}

func (m Model) contentHeight() int {
	return max(4, m.height-1)
}

func (m Model) isSelectable(filteredIndex int) bool {
	if filteredIndex < 0 || filteredIndex >= len(m.filtered) {
		return false
	}
	return m.items[m.filtered[filteredIndex]].kind != itemDir
}

func (m Model) firstSelectable() int {
	return m.nextSelectable(0, 1)
}

func (m Model) lastSelectable() int {
	return m.nextSelectable(len(m.filtered)-1, -1)
}

func (m Model) nextSelectable(start, delta int) int {
	if len(m.filtered) == 0 {
		return 0
	}
	start = clamp(start, 0, len(m.filtered)-1)
	for i := start; i >= 0 && i < len(m.filtered); i += delta {
		if m.isSelectable(i) {
			return i
		}
	}
	return 0
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
	listStyle           = lipgloss.NewStyle().BorderRight(true).BorderForeground(lipgloss.Color("#545862")).PaddingRight(1)
	contentStyle        = lipgloss.NewStyle().PaddingLeft(2)
	rootStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef")).Bold(true)
	bodyStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))
	itemStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))
	dirStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef"))
	symlinkStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#98c379"))
	linkStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef"))
	mutedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#545862"))
	activeTreeItemStyle = lipgloss.NewStyle().Background(lipgloss.Color("#3e4451"))
	footerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#545862")).Background(lipgloss.Color("#282c34"))
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true)
	emptyStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#565c64"))
)

type itemKind int

const (
	itemPage itemKind = iota
	itemRoot
	itemDir
	itemSymlink
)

type treeItem struct {
	kind      itemKind
	label     string
	sortKey   string
	target    string
	depth     int
	pageIndex int
}

func buildTreeItems(pages []content.Page, entries []site.TreeEntry) []treeItem {
	items := make([]treeItem, 0, len(pages)+4)
	pagesByURL := map[string]int{}
	for i, page := range pages {
		pagesByURL[page.URL()] = i
	}

	if rootIndex, ok := pagesByURL["/"]; ok {
		items = append(items, treeItem{kind: itemRoot, label: "~/parthiv.cc", sortKey: "", depth: 0, pageIndex: rootIndex})
	}

	if len(entries) > 0 {
		pathStack := map[int]string{}
		for _, entry := range entries {
			switch entry.Kind {
			case "dir":
				dirPath := entry.Name
				if entry.Depth > 0 && pathStack[entry.Depth-1] != "" {
					dirPath = pathStack[entry.Depth-1] + "/" + entry.Name
				}
				pathStack[entry.Depth] = dirPath
				pageIndex, ok := pagesByURL["/"+dirPath]
				if ok {
					items = append(items, treeItem{kind: itemPage, label: entry.Name + "/", sortKey: entry.Name, depth: entry.Depth, pageIndex: pageIndex})
					continue
				}
				items = append(items, treeItem{kind: itemDir, label: entry.Name + "/", sortKey: entry.Name, depth: entry.Depth})
			case "symlink":
				items = append(items, treeItem{kind: itemSymlink, label: entry.Name, sortKey: entry.Name, target: entry.Target, depth: entry.Depth})
			case "file":
				pageIndex, ok := pagesByURL[entry.URL]
				if ok {
					items = append(items, treeItem{kind: itemPage, label: entry.Name, sortKey: entry.Name, depth: entry.Depth, pageIndex: pageIndex})
				}
			}
		}
		return items
	}

	seenDirs := map[string]bool{}

	for i, page := range pages {
		if page.IsRoot() {
			continue
		}

		path := page.DisplayPath()
		parts := strings.Split(path, "/")
		if strings.HasSuffix(path, "/") {
			label := strings.TrimSuffix(path, "/") + "/"
			items = append(items, treeItem{kind: itemPage, label: label, sortKey: strings.TrimSuffix(path, "/"), depth: len(parts) - 2, pageIndex: i})
			continue
		}

		for depth := 0; depth < len(parts)-1; depth++ {
			dir := strings.Join(parts[:depth+1], "/")
			if seenDirs[dir] {
				continue
			}
			seenDirs[dir] = true
			items = append(items, treeItem{kind: itemDir, label: parts[depth] + "/", sortKey: dir, depth: depth})
		}
		items = append(items, treeItem{kind: itemPage, label: parts[len(parts)-1], sortKey: path, depth: len(parts) - 1, pageIndex: i})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].sortKey < items[j].sortKey
	})

	return items
}
