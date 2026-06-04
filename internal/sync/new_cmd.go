package sync

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Returns a new command that kills all child processes on termination
// and disables the git terminal prompt
func newCmd(ctx context.Context, dir string, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	return cmd
}

// Runs a command, captures the stderr output and returns that as an error on failure.
func runCmd(ctx context.Context, dir string, name string, args ...string) error {
	cmd := newCmd(ctx, dir, name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git error: %w: %s", err, out)
	}
	return nil
}
