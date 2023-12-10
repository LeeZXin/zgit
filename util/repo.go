package util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	InvalidRepoNamePattern = regexp.MustCompile(`[^\w-\.]`)
	AlphaPattern           = regexp.MustCompile("^\\w.+$")
)

var (
	reservedRepoNames    = []string{".", "..", "-"}
	reservedRepoPatterns = []string{"*.git", "*.wiki"}
)

type RelativeRepoPath struct {
	RowData   string
	CompanyId string
	ClusterId string
	RepoName  string
}

func IsValidRepoName(name string) error {
	if InvalidRepoNamePattern.MatchString(name) {
		return fmt.Errorf("name is invalid [%s]: must be valid alpha or numeric or dash(-_) or dot characters", name)
	}
	return isValidName(reservedRepoNames, reservedRepoPatterns, name)
}

func isValidName(names, patterns []string, name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return fmt.Errorf("name is empty")
	}
	for i := range names {
		if name == names[i] {
			return fmt.Errorf("name is reserved [name: %s]", name)
		}
	}
	for _, pat := range patterns {
		if pat[0] == '*' && strings.HasSuffix(name, pat[1:]) ||
			(pat[len(pat)-1] == '*' && strings.HasPrefix(name, pat[:len(pat)-1])) {
			return fmt.Errorf("name pattern is not allowed [pattern: %s]", patterns)
		}
	}
	return nil
}

func ParseRelativeRepoPath(relativeRepoPath string) (RelativeRepoPath, error) {
	path := strings.TrimPrefix(relativeRepoPath, "/")
	splits := strings.Split(path, "/")
	if len(splits) != 3 {
		return RelativeRepoPath{}, errors.New("invalid repo path")
	}
	if !AlphaPattern.MatchString(splits[0]) {
		return RelativeRepoPath{}, errors.New("invalid company id")
	}
	if !AlphaPattern.MatchString(splits[1]) {
		return RelativeRepoPath{}, errors.New("invalid cluster id")
	}
	if InvalidRepoNamePattern.MatchString(splits[2]) {
		return RelativeRepoPath{}, errors.New("invalid repo path")
	}
	return RelativeRepoPath{
		RowData:   relativeRepoPath,
		CompanyId: strings.ToLower(splits[0]),
		ClusterId: strings.ToLower(splits[1]),
		RepoName:  strings.ToLower(strings.TrimSuffix(splits[2], ".git")),
	}, nil
}
