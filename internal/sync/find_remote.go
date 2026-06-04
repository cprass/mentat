package sync

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

func findRemote(ctx context.Context, folder string, repoURL string) (string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	// searches for codeberg.org/user/repo
	searchStr1 := u.Host + u.Path
	// searches for github.com:user/repo
	searchStr2 := u.Host + ":" + strings.TrimPrefix(u.Path, "/")

	cmd := newCmd(ctx, folder, "git", "remote", "-v")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git error: %w, %s", err, out)
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, searchStr1) || strings.Contains(line, searchStr2) {
			fields := strings.Fields(line)
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("no remote matching the URL found")
}
