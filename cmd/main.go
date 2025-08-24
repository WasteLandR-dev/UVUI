// cmd/main.go
package main

import (
	"log"

	"uvui/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := app.NewModel()

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
