package command

import "io"

type ReadCloserWrapper interface {
	io.Reader
	CloseWithError(error) error
}

type runOpts struct {
	Env          []string
	Dir          string
	Stdin        ReadCloserWrapper
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

func withStdin(reader ReadCloserWrapper) RunOpts {
	return func(opts *runOpts) {
		opts.Stdin = reader
	}
}

func WithPipelineFunc(fn func() error) RunOpts {
	return func(opts *runOpts) {
		opts.PipelineFunc = fn
	}
}
