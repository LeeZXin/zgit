package git

import "encoding/hex"

type Commit struct {
}

type CommitId struct {
	o string
	b []byte
}

func (c *CommitId) OriginalStr() string {
	return c.o
}

func (c *CommitId) BytesContent() []byte {
	b := make([]byte, 20)
	copy(c.b, b)
	return b
}

func NewCommitIdFromHexStr(str string) (CommitId, error) {
	bs, err := hex.DecodeString(str)
	if err != nil {
		return CommitId{}, err
	}
	b := make([]byte, 20)
	copy(bs, b)
	return CommitId{
		o: str,
		b: b,
	}, nil
}
