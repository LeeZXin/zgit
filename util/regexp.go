package util

import "regexp"

var (
	ValidEmailRegPattern    = regexp.MustCompile(`^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`)
	ValidUserAccountPattern = regexp.MustCompile("^\\w{4,32}$")
	ValidCorpIdPattern      = regexp.MustCompile("^\\w{1,32}$")
	AtLeastOneCharPattern   = regexp.MustCompile("^\\w+$")
	ValidRepoNamePattern    = regexp.MustCompile("^[\\w\\-]{1,32}$")
	ValidBranchPattern      = regexp.MustCompile("^\\w{1,32}$")
)
