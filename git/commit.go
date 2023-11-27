package git

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zgit/git/command"
	"zgit/gpg"
)

const (
	MissingType = "missing"
	CommitType  = "commit"
	TagType     = "tag"
)

var commitIdPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

type Commit struct {
	Id               string     `json:"id"`
	Tree             *Tree      `json:"tree"`
	Parent           []string   `json:"parent"`
	Author           *User      `json:"author"`
	AuthorSigTime    *time.Time `json:"authorSigTime"`
	Committer        *User      `json:"committer"`
	CommitterSigTime *time.Time `json:"committerSigTime"`
	GpgSig           string     `json:"gpgSig"`
	CommitMsg        string     `json:"commitMsg"`
}

func NewCommit(id string) *Commit {
	return &Commit{
		Id:     id,
		Parent: make([]string, 0),
	}
}

type Tag struct {
	Id        string     `json:"id"`
	Object    string     `json:"object"`
	Typ       string     `json:"typ"`
	Tag       string     `json:"tag"`
	Tagger    *User      `json:"tagger"`
	TagTime   *time.Time `json:"tagTime"`
	CommitMsg string     `json:"commitMsg"`
}

type Tree struct {
	Id string `json:"id"`
}

func NewTree(id string) *Tree {
	return &Tree{
		Id: id,
	}
}

type CommitID struct {
	o string
	b []byte
	s string
}

func (c *CommitID) OriginalStr() string {
	return c.o
}

func (c *CommitID) BytesContent() []byte {
	b := make([]byte, 20)
	copy(c.b, b)
	return b
}

func (c *CommitID) ShortStr() string {
	return c.s
}

func NewCommitIDFromHexStr(str string) (CommitID, error) {
	if !commitIdPattern.MatchString(str) {
		return CommitID{}, errors.New("commitId is not valid")
	}
	bs, err := hex.DecodeString(str)
	if err != nil {
		return CommitID{}, err
	}
	b := make([]byte, 20)
	copy(bs, b)
	return CommitID{
		o: str,
		b: b,
		s: str[0:7],
	}, nil
}

func GetRefCommitId(ctx context.Context, repoPath string, name string) (string, error) {
	commitId, _, _, err := CatFileBatchCheck(ctx, repoPath, name)
	return commitId, err
}

func GetCommitByCommitId(ctx context.Context, repoPath string, commitId string) (*Commit, error) {
	c := NewCommit(commitId)
	return c, CatFileBatch(ctx, repoPath, commitId, func(r io.Reader, closer command.PipeResultCloser) error {
		defer closer.ClosePipe()
		reader := bufio.NewReader(r)
		var (
			typ  string
			size int64
		)
		for {
			line, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("read line err: %v", err)
			}
			if isPrefix {
				continue
			}
			_, typ, size, err = readBatchLine(string(line))
			if err != nil {
				return fmt.Errorf("readBatchLine err: %v", err)
			}
			break
		}
		switch typ {
		case MissingType:
			return fmt.Errorf("%s is missing", commitId)
		case CommitType:
			return genCommit(io.LimitReader(reader, size), c)
		default:
			return fmt.Errorf("unsupported type: %s", typ)
		}
	})
}

func GetCommitByTag(ctx context.Context, repoPath string, tag string) (c *Commit, e error) {
	e = CatFileBatch(ctx, repoPath, tag, func(r io.Reader, closer command.PipeResultCloser) error {
		defer closer.ClosePipe()
		reader := bufio.NewReader(r)
		var (
			typ  string
			size int64
			id   string
		)
		for {
			line, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("read line err: %v", err)
			}
			if isPrefix {
				continue
			}
			id, typ, size, err = readBatchLine(string(line))
			if err != nil {
				return fmt.Errorf("readBatchLine err: %v", err)
			}
			break
		}
		switch typ {
		case MissingType:
			return fmt.Errorf("%s is missing", tag)
		case TagType:
			t := &Tag{
				Id:  id,
				Tag: tag,
			}
			err := genTag(io.LimitReader(reader, size), t)
			if err != nil {
				return fmt.Errorf("parse Tag err: %v", err)
			}
			if t.Object == "" {
				return fmt.Errorf("%s object is empty", tag)
			}
			c, err = GetCommitByCommitId(ctx, repoPath, t.Object)
			return err
		default:
			return fmt.Errorf("unsupported type: %s", typ)
		}
	})
	return
}

func genTag(r io.Reader, tag *Tag) error {
	reader := bufio.NewReader(r)
	commitMsg := strings.Builder{}
	defer func() {
		tag.CommitMsg = commitMsg.String()
	}()
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read line err: %v", err)
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		fields := strings.Fields(lineStr)
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {
		case "object":
			tag.Object = fields[1]
		case "type":
			tag.Typ = fields[1]
		case "tag":
			tag.Tag = fields[1]
		case "tagger":
			tag.Tagger, tag.TagTime = parseUserAndTime(fields[1:])
		default:
			commitMsg.WriteString(lineStr)
		}
	}
}

func genCommit(r io.Reader, commit *Commit) error {
	reader := bufio.NewReader(r)
	commitMsg := strings.Builder{}
	defer func() {
		commit.CommitMsg = commitMsg.String()
	}()
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read line err: %v", err)
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		fields := strings.Fields(lineStr)
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {
		case "tree":
			commit.Tree = NewTree(fields[1])
		case "parent":
			commit.Parent = append(commit.Parent, fields[1])
		case "author":
			commit.Author, commit.AuthorSigTime = parseUserAndTime(fields[1:])
		case "committer":
			commit.Committer, commit.CommitterSigTime = parseUserAndTime(fields[1:])
		case "gpgsig":
			sigPayload := strings.Builder{}
			sigPayload.WriteString(fields[1])
			for {
				line, isPrefix, err = reader.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("read gpgsig err: %v", err)
				}
				if isPrefix {
					continue
				}
				lineStr := string(line)
				sigPayload.WriteString(lineStr)
				if strings.TrimSpace(lineStr) == gpg.EndLineTag {
					break
				}
			}
			commit.GpgSig = sigPayload.String()
		default:
			commitMsg.WriteString(lineStr + "\n")
		}
	}
}

func parseUserAndTime(f []string) (*User, *time.Time) {
	u := User{}
	l := len(f)
	if l >= 1 {
		u.Name = f[0]
	}
	if l >= 2 {
		c := make([]byte, len(f[1])-2)
		copy(c, f[1][1:len(f[1])-1])
		u.Email = string(c)
	}
	var eventTime time.Time
	if l >= 3 {
		firstChar := f[2][0]
		if firstChar >= 48 && firstChar <= 57 {
			i, err := strconv.ParseInt(f[2], 10, 64)
			if err == nil {
				eventTime = time.Unix(i, 0)
			}
			if l >= 4 {
				zone := f[3]
				h, herr := strconv.ParseInt(zone[0:3], 10, 64)
				m, merr := strconv.ParseInt(zone[3:], 10, 64)
				if herr == nil && merr == nil {
					if h < 0 {
						m = -m
					}
					eventTime = eventTime.In(time.FixedZone("", int(h*3600+m*60)))
				}
			}
		} else {
			i, err := time.Parse(TimeLayout, f[2])
			if err == nil {
				eventTime = i
			}
		}
	}
	return &u, &eventTime
}
