package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error in geting home directory", err)
	}
	vaultDir = fmt.Sprintf("%s/.gopad", homeDir)
}

type model struct {
	newFileInput           textinput.Model
	currentFile            *os.File
	createFileInputVisible bool
	noteTextArea           textarea.Model
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
		case "ctrl+c":
			fmt.Println("Pressed key:", msg.String())
			return m, tea.Quit

		case "ctrl+n":
			// fmt.Println("User pressed Ctrl+N")
			m.createFileInputVisible = true

			return m, nil

		case "enter":
			if m.createFileInputVisible {
				filename := m.newFileInput.Value()
				if filename != "" {
					filePath := fmt.Sprintf("%s/%s.txt", vaultDir, filename)
					fileInfo, err := os.Stat(filePath)
					if err == nil && !fileInfo.IsDir() {
						fmt.Printf("File '%s' already exists!\n", filename)
						return m, nil
					} else {
						file, err := os.Create(filePath)

						if err != nil {
							log.Printf("Error creating file: %v\n", err)
						} else {
							fmt.Printf("File '%s' created successfully!\n", filename)
							m.currentFile = file
						}
					}

				}
				m.createFileInputVisible = false
				m.newFileInput.SetValue("") // Clear the input field after creating the file

			}
			return m, nil

		case "ctrl+s":
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				log.Printf("Error truncating file: %v\n", err)
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				log.Printf("Error seeking file: %v\n", err)
				return m, nil
			}

			if m.currentFile != nil {
				content := m.noteTextArea.Value()
				_, err := m.currentFile.WriteString(content)
				if err != nil {
					log.Printf("Error saving file: %v\n", err)
				} else {
					fmt.Printf("File '%s' saved successfully!\n", m.currentFile.Name())
				}
				m.currentFile.Close()
				m.currentFile = nil
				m.noteTextArea.SetValue("") // Clear the text area after saving the file
			}
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
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
	helpKeys := "Ctrl+N : New File | Ctrl+S : Save File | Ctrl+O : Open File | Ctrl+L : List Files | Ctrl+C : Quit"
	note := "Use the following keybindings to navigate and manage your files:"

	view := ""
	if m.createFileInputVisible {
		view += m.newFileInput.View()
	}

	if m.currentFile != nil {

		view = m.noteTextArea.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n%s", text, view, note, helpKeys)
}

func initialModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatalf("Error creating vault directory: %v\n", err)
		os.Exit(1)
	}

	// initialize the text input component
	ti := textinput.New()
	ti.Placeholder = "What would you like to call your file?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle

	// initialize the text area component
	ta := textarea.New()
	ta.Placeholder = "Start writing your notes here..."
	ta.Focus()
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false
	ta.Cursor.Style = cursorStyle

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta}
}

func main() {
	// fmt.Println("Welcome to GoPad!")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
		os.Exit(1)
	}
}
