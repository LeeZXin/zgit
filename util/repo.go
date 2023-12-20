package util

import (
	"errors"
	"strings"
)

type RelativeRepoPath struct {
	RowData   string
	CorpId    string
	ClusterId string
	RepoName  string
}

func ParseRelativeRepoPath(relativeRepoPath string) (RelativeRepoPath, error) {
	path := strings.TrimPrefix(relativeRepoPath, "/")
	splits := strings.Split(path, "/")
	if len(splits) != 3 {
		return RelativeRepoPath{}, errors.New("invalid repo path")
	}
	return RelativeRepoPath{
		RowData:   relativeRepoPath,
		CorpId:    strings.ToLower(splits[0]),
		ClusterId: strings.ToLower(splits[1]),
		RepoName:  strings.ToLower(strings.TrimSuffix(splits[2], ".git")),
	}, nil
}
