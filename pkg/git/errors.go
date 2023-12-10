package git

type ErrMergeConflict struct {
	err error
}

func (e *ErrMergeConflict) Error() string {
	return e.err.Error()
}

type ErrMergeUnrelatedHistories struct {
	err error
}

func (e *ErrMergeUnrelatedHistories) Error() string {
	return e.err.Error()
}

type ErrPushOutOfDate struct {
	err error
}

func (e *ErrPushOutOfDate) Error() string {
	return e.err.Error()
}

type ErrPushRejected struct {
	err error
}

func (e *ErrPushRejected) Error() string {
	return e.err.Error()
}
