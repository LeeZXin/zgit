package command

import "io"

type runOpts struct {
	Env          []string
	Dir          string
	Stdin        io.Reader
	PipelineFunc func() error
}

type RunOpts func(*runOpts)

func WithEnv(env []string) RunOpts {
	return func(opts *runOpts) {
		opts.Env = env
	}
}

func WithDir(dir string) RunOpts {
	return func(opts *runOpts) {
		opts.Dir = dir
	}
}

func WithStdin(reader io.Reader) RunOpts {
	return func(opts *runOpts) {
		opts.Stdin = reader
	}
}
