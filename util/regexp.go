package util

import "regexp"

var (
	ValidUserEmailRegPattern = regexp.MustCompile(`^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`)
	ValidUserAccountPattern  = regexp.MustCompile("^\\w{4,32}$")
	ValidCorpIdPattern       = regexp.MustCompile("^\\w{1,32}$")
	ValidRepoNamePattern     = regexp.MustCompile("^[\\w\\-]{1,32}$")
	ValidBranchPattern       = regexp.MustCompile("^\\w{1,32}$")
)
