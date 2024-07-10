package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ary82/pacman/internal/game"
	"github.com/ary82/pacman/internal/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	port := os.Getenv("PORT")

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf(":%s", port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("starting SSH server", "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	renderer := bubbletea.MakeRenderer(s)

	block := renderer.NewStyle().Foreground(lipgloss.Color("4"))   // blue
	pacman := renderer.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	blinky := renderer.NewStyle().Foreground(lipgloss.Color("9"))  // red
	pinky := renderer.NewStyle().Foreground(lipgloss.Color("13"))  // pink
	inky := renderer.NewStyle().Foreground(lipgloss.Color("14"))   // cyan
	clyde := renderer.NewStyle().Foreground(lipgloss.Color("202")) // orange

	model := game.InitialGameModel()

	model.Styles = style.Styles{
		Block:  block,
		Pacman: pacman,
		Blinky: blinky,
		Pinky:  pinky,
		Inky:   inky,
		Clyde:  clyde,
	}

	return model, []tea.ProgramOption{}
}
