package pullrequestapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/pkg/git"
)

type PreparePullRequestReqVO struct {
	RepoPath string `json:"repoPath"`
	Target   string `json:"target"`
	Head     string `json:"head"`
}

type DiffReqVO struct {
	RepoPath string `json:"repoPath"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FileName string `json:"fileName"`
}

type CommitVO struct {
	Author        git.User
	Committer     git.User
	AuthoredDate  string
	CommittedDate string
	CommitMsg     string
	CommitId      string
	ShortId       string
}

type PreparePullRequestRespVO struct {
	ginutil.BaseResp
	Target        string             `json:"target"`
	Head          string             `json:"head"`
	TargetCommit  CommitVO           `json:"targetCommit"`
	HeadCommit    CommitVO           `json:"headCommit"`
	Commits       []CommitVO         `json:"commits"`
	NumFiles      int                `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoVO `json:"diffNumsStats"`
}

type DiffNumsStatInfoVO struct {
	FileChangeNums int              `json:"fileChangeNums"`
	InsertNums     int              `json:"insertNums"`
	DeleteNums     int              `json:"deleteNums"`
	Stats          []DiffNumsStatVO `json:"stats"`
}

type DiffNumsStatVO struct {
	RawPath    string `json:"rawPath"`
	Path       string `json:"path"`
	TotalNums  int    `json:"totalNums"`
	InsertNums int    `json:"insertNums"`
	DeleteNums int    `json:"deleteNums"`
}

type DiffRespVO struct {
	FilePath    string       `json:"filePath"`
	OldMode     string       `json:"oldMode"`
	Mode        string       `json:"mode"`
	IsSubModule bool         `json:"isSubModule"`
	FileType    string       `json:"fileType"`
	IsBinary    bool         `json:"isBinary"`
	RenameFrom  string       `json:"renameFrom"`
	RenameTo    string       `json:"renameTo"`
	CopyFrom    string       `json:"copyFrom"`
	CopyTo      string       `json:"copyTo"`
	Lines       []DiffLineVO `json:"lines"`
}

type DiffLineVO struct {
	Index   int    `json:"index"`
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

type CatFileReqVO struct {
	RepoPath  string `json:"repoPath"`
	CommitId  string `json:"commitId"`
	FileName  string `json:"fileName"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Direction string `json:"direction"`
}

type CatFileRespVO struct {
	ginutil.BaseResp
	Lines []DiffLineVO
}
