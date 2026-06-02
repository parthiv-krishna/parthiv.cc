package markdown

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var headingRE = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*#*\s*$`)

type Renderer struct {
	width int
}

func New(width int) Renderer {
	if width < 24 {
		width = 24
	}
	return Renderer{width: width}
}

func (r Renderer) Render(src string) string {
	lines := strings.Split(src, "\n")
	out := make([]string, 0, len(lines))
	paragraph := make([]string, 0, 4)
	inFence := false

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		text := strings.Join(paragraph, " ")
		out = append(out, bodyStyle.Width(r.width).Render(text))
		paragraph = paragraph[:0]
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			flushParagraph()
			inFence = !inFence
			out = append(out, fenceStyle.Render(line))
			continue
		}

		if inFence {
			out = append(out, codeStyle.Render(line))
			continue
		}

		if trimmed == "" {
			flushParagraph()
			if len(out) > 0 && out[len(out)-1] != "" {
				out = append(out, "")
			}
			continue
		}

		if match := headingRE.FindStringSubmatch(trimmed); match != nil {
			flushParagraph()
			out = append(out, renderHeading(len(match[1]), strings.TrimSpace(match[2])))
			continue
		}

		if strings.HasPrefix(trimmed, ">") {
			flushParagraph()
			out = append(out, quoteStyle.Width(r.width).Render(trimmed))
			continue
		}

		if isListLine(trimmed) {
			flushParagraph()
			out = append(out, bodyStyle.Width(r.width).Render(trimmed))
			continue
		}

		paragraph = append(paragraph, trimmed)
	}
	flushParagraph()

	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

func renderHeading(level int, text string) string {
	marker := strings.Repeat("#", level) + " "
	return markerStyle.Render(marker) + headingStyle(level).Render(text)
}

func headingStyle(level int) lipgloss.Style {
	switch level {
	case 1:
		return h1Style
	case 2:
		return h2Style
	case 3:
		return h3Style
	case 4:
		return h4Style
	case 5:
		return h5Style
	default:
		return h6Style
	}
}

func isListLine(line string) bool {
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
		return true
	}
	for i, r := range line {
		if r == '.' && i > 0 && i+1 < len(line) && line[i+1] == ' ' {
			for _, prefix := range line[:i] {
				if prefix < '0' || prefix > '9' {
					return false
				}
			}
			return true
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return false
}

var (
	bodyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))
	markerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#545862"))
	h1Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef")).Bold(true)
	h2Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#c678dd")).Bold(true)
	h3Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#56b6c2")).Bold(true)
	h4Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#98c379")).Bold(true)
	h5Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#e5c07b")).Bold(true)
	h6Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true)
	quoteStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#56b6c2")).BorderLeft(true).BorderForeground(lipgloss.Color("#545862")).PaddingLeft(1)
	fenceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#545862"))
	codeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf")).Background(lipgloss.Color("#353b45"))
)
