package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	board  [31][28]int
	lives  int
	score  int
	pacman pacman
	ghosts [4]ghost
}

type pacman struct {
	xPosition int
	yPosition int
	direction int // 0: static, 1: up, 2: down, 3: left 4: right
}

type ghost struct {
	xPosition int
	yPosition int
}

type updatePacmanPosition int

const (
	blockStr  string = "███"
	pacmanStr string = " 󰮯 "
	ghostStr  string = " 󰊠 "
	emptyStr  string = "   "
	pelletStr string = " 🞄 "
)

var (
	pacmanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	blockStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	blinkyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	pinkyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	inkyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	clydeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
)

func initialModel() model {
	return model{
		board: [31][28]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 3, 3, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 3, 3, 3, 3, 3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 3, 3, 3, 3, 3, 3, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 3, 3, 3, 3, 3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0},
			{0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0},
			{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0},
			{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
			{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		pacman: pacman{
			xPosition: 1,
			yPosition: 1,
			direction: 0,
		},
		score: 0,
		lives: 3,
		ghosts: [4]ghost{
			{xPosition: 15, yPosition: 12},
			{xPosition: 15, yPosition: 13},
			{xPosition: 15, yPosition: 14},
			{xPosition: 15, yPosition: 15},
		},
	}
}

func movePacman() tea.Cmd {
	return tea.Every(50*time.Millisecond, func(t time.Time) tea.Msg {
		return updatePacmanPosition(0)
	})
}

func (m model) Init() tea.Cmd {
	return movePacman()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		case "j", "down", "s":
			m.pacman.direction = 2
		case "k", "up", "w":
			m.pacman.direction = 1
		case "h", "left", "a":
			m.pacman.direction = 3
		case "l", "right", "d":
			m.pacman.direction = 4

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case updatePacmanPosition:
		prevX := m.pacman.xPosition
		prevY := m.pacman.yPosition
		newX := prevX
		newY := prevY

		switch m.pacman.direction {
		case 0:
		case 1:
			newX -= 1
		case 2:
			newX += 1
		case 3:
			newY -= 1
		case 4:
			newY += 1
		}

		if newX == 14 && newY == -1 {
			newY = 27
		}
		if newX == 14 && newY == 28 {
			newY = 0
		}

		switch newPosition := m.board[newX][newY]; newPosition {
		case 1:
			m.score += 1
			m.pacman.xPosition = newX
			m.pacman.yPosition = newY
			m.board[prevX][prevY] = 3
			m.board[newX][newY] = 2
		case 3:
			m.pacman.xPosition = newX
			m.pacman.yPosition = newY
			m.board[prevX][prevY] = 3
			m.board[newX][newY] = 2
		}
		return m, movePacman()
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for _, v := range m.board {
		for _, num := range v {
			switch num {
			case 0:
				s += blockStyle.Render(blockStr)
			case 1:
				s += pelletStr
			case 2:
				s += pacmanStyle.Render(pacmanStr)
			case 3:
				s += emptyStr
			}
		}
		s += "\n"
	}
	s += fmt.Sprintf("direction: %v\n", m.pacman.direction)
	s += fmt.Sprintf("score: %v\n", m.score)
	s += fmt.Sprintf("xPosition: %v\n", m.pacman.xPosition)
	s += fmt.Sprintf("yPosition: %v\n", m.pacman.yPosition)

	// Send the UI for rendering
	return s
}

func main() {
	model := initialModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
