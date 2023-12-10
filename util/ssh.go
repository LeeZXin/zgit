package util

import (
	"fmt"
	"github.com/gliderlabs/ssh"
)

func ExitWithErrMsg(session ssh.Session, msg string) {
	fmt.Fprintf(session.Stderr(), msg)
	session.Exit(1)
}
