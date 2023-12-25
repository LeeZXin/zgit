package git

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf-utils/listutil"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zgit/pkg/git/command"
	"zgit/pkg/git/signature"
)

const (
	CommitType = "commit"
	TagType    = "tag"
)

var (
	ShortCommitIdPattern = regexp.MustCompile(`^[0-9a-f]{7}$`)
)

type Commit struct {
	Id            string
	Tree          *Tree
	Parent        []string
	Author        User
	AuthorSigTime time.Time
	Committer     User
	CommitSigTime time.Time
	Tagger        User
	TagSigTime    time.Time
	GpgSig        signature.GPGSig
	CommitMsg     string
	Tag           *Tag
	Payload       string
}

func (c *Commit) VerifyGPGSignature(pubKeys ...*signature.GPGPublicKey) error {
	if c.GpgSig == "" {
		return errors.New("no need to verify gpg")
	}
	if len(pubKeys) == 0 {
		return errors.New("empty pubKeys")
	}
	sig, err := signature.ParseGPGSignature(c.GpgSig)
	if err != nil {
		return err
	}
	hashContent, err := sig.HashObject([]byte(c.Payload))
	if err != nil {
		return err
	}
	for _, key := range pubKeys {
		err = key.VerifySignature(hashContent, sig.Signature)
		if err == nil {
			return nil
		}
	}
	return errors.New("verify failed")
}

func (c *Commit) VerifySshSignature(publicKey string) error {
	if c.GpgSig == "" {
		return errors.New("no need to verify gpg")
	}
	if publicKey == "" {
		return errors.New("empty publicKey")
	}
	return signature.VerifySshSignature(c.GpgSig.String(), c.Payload, publicKey)
}

func newCommit(id string) Commit {
	return Commit{
		Id:     id,
		Parent: make([]string, 0),
	}
}

type Tag struct {
	Id        string
	Object    string
	Typ       string
	Tag       string
	Tagger    User
	TagTime   time.Time
	CommitMsg string
}

type Tree struct {
	Id string `json:"id"`
}

func NewTree(id string) *Tree {
	return &Tree{
		Id: id,
	}
}

func GetRefCommitId(ctx context.Context, repoPath string, name string) (string, error) {
	cmd := command.NewCommand("rev-parse", "--verify", name)
	result, err := cmd.Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func CheckRefIsCommit(ctx context.Context, repoPath string, name string) bool {
	return CheckExists(ctx, repoPath, name)
}

func GetCommitByCommitId(ctx context.Context, repoPath string, commitId string) (Commit, error) {
	c := newCommit(commitId)
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
		case CommitType:
			return genCommit(io.LimitReader(reader, size), &c)
		default:
			return fmt.Errorf("unsupported type: %s", typ)
		}
	})
}

func GetCommitByTag(ctx context.Context, repoPath string, tag string) (c Commit, e error) {
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
			c.Tag = t
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
	payload := strings.Builder{}
	defer func() {
		commit.CommitMsg = commitMsg.String()
		commit.Payload = payload.String()
	}()
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			continue
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read line err: %v", err)
		}
		rowLine := string(line)
		lineStr := strings.TrimSpace(rowLine)
		fields := strings.Fields(lineStr)
		// 记录非签名数据
		if len(fields) < 1 || fields[0] != "gpgsig" {
			payload.WriteString(rowLine)
			payload.WriteString("\n")
		}
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
		case "tagger":
			commit.Tagger, commit.TagSigTime = parseUserAndTime(fields[1:])
		case "committer":
			commit.Committer, commit.CommitSigTime = parseUserAndTime(fields[1:])
		case "gpgsig":
			if len(fields) <= 1 {
				continue
			}
			sigPayload := strings.Builder{}
			sigPayload.WriteString(strings.TrimPrefix(lineStr, "gpgsig ") + "\n")
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
				lineStr = strings.TrimSpace(string(line))
				sigPayload.WriteString(lineStr + "\n")
				if lineStr == signature.EndGPGSigLineTag || lineStr == signature.EndSSHSigLineTag {
					break
				}
			}
			commit.GpgSig = signature.GPGSig(sigPayload.String())
		default:
			commitMsg.WriteString(lineStr + "\n")
		}
	}
}

func parseUserAndTime(f []string) (User, time.Time) {
	u := User{}
	l := len(f)
	if l >= 1 {
		u.Account = f[0]
	}
	if l >= 2 {
		u.Email = f[1][1 : len(f[1])-1]
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
	return u, eventTime
}

func GetFullShaCommitId(ctx context.Context, repoPath, short string) (string, error) {
	if ShortCommitIdPattern.MatchString(short) {
		line, _, _, err := CatFileBatchCheck(ctx, repoPath, short)
		return line, err
	}
	return short, nil
}

func GetGitLogCommitList(ctx context.Context, repoPath, target, head string) ([]Commit, error) {
	result, err := command.NewCommand("log", PrettyLogFormat, target+".."+head, "--max-count=500", "--").
		Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return nil, err
	}
	idList := strings.Fields(strings.TrimSpace(result.ReadAsString()))
	return listutil.Map(idList, func(t string) (Commit, error) {
		return GetCommitByCommitId(ctx, repoPath, t)
	})
}

func GetFileLastCommit(ctx context.Context, repoPath, refName, filePath string) (Commit, error) {
	result, err := command.NewCommand("log", PrettyLogFormat, "-1", refName, "--", filePath).
		Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return Commit{}, err
	}
	commitId := strings.TrimSpace(result.ReadAsString())
	return GetCommitByCommitId(ctx, repoPath, commitId)
}
