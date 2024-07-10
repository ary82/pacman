package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ary82/pacman/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	board  [31][28]int
	lives  int
	score  int
	pacman pacman
	ghosts [4]*ghost
}

type pacman struct {
	xPosition int
	yPosition int
	direction int // 0: static, 1: up, 2: down, 3: left 4: right
}

type ghost struct {
	xPosition               int
	yPosition               int
	viewCode                int
	isCurrentPositionPallet bool
}

type (
	updatePacmanPosition int
	updateGhostsPosition int
)

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
		ghosts: [4]*ghost{
			{xPosition: 15, yPosition: 12, viewCode: 4},
			{xPosition: 15, yPosition: 13, viewCode: 5},
			{xPosition: 15, yPosition: 14, viewCode: 6},
			{xPosition: 15, yPosition: 15, viewCode: 7},
		},
	}
}

func movePacman() tea.Cmd {
	return tea.Every(200*time.Millisecond, func(t time.Time) tea.Msg {
		return updatePacmanPosition(0)
	})
}

func moveGhosts() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return updateGhostsPosition(0)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		movePacman(),
		moveGhosts(),
	)
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

	case updateGhostsPosition:
		for _, ghost := range m.ghosts {
			prevX := ghost.xPosition
			prevY := ghost.yPosition

			positions := utils.CalculatePosibbleNextTile(m.board, ghost.xPosition, ghost.yPosition)
			randNum := rand.Intn(len(positions))

			newX := positions[randNum][0]
			newY := positions[randNum][1]

			ghost.xPosition = newX
			ghost.yPosition = newY

			if ghost.isCurrentPositionPallet {
				m.board[prevX][prevY] = 1
			} else {
				m.board[prevX][prevY] = 3
			}

			if m.board[newX][newY] == 1 {
				ghost.isCurrentPositionPallet = true
			} else {
				ghost.isCurrentPositionPallet = false
			}

			m.board[newX][newY] = ghost.viewCode
		}
		return m, moveGhosts()
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("\nSCORE: %v\n\n", m.score)

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
			case 4:
				s += blinkyStyle.Render(ghostStr)
			case 5:
				s += pinkyStyle.Render(ghostStr)
			case 6:
				s += inkyStyle.Render(ghostStr)
			case 7:
				s += clydeStyle.Render(ghostStr)
			}
		}
		s += "\n"
	}

	// Log pacman position
	s += fmt.Sprintf("pacman: {%v, %v}, direction: %v\n",
		m.pacman.xPosition,
		m.pacman.yPosition,
		m.pacman.direction,
	)

	// Log ghosts position
	for i, v := range m.ghosts {
		s += fmt.Sprintf("ghost %v: {%v, %v}\n", i, v.xPosition, v.yPosition)
	}

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
