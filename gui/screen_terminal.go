package gui

import (
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/KiraCore/kensho/helper/gssh"
	"github.com/fyne-io/terminal"
	"golang.org/x/crypto/ssh"
)

type Terminal struct {
	Term                  *terminal.Terminal
	SSHSessionForTerminal *ssh.Session
	SSHIn                 io.WriteCloser
	SSHOut                io.Reader
}

func TryToRunSSHSessionForTerminal(g *Gui) (err error) {
	g.Terminal.SSHSessionForTerminal, err = gssh.MakeSSHsessionForTerminalV2(g.sshClient)
	if err != nil {
		return err
	}
	go func() {
		err := g.Terminal.SSHSessionForTerminal.Shell()
		if err != nil {
			log.Println("Shell err: ", err.Error())
		}
	}()

	g.Terminal.SSHIn, err = g.Terminal.SSHSessionForTerminal.StdinPipe()
	if err != nil {
		log.Println("SSHSessionForTerminal.StdinPipe", err)
		return err
	}
	g.Terminal.SSHOut, err = g.Terminal.SSHSessionForTerminal.StdoutPipe()
	if err != nil {
		log.Println("SSHSessionForTerminal.StdoutPipe", err)
		return err
	}
	g.Terminal.Term = terminal.New()

	ch := make(chan terminal.Config)
	go func() {
		rows, cols := uint(0), uint(0)
		for {
			config := <-ch
			if rows == config.Rows && cols == config.Columns {
				continue
			}
			rows, cols = config.Rows, config.Columns
			g.Terminal.SSHSessionForTerminal.WindowChange(int(rows), int(cols))
		}
	}()
	g.Terminal.Term.AddListener(ch)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in terminal: %v", r)
			}
		}()
		err := g.Terminal.Term.RunWithConnection(g.Terminal.SSHIn, g.Terminal.SSHOut)
		if err != nil {
			log.Println("RunWithConnection err: ", err.Error())
		}
	}()

	return nil
}

func makeTerminalScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {

	return container.NewVScroll(container.NewBorder(nil, nil, nil, nil, g.Terminal.Term))
	// return container.NewVScroll(container.NewStack())
}
