package pullrequestsrv

import (
	"time"
	"zgit/pkg/git"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

const (
	UpDirection   = "up"
	DownDirection = "down"
)

type PreparePullRequestReqDTO struct {
	RepoPath string
	Target   string
	Head     string
	Operator usermd.UserInfo
}

func (r *PreparePullRequestReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoPath) > 255 || len(r.RepoPath) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Target) > 128 || len(r.Target) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Head) > 128 || len(r.Head) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type CatFileReqDTO struct {
	RepoPath  string
	CommitId  string
	FileName  string
	Offset    int
	Limit     int
	Direction string
	Operator  usermd.UserInfo
}

func (r *CatFileReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoPath) > 255 || len(r.RepoPath) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.CommitId) == 0 || len(r.CommitId) > 128 {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if len(r.Direction) == 0 || len(r.Direction) > 10 {
		return util.InvalidArgsError()
	}
	if r.Direction != UpDirection && r.Direction != DownDirection {
		return util.InvalidArgsError()
	}
	if len(r.FileName) > 255 || len(r.FileName) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffReqDTO struct {
	RepoPath string
	Target   string
	Head     string
	FileName string
	Operator usermd.UserInfo
}

func (r *DiffReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoPath) > 255 || len(r.RepoPath) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Target) > 128 || len(r.Target) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Head) > 128 || len(r.Head) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.FileName) > 255 || len(r.FileName) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type CommitDTO struct {
	Author        git.User
	Committer     git.User
	AuthoredDate  time.Time
	CommittedDate time.Time
	CommitMsg     string
	CommitId      string
	ShortId       string
}

type PreparePullRequestRespDTO struct {
	Target        string              `json:"target"`
	Head          string              `json:"head"`
	TargetCommit  CommitDTO           `json:"targetCommit"`
	HeadCommit    CommitDTO           `json:"headCommit"`
	Commits       []CommitDTO         `json:"commits"`
	NumFiles      int                 `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoDTO `json:"diffNumsStats"`
}

type DiffNumsStatInfoDTO struct {
	FileChangeNums int `json:"fileChangeNums"`
	InsertNums     int `json:"insertNums"`
	DeleteNums     int `json:"deleteNums"`
	Stats          []DiffNumsStatDTO
}

type DiffNumsStatDTO struct {
	RawPath    string
	Path       string
	TotalNums  int
	InsertNums int
	DeleteNums int
}

type DiffRespDTO struct {
	FilePath    string
	OldMode     string
	Mode        string
	IsSubModule bool
	FileType    git.DiffFileType
	IsBinary    bool
	RenameFrom  string
	RenameTo    string
	CopyFrom    string
	CopyTo      string
	Lines       []DiffLineDTO
}

type DiffLineDTO struct {
	Index   int
	LeftNo  int
	Prefix  string
	RightNo int
	Text    string
}
