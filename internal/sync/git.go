package sync

import (
	"context"
	"fmt"
	"mentat/internal/config"
	"time"
)

func SyncPull(ctx context.Context) error {
	// get the vault directory
	vault, err := config.Vault()
	if err != nil {
		return err
	}

	gitConfig := config.Git()
	if gitConfig.Url == "" {
		return fmt.Errorf("sync.git.url config is missing")
	}

	if ok, _ := isGitRepo(vault); !ok {
		return fmt.Errorf("vault %q is not a git repository", vault)
	}

	remote, err := findRemote(ctx, vault, gitConfig.Url)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := runCmd(timeoutCtx, vault, "git", "pull", "--ff-only", "-q", remote, gitConfig.Branch); err != nil {
		return err
	}

	return nil
}

func SyncPush(ctx context.Context) error {
	vault, err := config.Vault()
	if err != nil {
		return err
	}

	changes, err := newChanges(ctx, vault)
	if err != nil {
		return err
	}

	if changes.HasChanges() {
		gitConfig := config.Git()
		if gitConfig.Url == "" {
			return fmt.Errorf("sync.git.url config is missing")
		}

		if err := changes.Commit(ctx); err != nil {
			return err
		}

		remote, err := findRemote(ctx, vault, gitConfig.Url)
		if err != nil {
			return err
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := runCmd(timeoutCtx, vault, "git", "push", "-q", remote, gitConfig.Branch); err != nil {
			return err
		}
	}

	return nil
}
