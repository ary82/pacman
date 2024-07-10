package style

import "github.com/charmbracelet/lipgloss"

var (
	Block  = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))   // blue
	Pacman = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))  // yellow
	Blinky = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))   // red
	Pinky  = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))  // pink
	Inky   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))  // cyan
	Clyde  = lipgloss.NewStyle().Foreground(lipgloss.Color("202")) // orange
)
