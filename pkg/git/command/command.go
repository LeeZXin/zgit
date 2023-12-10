package command

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
	"zgit/pkg/git/process"
	"zgit/setting"
)

var (
	globalCmdArgs = make([]string, 0)

	passThroughEnvKeys = []string{
		"GNUPGHOME",
	}
)

// AddGlobalCmdArgs not thread safe
func AddGlobalCmdArgs(args ...string) {
	globalCmdArgs = append(globalCmdArgs, args...)
}

type Command struct {
	args []string
}

func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) Run(ctx context.Context, ros ...RunOpts) (*Result, error) {
	stdOut := new(bytes.Buffer)
	if err := c.run(ctx, stdOut, ros...); err != nil {
		return nil, err
	}
	return &Result{
		stdOut: stdOut,
	}, nil
}

func (c *Command) RunWithReadPipe(ctx context.Context, ros ...RunOpts) *ReadPipeResult {
	reader, writer := io.Pipe()
	go func() {
		if err := c.run(ctx, writer, ros...); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()
	return &ReadPipeResult{
		reader: reader,
	}
}

func (c *Command) RunWithStdinPipe(ctx context.Context, ros ...RunOpts) *ReadWritePipeResult {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	go func() {
		if err := c.run(ctx, stdoutWriter, append(ros, WithStdin(stdinReader))...); err != nil {
			stdinReader.CloseWithError(err)
			stdoutWriter.CloseWithError(err)
		} else {
			stdinReader.Close()
			stdoutWriter.Close()
		}
	}()
	return &ReadWritePipeResult{
		reader: stdoutReader,
		writer: stdinWriter,
	}
}

func (c *Command) run(ctx context.Context, stdOut io.Writer, ros ...RunOpts) error {
	opts := new(runOpts)
	for _, o := range ros {
		o(opts)
	}
	if ctx == nil {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(context.Background(), 360*time.Second)
		defer cancelFunc()
	}
	cmd := exec.CommandContext(ctx, setting.GitExecutablePath(), c.args...)
	if opts.Env == nil {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	stdErr := new(bytes.Buffer)
	process.SetSysProcAttribute(cmd)
	cmd.Env = append(cmd.Env, CommonGitCmdEnvs()...)
	cmd.Dir = opts.Dir
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	cmd.Stdin = opts.Stdin
	if err := cmd.Start(); err != nil {
		return stdErrorResult(err, bytesToString(stdErr.Bytes()))
	}
	err := cmd.Wait()
	if err != nil && ctx.Err() != context.DeadlineExceeded {
		return stdErrorResult(err, bytesToString(stdErr.Bytes()))
	}
	if ctx.Err() != nil {
		return stdErrorResult(ctx.Err(), bytesToString(stdErr.Bytes()))
	}
	return nil
}

func NewCommand(args ...string) *Command {
	if args == nil {
		args = []string{}
	}
	return &Command{
		args: append(globalCmdArgs, args...),
	}
}

func NewCommandWithNoGlobalArgs(args ...string) *Command {
	if args == nil {
		args = []string{}
	}
	return &Command{
		args: args,
	}
}

func CommonEnvs() []string {
	envs := []string{
		"HOME=" + setting.HomeDir(),
		"GIT_NO_REPLACE_OBJECTS=1",
	}
	for _, key := range passThroughEnvKeys {
		if val, ok := os.LookupEnv(key); ok {
			envs = append(envs, key+"="+val)
		}
	}
	return envs
}

func CommonGitCmdEnvs() []string {
	return append(CommonEnvs(), "LC_ALL=C", "GIT_TERMINAL_PROMPT=0")
}

func IsExitCode(err error, code int) bool {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode() == code
	}
	return false
}
