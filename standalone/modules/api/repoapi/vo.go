package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/pkg/git"
)

type AllGitIgnoreTemplateListRespVO struct {
	ginutil.BaseResp
	Data []string `json:"data"`
}

type InitRepoReqVO struct {
	Name          string `json:"name"`
	Desc          string `json:"Desc"`
	RepoType      int    `json:"repoType"`
	CreateReadme  bool   `json:"createReadme"`
	ProjectId     string `json:"projectId"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReqVO struct {
	RepoId string `json:"repoId"`
}

type TreeRepoReqVO struct {
	RepoId  string `json:"repoId"`
	RefName string `json:"refName"`
	Dir     string `json:"dir"`
}

type EntriesRepoReqVO struct {
	RepoId  string `json:"repoId"`
	RefName string `json:"refName"`
	Dir     string `json:"dir"`
	Offset  int    `json:"offset"`
}

type ListRepoReqVO struct {
	ProjectId string `json:"projectId"`
}

type ListRepoRespVO struct {
	ginutil.BaseResp
	RepoList   []RepoVO `json:"repoList"`
	TotalCount int64    `json:"totalCount"`
	Cursor     int64    `json:"cursor"`
	Limit      int      `json:"limit"`
}

type CommitVO struct {
	Author        git.User `json:"author"`
	Committer     git.User `json:"committer"`
	AuthoredDate  string   `json:"authoredDate"`
	CommittedDate string   `json:"committedDate"`
	CommitMsg     string   `json:"commitMsg"`
	CommitId      string   `json:"commitId"`
	ShortId       string   `json:"shortId"`
}

type FileVO struct {
	Mode    string   `json:"mode"`
	RawPath string   `json:"rawPath"`
	Path    string   `json:"path"`
	Commit  CommitVO `json:"commit"`
}

type TreeVO struct {
	Files   []FileVO `json:"files"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
	HasMore bool     `json:"hasMore"`
}

type TreeRepoRespVO struct {
	ginutil.BaseResp
	IsEmpty      bool     `json:"isEmpty"`
	ReadmeText   string   `json:"readmeText"`
	RecentCommit CommitVO `json:"recentCommit"`
	Tree         TreeVO   `json:"tree"`
}

type RepoVO struct {
	RepoId    string `json:"repoId"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Author    string `json:"author"`
	ProjectId string `json:"projectId"`
	RepoType  string `json:"repoType"`
	IsEmpty   bool   `json:"isEmpty"`
	TotalSize int64  `json:"totalSize"`
	WikiSize  int64  `json:"wikiSize"`
	GitSize   int64  `json:"gitSize"`
	LfsSize   int64  `json:"lfsSize"`
	Created   string `json:"created"`
}

type CatFileReqVO struct {
	RepoId   string `json:"repoId"`
	RefName  string `json:"refName"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
}

type CatFileRespVO struct {
	ginutil.BaseResp
	Mode    string `json:"mode"`
	Content string `json:"content"`
}

type RepoTypeVO struct {
	Option int    `json:"option"`
	Name   string `json:"name"`
}

type AllTypeListRespVO struct {
	ginutil.BaseResp
	Data []RepoTypeVO `json:"data"`
}

type AllBranchesReqVO struct {
	RepoId string `json:"repoId"`
}

type AllBranchesRespVO struct {
	ginutil.BaseResp
	Data []string `json:"data"`
}

type AllTagsReqVO struct {
	RepoId string `json:"repoId"`
}

type AllTagsRespVO struct {
	ginutil.BaseResp
	Data []string `json:"data"`
}

type GcReqVO struct {
	RepoId string `json:"repoId"`
}

type PrepareMergeReqVO struct {
	RepoId string `json:"repoId"`
	Target string `json:"target"`
	Head   string `json:"head"`
}

type DiffFileReqVO struct {
	RepoId   string `json:"repoId"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FileName string `json:"fileName"`
}

type PrepareMergeRespVO struct {
	ginutil.BaseResp
	Target        string             `json:"target"`
	Head          string             `json:"head"`
	TargetCommit  CommitVO           `json:"targetCommit"`
	HeadCommit    CommitVO           `json:"headCommit"`
	Commits       []CommitVO         `json:"commits"`
	NumFiles      int                `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoVO `json:"diffNumsStats"`
	ConflictFiles []string           `json:"conflictFiles"`
	CanMerge      bool               `json:"canMerge"`
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

type DiffFileRespVO struct {
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

type ShowDiffTextContentReqVO struct {
	RepoId    string `json:"repoId"`
	CommitId  string `json:"commitId"`
	FileName  string `json:"fileName"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Direction string `json:"direction"`
}

type ShowDiffTextContentRespVO struct {
	ginutil.BaseResp
	Lines []DiffLineVO `json:"lines"`
}
