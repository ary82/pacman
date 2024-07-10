package main

import (
	"fmt"
	"os"

	"github.com/ary82/pacman/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := game.InitialGameModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
