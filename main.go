package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ary82/pacman/internal/style"
	"github.com/ary82/pacman/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
	board  [31][28]int
	score  int
	pacman pacman
	ghosts [4]*ghost
}

type pacman struct {
	xPosition int
	yPosition int
	direction int // 0: static, 1: up, 2: down, 3: left 4: right
	lives     int
}

type ghost struct {
	xPosition               int
	yPosition               int
	viewCode                int
	isCurrentPositionPellet bool
	direction               int // 0: static, 1: up, 2: down, 3: left 4: right
}

type (
	updatePacmanPosition int
	updateGhostsPosition int
	gameOverMsg          int
)

const (
	blockStr  string = "███"
	pacmanStr string = " 󰮯 "
	ghostStr  string = " 󰊠 "
	emptyStr  string = "   "
	pelletStr string = " 🞄 "
)

func initialGameModel() game {
	return game{
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
			lives:     3,
		},
		score: 0,
		ghosts: [4]*ghost{
			{xPosition: 15, yPosition: 12, viewCode: 4},
			{xPosition: 15, yPosition: 13, viewCode: 5},
			{xPosition: 15, yPosition: 14, viewCode: 6},
			{xPosition: 15, yPosition: 15, viewCode: 7},
		},
	}
}

func movePacman() tea.Cmd {
	return tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg {
		return updatePacmanPosition(0)
	})
}

func moveGhosts() tea.Cmd {
	return tea.Every(350*time.Millisecond, func(t time.Time) tea.Msg {
		return updateGhostsPosition(0)
	})
}

func gameOver() tea.Msg {
	return gameOverMsg(0)
}

func (m game) Init() tea.Cmd {
	return tea.Batch(
		movePacman(),
		moveGhosts(),
	)
}

func (m game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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

		if m.score == 300 {
			return m, gameOver
		}
		return m, movePacman()

	case updateGhostsPosition:
		for _, ghost := range m.ghosts {
			prevX := ghost.xPosition
			prevY := ghost.yPosition

			positions := utils.CalculatePosibbleNextTile(m.board, ghost.xPosition, ghost.yPosition, ghost.direction)
			randNum := rand.Intn(len(positions))

			newX := positions[randNum][0]
			newY := positions[randNum][1]
			newDirection := positions[randNum][2]

			ghost.xPosition = newX
			ghost.yPosition = newY
			ghost.direction = newDirection

			if ghost.isCurrentPositionPellet {
				m.board[prevX][prevY] = 1
			} else {
				m.board[prevX][prevY] = 3
			}

			if m.board[newX][newY] == 1 {
				ghost.isCurrentPositionPellet = true
			} else {
				ghost.isCurrentPositionPellet = false
			}

			if m.board[newX][newY] == 2 {
				m.pacman.xPosition = 1
				m.pacman.yPosition = 1
				m.pacman.direction = 0
				m.pacman.lives -= 1
			}

			if m.pacman.lives == 0 {
				return m, gameOver
			}

			m.board[newX][newY] = ghost.viewCode
		}
		return m, moveGhosts()
	case gameOverMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m game) View() string {
	s := fmt.Sprintf("\nSCORE: %v\tLIVES: %v\n\n", m.score, m.pacman.lives)

	for _, v := range m.board {
		for _, num := range v {
			switch num {
			case 0:
				s += style.Block.Render(blockStr)
			case 1:
				s += pelletStr
			case 2:
				s += style.Pacman.Render(pacmanStr)
			case 3:
				s += emptyStr
			case 4:
				s += style.Blinky.Render(ghostStr)
			case 5:
				s += style.Pinky.Render(ghostStr)
			case 6:
				s += style.Inky.Render(ghostStr)
			case 7:
				s += style.Clyde.Render(ghostStr)
			}
		}
		s += "\n"
	}

	// Log pacman position
	// s += fmt.Sprintf("pacman: {%v, %v}, direction: %v\n",
	// 	m.pacman.xPosition,
	// 	m.pacman.yPosition,
	// 	m.pacman.direction,
	// )

	// Log ghosts position
	// for i, v := range m.ghosts {
	// 	s += fmt.Sprintf("ghost %v: {%v, %v}\n", i, v.xPosition, v.yPosition)
	// }

	return s
}

func main() {
	model := initialGameModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
