package git

import (
	"fmt"
	"os/exec"
)

func Commit(message string) error {
	out, err := exec.Command("git", "commit", "-m", message).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", string(out))
	}
	return nil
}
