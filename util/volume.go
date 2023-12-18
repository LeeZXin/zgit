package util

import "fmt"

const (
	Gib = 1024 * 1024 * 1024
	Mib = 1024 * 1024
	Kib = 1024
)

func VolumeReadable(b int64) string {
	if b > Gib {
		return fmt.Sprintf("%d Gib", b/Gib)
	}
	if b > Mib {
		return fmt.Sprintf("%d Mib", b/Mib)
	}
	if b > Kib {
		return fmt.Sprintf("%d Kib", b/Kib)
	}
	return fmt.Sprintf("%d bytes", b)
}
