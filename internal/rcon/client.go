// Package rcon implements the interactive and single-shot RCON sessions
// used by the rcon-cli command.
package rcon

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gorcon/rcon"
	"github.com/peterh/liner"
)

// Start opens an RCON connection to hostPort and runs an interactive REPL
// on it, reading commands from stdin via liner and writing colorized
// responses to out. It returns when the user sends "exit", aborts the
// prompt (Ctrl-D/Ctrl-C), or the connection is closed by the server.
func Start(hostPort, password string, out io.Writer) error {
	conn, err := rcon.Dial(hostPort, password)
	if err != nil {
		return fmt.Errorf("connecting to RCON server: %w", err)
	}
	defer func() { _ = conn.Close() }()

	lineEditor := liner.NewLiner()
	defer func() { _ = lineEditor.Close() }()

	for {
		cmd, err := lineEditor.Prompt("> ")
		if err != nil {
			if errors.Is(err, liner.ErrPromptAborted) || errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("reading input: %w", err)
		}

		if cmd == "exit" {
			return nil
		}

		resp, err := conn.Execute(cmd)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			if _, werr := fmt.Fprintln(os.Stderr, "Failed to execute command:", err); werr != nil {
				return fmt.Errorf("writing to stderr: %w", werr)
			}
			continue
		}

		if _, err := fmt.Fprintln(out, colorize(resp)); err != nil {
			return fmt.Errorf("writing response: %w", err)
		}
		lineEditor.AppendHistory(cmd)
	}
}

// Execute opens an RCON connection to hostPort, sends command (its parts
// joined by spaces) as a single command, writes the colorized response to
// out, and closes the connection.
func Execute(hostPort, password string, out io.Writer, command ...string) error {
	conn, err := rcon.Dial(hostPort, password)
	if err != nil {
		return fmt.Errorf("connecting to RCON server: %w", err)
	}
	defer func() { _ = conn.Close() }()

	resp, err := conn.Execute(strings.Join(command, " "))
	if err != nil {
		return fmt.Errorf("executing command: %w", err)
	}

	if _, err := fmt.Fprintln(out, colorize(resp)); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	return nil
}
