package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"parthiv.cc/tui/internal/app"
	"parthiv.cc/tui/internal/content"
	"parthiv.cc/tui/internal/site"
)

func main() {
	configPath := flag.String("config", "config.toml", "Zola config file")
	contentRoot := flag.String("content", "content", "Zola content directory")
	flag.Parse()

	pages, err := content.Load(*contentRoot)
	entries, configErr := site.LoadTree(*configPath)
	model := app.New(pages, entries)
	if err != nil || configErr != nil {
		err = errors.Join(err, configErr)
		model = app.WithError(err)
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tui failed: %v\n", err)
		os.Exit(1)
	}
}
