package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			fmt.Println("Pressed key:", msg.String())
			return m, tea.Quit

		case "ctrl+n":
			// fmt.Println("User pressed Ctrl+N")
			m.createFileInputVisible = true

			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).Padding(0, 1, 0, 1)

	text := style.Render("Weclcome to GoPad! ")
	helpKeys := "Ctrl+N : New File | Ctrl+S : Save File | Ctrl+O : Open File | Ctrl+L : List Files | Ctrl+C or q : Quit"
	note := "Use the following keybindings to navigate and manage your files:"

	view := ""
	if m.createFileInputVisible {
		view += m.newFileInput.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n%s", text, view, note, helpKeys)
}

func initialModel() model {
	// initialize the text input component
	ti := textinput.New()
	ti.Placeholder = "What would you like to call your file?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return model{
		newFileInput:           ti,
		createFileInputVisible: false}
}

func main() {
	// fmt.Println("Welcome to GoPad!")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
		os.Exit(1)
	}
}
