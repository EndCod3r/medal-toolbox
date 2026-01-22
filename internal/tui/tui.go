// internal/tui/tui.go
package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/EndCod3r/medal-toolbox/internal/clip"
	"github.com/EndCod3r/medal-toolbox/internal/operations"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	clips         map[string]clip.Clip
	filteredClips []clip.Clip
	jsonPath      string
	copyDir       string
	filterOptions clip.FilterOptions
	currentView   string
	list          list.Model
	textInput     textinput.Model
	message       string
	err           error
	width         int
	height        int
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("FFF")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)
	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel(jsonPath string) model {
	ti := textinput.New()
	ti.Placeholder = "Enter value"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	// Load clips
	clips, err := clip.LoadClipsFromFile(jsonPath)
	if err != nil {
		return model{err: err}
	}

	// Create menu items
	items := []list.Item{
		item{title: "Set JSON Path", desc: "Set the path to the clips JSON file"},
		item{title: "Set Copy Directory", desc: "Set the directory to copy clips to"},
		item{title: "Filter by Path", desc: "Filter clips by text in file path"},
		item{title: "Filter by Title", desc: "Filter clips by title"},
		item{title: "Filter by Collection ID", desc: "Filter clips by collection ID"},
		item{title: "Filter by Collection Name", desc: "Filter clips by collection name"},
		item{title: "Filter by Game", desc: "Filter clips by game name"},
		item{title: "Search Collections", desc: "Search for collections by name"},
		item{title: "Apply Filters", desc: "Apply current filters"},
		item{title: "Copy Clips", desc: "Copy filtered clips to destination"},
		item{title: "Quit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Medal Clip Manager"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)

	return model{
		clips:       clips,
		jsonPath:    jsonPath,
		currentView: "main",
		textInput:   ti,
		list:        l,
		filterOptions: clip.FilterOptions{
			PathSearch:     "",
			Title:          "",
			CollectionID:   "",
			CollectionName: "",
			Game:           "",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-6)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.currentView == "main" {
				// Handle menu selection
				selected := m.list.SelectedItem().(item).Title()
				switch selected {
				case "Set JSON Path":
					m.currentView = "json_path"
					m.textInput.SetValue(m.jsonPath)
				case "Set Copy Directory":
					m.currentView = "copy_dir"
					m.textInput.SetValue(m.copyDir)
				case "Filter by Path":
					m.currentView = "path_search"
					m.textInput.SetValue(m.filterOptions.PathSearch)
				case "Filter by Title":
					m.currentView = "title"
					m.textInput.SetValue(m.filterOptions.Title)
				case "Filter by Collection ID":
					m.currentView = "collection_id"
					m.textInput.SetValue(m.filterOptions.CollectionID)
				case "Filter by Collection Name":
					m.currentView = "collection_name"
					m.textInput.SetValue(m.filterOptions.CollectionName)
				case "Filter by Game":
					m.currentView = "game"
					m.textInput.SetValue(m.filterOptions.Game)
				case "Search Collections":
					m.currentView = "search_collections"
					m.textInput.SetValue("")
				case "Apply Filters":
					m.filteredClips = clip.FilterClips(m.clips, m.filterOptions)
					m.message = fmt.Sprintf("Found %d matching clips", len(m.filteredClips))
				case "Copy Clips":
					if m.copyDir == "" {
						m.message = "Error: Copy directory must be set first"
					} else if len(m.filteredClips) == 0 {
						m.message = "No clips to copy. Apply filters first."
					} else {
						result := operations.CopyClips(m.filteredClips, m.copyDir)
						m.message = fmt.Sprintf("Copied %d clips with %d errors", result.SuccessCount, result.ErrorCount)
						
						// Write log file
						if err := operations.WriteLogFile(m.copyDir, result); err != nil {
							m.message += fmt.Sprintf("\nError creating log: %v", err)
						} else {
							m.message += "\nLog file created in destination directory"
						}
					}
				case "Quit":
					return m, tea.Quit
				}
			} else {
				// Handle text input
				value := m.textInput.Value()
				switch m.currentView {
				case "json_path":
					m.jsonPath = value
					// Reload clips if path changed
					if newClips, err := clip.LoadClipsFromFile(m.jsonPath); err == nil {
						m.clips = newClips
						m.message = "Clips reloaded successfully"
					} else {
						m.message = fmt.Sprintf("Error loading clips: %v", err)
					}
				case "copy_dir":
					m.copyDir = value
				case "path_search":
					m.filterOptions.PathSearch = value
				case "title":
					m.filterOptions.Title = value
				case "collection_id":
					m.filterOptions.CollectionID = value
				case "collection_name":
					m.filterOptions.CollectionName = value
				case "game":
					m.filterOptions.Game = value
				case "search_collections":
					collections := clip.GetAllCollections(m.clips)
					matching := clip.SearchCollectionsByName(collections, value)
					if len(matching) > 0 {
						var collectionList []string
						for _, c := range matching {
							collectionList = append(collectionList, fmt.Sprintf("ID: %s, Name: %s", c.ID, c.Name))
						}
						m.message = "Collections found:\n" + strings.Join(collectionList, "\n")
					} else {
						m.message = "No collections found matching: " + value
					}
				}
				m.currentView = "main"
			}
		case "esc":
			if m.currentView != "main" {
				m.currentView = "main"
			}
		}
	}

	if m.currentView == "main" {
		m.list, cmd = m.list.Update(msg)
	} else {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	if m.currentView != "main" {
		return fmt.Sprintf(
			"%s\n\n%s",
			m.textInput.View(),
			"(press enter to confirm, esc to cancel)",
		) + "\n"
	}

	var view strings.Builder
	
	// Show current settings
	view.WriteString(titleStyle.Render("Current Settings"))
	view.WriteString(fmt.Sprintf("\nJSON Path: %s\n", m.jsonPath))
	view.WriteString(fmt.Sprintf("Copy Directory: %s\n", m.copyDir))
	view.WriteString(fmt.Sprintf("Path Filter: %s\n", m.filterOptions.PathSearch))
	view.WriteString(fmt.Sprintf("Title Filter: %s\n", m.filterOptions.Title))
	view.WriteString(fmt.Sprintf("Collection ID Filter: %s\n", m.filterOptions.CollectionID))
	view.WriteString(fmt.Sprintf("Collection Name Filter: %s\n", m.filterOptions.CollectionName))
	view.WriteString(fmt.Sprintf("Game Filter: %s\n\n", m.filterOptions.Game))
	
	if len(m.filteredClips) > 0 {
		view.WriteString(fmt.Sprintf("Filtered Clips: %d\n\n", len(m.filteredClips)))
		
		// Show a preview of the filtered clips with their game names
		view.WriteString("Preview of filtered clips:\n")
		for i, clip := range m.filteredClips {
			if i >= 5 { // Show only first 5 as preview
				view.WriteString(fmt.Sprintf("... and %d more\n", len(m.filteredClips)-5))
				break
			}
			gameName := clip.Content.Category.CategoryName
			if gameName == "" {
				gameName = "Unknown Game"
			}
			view.WriteString(fmt.Sprintf("  • %s (%s)\n", clip.GameTitle, gameName))
		}
		view.WriteString("\n")
	}
	
	if m.message != "" {
		if strings.HasPrefix(m.message, "Error:") {
			view.WriteString(errorStyle.Render(m.message) + "\n\n")
		} else {
			view.WriteString(messageStyle.Render(m.message) + "\n\n")
		}
		m.message = "" // Clear message after displaying
	}
	
	view.WriteString(m.list.View())
	
	return view.String()
}

func StartTUI(jsonPath string) {
	// Initialize the model
	m := initialModel(jsonPath)

	// Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting TUI: %v", err)
		os.Exit(1)
	}
}