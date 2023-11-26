package setting

import "github.com/LeeZXin/zsf/property/static"

var (
	lfsEnabled = static.GetBool("lfs.enabled")
)

func LfsEnabled() bool {
	return lfsEnabled
}
