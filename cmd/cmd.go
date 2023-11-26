package cmd

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func initWaitContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(
			signalChannel,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		select {
		case <-signalChannel:
		case <-ctx.Done():
		}
		cancel()
		signal.Reset()
	}()
	return ctx, cancel
}

func exitWithMsg(code int, msg string) error {
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stderr, "zgit:", msg)
	logger.Logger.Error(msg)
	return cli.Exit("", code)
}

func exitWithDefaultCode(msg string) error {
	return exitWithMsg(1, msg)
}
