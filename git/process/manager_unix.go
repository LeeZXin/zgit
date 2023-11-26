//go:build !windows

package process

import (
	"os/exec"
	"syscall"
)

func SetSysProcAttribute(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
