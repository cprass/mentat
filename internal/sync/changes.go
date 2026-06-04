package sync

import (
	"context"
	"fmt"
	"strings"
)

type changes struct {
	Reviews bool
	Vault   bool
	Folder  string
}

func newChanges(ctx context.Context, folder string) (changes, error) {
	c := changes{
		Folder: folder,
	}
	cmd := newCmd(ctx, folder, "git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c, fmt.Errorf("git error: %w, %s", err, out)
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.HasSuffix(line, "reviews.log") {
			c.Reviews = true
		}
		if strings.HasSuffix(line, ".md") {
			c.Vault = true
		}
		if c.Reviews && c.Vault {
			break
		}
	}

	return c, nil
}

func (c *changes) CommitMessage() string {
	parts := []string{}
	if c.Vault {
		parts = append(parts, "update vault")
	}
	if c.Reviews {
		parts = append(parts, "review")
	}
	return strings.Join(parts, " & ")
}

func (c *changes) HasChanges() bool {
	return c.Reviews || c.Vault
}

func (c *changes) Commit(ctx context.Context) error {
	if err := runCmd(ctx, c.Folder, "git", "add", "."); err != nil {
		return err
	}

	if err := runCmd(ctx, c.Folder, "git", "commit", "-q", "-m", c.CommitMessage()); err != nil {
		return err
	}

	return nil
}
