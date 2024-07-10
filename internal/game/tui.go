package game

import (
	"fmt"
	"math/rand"

	"github.com/ary82/pacman/internal/constants"
	"github.com/ary82/pacman/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Game) Init() tea.Cmd {
	return tea.Batch(
		MovePacman(),
		MoveGhosts(),
	)
}

func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "s":
			m.Pacman.Direction = 2
		case "k", "up", "w":
			m.Pacman.Direction = 1
		case "h", "left", "a":
			m.Pacman.Direction = 3
		case "l", "right", "d":
			m.Pacman.Direction = 4
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case UpdatePacmanPositionMsg:
		prevX := m.Pacman.XPosition
		prevY := m.Pacman.YPosition
		newX := prevX
		newY := prevY

		switch m.Pacman.Direction {
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

		switch newPosition := m.Board[newX][newY]; newPosition {
		case 1:
			m.Score += 1
			m.Pacman.XPosition = newX
			m.Pacman.YPosition = newY
			m.Board[prevX][prevY] = 3
			m.Board[newX][newY] = 2
		case 3:
			m.Pacman.XPosition = newX
			m.Pacman.YPosition = newY
			m.Board[prevX][prevY] = 3
			m.Board[newX][newY] = 2
		}

		if m.Score == 300 {
			return m, GameOver
		}
		return m, MovePacman()

	case UpdateGhostsPositionMsg:
		for _, ghost := range m.Ghosts {
			prevX := ghost.XPosition
			prevY := ghost.YPosition

			positions := utils.CalculatePosibbleNextTile(m.Board, ghost.XPosition, ghost.YPosition, ghost.Direction)
			randNum := rand.Intn(len(positions))

			newX := positions[randNum][0]
			newY := positions[randNum][1]
			newDirection := positions[randNum][2]

			ghost.XPosition = newX
			ghost.YPosition = newY
			ghost.Direction = newDirection

			if ghost.IsCurrentPositionPellet {
				m.Board[prevX][prevY] = 1
			} else {
				m.Board[prevX][prevY] = 3
			}

			if m.Board[newX][newY] == 1 {
				ghost.IsCurrentPositionPellet = true
			} else {
				ghost.IsCurrentPositionPellet = false
			}

			if m.Board[newX][newY] == 2 {
				m.Pacman.XPosition = 1
				m.Pacman.YPosition = 1
				m.Pacman.Direction = 0
				m.Pacman.Lives -= 1
			}

			if m.Pacman.Lives == 0 {
				return m, GameOver
			}

			m.Board[newX][newY] = ghost.ViewCode
		}
		return m, MoveGhosts()
	case GameOverMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m Game) View() string {
	s := fmt.Sprintf("\nSCORE: %v\tLIVES: %v\n\n", m.Score, m.Pacman.Lives)

	for _, v := range m.Board {
		for _, num := range v {
			switch num {
			case 0:
				s += m.Styles.Block.Render(constants.BlockStr)
			case 1:
				s += constants.PelletStr
			case 2:
				s += m.Styles.Pacman.Render(constants.PacmanStr)
			case 3:
				s += constants.EmptyStr
			case 4:
				s += m.Styles.Blinky.Render(constants.GhostStr)
			case 5:
				s += m.Styles.Pinky.Render(constants.GhostStr)
			case 6:
				s += m.Styles.Inky.Render(constants.GhostStr)
			case 7:
				s += m.Styles.Clyde.Render(constants.GhostStr)
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
