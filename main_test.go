package main_test

import (
	"os/exec"
	"testing"
)

// TODO: Add an "end-to-end" test where we spin up a server
// that serves mock data and the client downloads it.

func TestValidateFlags(t *testing.T) {
	var (
		cmd = exec.Command("go", "run", "main.go")
		err = cmd.Run()
	)
	if err == nil {
		t.Error("cmd.Run = nil; want non-nil error")
	}
	cmd = exec.Command(
		"go", "run", "main.go",
		"-auth.email", "support@cloudflare.com",
	)
	err = cmd.Run()
	if err == nil {
		t.Error("cmd.Run = nil; want non-nil error")
	}
	cmd = exec.Command(
		"go", "run", "main.go",
		"-auth.email", "support@cloudflare.com",
		"-auth.key", "CF",
		"-url", "",
	)
	err = cmd.Run()
	if err == nil {
		t.Error("cmd.Run = nil; want non-nil error")
	}
	cmd = exec.Command(
		"go", "run", "main.go",
		"-auth.email", "support@cloudflare.com",
		"-auth.key", "CF",
		"start", "0",
		"-end", "-1",
	)
	err = cmd.Run()
	if err == nil {
		t.Error("cmd.Run = nil; want non-nil error")
	}
	cmd = exec.Command(
		"go", "run", "main.go",
		"-auth.email", "support@cloudflare.com",
		"-auth.key", "CF",
		"-start", "42",
		"-end", "42",
	)
	err = cmd.Run()
	if err == nil {
		t.Error("cmd.Run = nil; want non-nil error")
	}
}
