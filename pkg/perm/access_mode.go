package perm

// AccessMode specifies the users access mode
type AccessMode int

const (
	// AccessModeNone no access
	AccessModeNone AccessMode = iota // 0
	// AccessModeRead read access
	AccessModeRead // 1
	// AccessModeWrite write access
	AccessModeWrite // 2
)

func (mode AccessMode) Readable() string {
	switch mode {
	case AccessModeRead:
		return "read"
	case AccessModeWrite:
		return "write"
	default:
		return "none"
	}
}
