package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	AlphaDashDotPattern = regexp.MustCompile(`[^\w-\.]`)
)

var (
	reservedRepoNames    = []string{".", "..", "-"}
	reservedRepoPatterns = []string{"*.git", "*.wiki", "*.rss", "*.atom"}
)

func IsValidRepoName(name string) error {
	if AlphaDashDotPattern.MatchString(name) {
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
