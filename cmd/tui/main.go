package main

import (
	"fmt"
	"os"

	"github.com/ary82/pacman/internal/game"
	"github.com/ary82/pacman/internal/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	model := game.InitialGameModel()
	model.Styles = style.Styles{
		Block:  lipgloss.NewStyle().Foreground(lipgloss.Color("4")),   // blue
		Pacman: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),  // yellow
		Blinky: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),   // red
		Pinky:  lipgloss.NewStyle().Foreground(lipgloss.Color("13")),  // pink
		Inky:   lipgloss.NewStyle().Foreground(lipgloss.Color("14")),  // cyan
		Clyde:  lipgloss.NewStyle().Foreground(lipgloss.Color("202")), // orange
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
