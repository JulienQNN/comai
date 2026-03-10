package git

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func Commit(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "git", "commit", "-m", message).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", string(out))
	}
	return nil
}
