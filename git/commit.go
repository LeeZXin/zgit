package git

import (
	"encoding/hex"
	"errors"
	"regexp"
)

var commitIdPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

type Commit struct {
}

type CommitId struct {
	o string
	b []byte
	s string
}

func (c *CommitId) OriginalStr() string {
	return c.o
}

func (c *CommitId) BytesContent() []byte {
	b := make([]byte, 20)
	copy(c.b, b)
	return b
}

func (c *CommitId) ShortStr() string {
	return c.s
}

func NewCommitIdFromHexStr(str string) (CommitId, error) {
	if !commitIdPattern.MatchString(str) {
		return CommitId{}, errors.New("commitId is not valid")
	}
	bs, err := hex.DecodeString(str)
	if err != nil {
		return CommitId{}, err
	}
	b := make([]byte, 20)
	copy(bs, b)
	return CommitId{
		o: str,
		b: b,
		s: str[0:7],
	}, nil
}
