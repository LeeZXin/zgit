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

type DiffCommitsReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	Operator usermd.UserInfo
}

func (r *DiffCommitsReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
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

type SubmitPullRequestReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	Operator usermd.UserInfo
}

func (r *SubmitPullRequestReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
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
	RepoId    string
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
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
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

type DiffFileReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	FileName string
	Operator usermd.UserInfo
}

func (r *DiffFileReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
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

type DiffCommitsRespDTO struct {
	Target        string              `json:"target"`
	Head          string              `json:"head"`
	TargetCommit  CommitDTO           `json:"targetCommit"`
	HeadCommit    CommitDTO           `json:"headCommit"`
	Commits       []CommitDTO         `json:"commits"`
	NumFiles      int                 `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoDTO `json:"diffNumsStats"`
	ConflictFiles []string            `json:"conflictFiles"`
	CanMerge      bool                `json:"canMerge"`
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

type DiffFileRespDTO struct {
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

type ClosePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *ClosePullRequestReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.PrId) > 32 || len(r.PrId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.PrId) > 32 || len(r.PrId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}
