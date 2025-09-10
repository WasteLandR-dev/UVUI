// cmd/main.go
package main

import (
	"log"

	"uvui/internal/app"
	"uvui/internal/services"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	executor := services.NewCommandExecutor()
	uvInstaller := services.NewUVInstaller(executor)
	pythonManager := services.NewPythonManager(executor)
	projectManager := services.NewProjectManager(executor)

	model := app.NewModel(uvInstaller, pythonManager, projectManager, executor)

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
