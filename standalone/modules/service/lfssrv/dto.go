package lfssrv

import (
	"io"
	"regexp"
	"time"
	"zgit/standalone/modules/model/lfsmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

var (
	oidPattern = regexp.MustCompile(`^[a-f\d]{64}$`)
)

type LockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Path     string
}

func (r *LockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if r.Path == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type ListLockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Path     string
	Cursor   string
	Limit    int
	RefName  string
}

func (r *ListLockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
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

func (r *UnlockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type PointerDTO struct {
	Oid  string
	Size int64
}

func (p *PointerDTO) IsValid() error {
	if !oidPattern.MatchString(p.Oid) {
		return util.InvalidArgsError()
	}
	if p.Size < 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type VerifyReqDTO struct {
	PointerDTO
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
}

func (r *VerifyReqDTO) IsValid() error {
	if err := r.PointerDTO.IsValid(); err != nil {
		return util.InvalidArgsError()
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type DownloadReqDTO struct {
	Oid      string
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	FromByte int64
	ToByte   int64
}

func (r *DownloadReqDTO) IsValid() error {
	if !oidPattern.MatchString(r.Oid) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if r.FromByte < 0 || r.ToByte < 0 {
		return util.InvalidArgsError()
	}
	return nil
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

func (r *BatchReqDTO) IsValid() error {
	for _, obj := range r.Objects {
		if err := obj.IsValid(); err != nil {
			return err
		}
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
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

func validateRepo(repo repomd.RepoInfo) bool {
	return repo.RepoId != ""
}
