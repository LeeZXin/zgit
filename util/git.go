package util

func LongCommitId2ShortId(commitId string) string {
	if len(commitId) < 7 {
		return commitId
	}
	return commitId[:7]
}
