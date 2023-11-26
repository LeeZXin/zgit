//go:build windows

package process

import (
	"os/exec"
)

func SetSysProcAttribute(_ *exec.Cmd) {
	// Do nothing
}
