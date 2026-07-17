package rcon

import (
	"bytes"
	"testing"

	gorcon "github.com/gorcon/rcon"
	"github.com/gorcon/rcon/rcontest"
)

func TestExecute_Success(t *testing.T) {
	srv := rcontest.NewServer(
		rcontest.SetSettings(rcontest.Settings{Password: "testpass"}),
		rcontest.SetCommandHandler(func(c *rcontest.Context) {
			resp := "§aok §x§f§f§0§0§0§0red"
			_, _ = gorcon.NewPacket(gorcon.SERVERDATA_RESPONSE_VALUE, c.Request().ID, resp).WriteTo(c.Conn())
		}),
	)
	defer srv.Close()

	var buf bytes.Buffer
	if err := Execute(srv.Addr(), "testpass", &buf, "some", "command"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "\x1b[32mok \x1b[38;2;255;0;0mred\n"
	if buf.String() != want {
		t.Errorf("Execute() output = %q, want %q", buf.String(), want)
	}
}

func TestExecute_WrongPassword(t *testing.T) {
	srv := rcontest.NewServer(
		rcontest.SetSettings(rcontest.Settings{Password: "correct"}),
	)
	defer srv.Close()

	var buf bytes.Buffer
	err := Execute(srv.Addr(), "wrong", &buf, "some", "command")
	if err == nil {
		t.Fatal("Execute() expected an error for wrong password, got nil")
	}
}

func TestExecute_ConnectionRefused(t *testing.T) {
	var buf bytes.Buffer
	err := Execute("127.0.0.1:1", "pw", &buf, "cmd")
	if err == nil {
		t.Fatal("Execute() expected an error when dialing an unreachable address, got nil")
	}
}
