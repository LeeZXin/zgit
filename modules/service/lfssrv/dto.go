package lfssrv

import (
	"io"
	"regexp"
	"time"
	"zgit/modules/model/lfsmd"
	"zgit/modules/model/repomd"
	"zgit/modules/model/usermd"
)

var (
	oidPattern = regexp.MustCompile(`^[a-f\d]{64}$`)
)

type LockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
}

type ListLockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Path     string
	Cursor   string
	Limit    int
	RefName  string
}

type ListLockRespDTO struct {
	LockList []lfsmd.LfsLock
	Next     string
}

type UnlockReqDTO struct {
	Repo     repomd.RepoInfo
	LockId   int64
	Force    bool
	Operator usermd.UserInfo
}

type PointerDTO struct {
	Oid  string
	Size int64
}

type VerifyReqDTO struct {
	PointerDTO
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
}

type DownloadReqDTO struct {
	Oid      string
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	FromByte int64
	ToByte   int64
}

type UploadReqDTO struct {
	Oid      string
	Size     int64
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Body     io.Reader
}

type DownloadRespDTO struct {
	io.ReadCloser
	FromByte int64
	ToByte   int64
	Length   int64
}

type BatchReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Objects  []PointerDTO
	IsUpload bool
	RefName  string
}

type LinkDTO struct {
	Href      string
	Header    map[string]string
	ExpiresAt *time.Time
}

// ObjectErrDTO defines the JSON structure returned to the client in case of an error.
type ObjectErrDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectDTO struct {
	PointerDTO
	Err error
}

type BatchRespDTO struct {
	ObjectList []ObjectDTO
}
