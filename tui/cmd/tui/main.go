package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"parthiv.cc/tui/internal/app"
	"parthiv.cc/tui/internal/content"
)

func main() {
	contentRoot := flag.String("content", "content", "Zola content directory")
	flag.Parse()

	pages, err := content.Load(*contentRoot)
	model := app.New(pages)
	if err != nil {
		model = app.WithError(err)
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tui failed: %v\n", err)
		os.Exit(1)
	}
}
