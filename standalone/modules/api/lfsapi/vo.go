package lfsapi

import "time"

type BatchReqVO struct {
	Operation string `json:"operation"`
	// 没使用到
	Transfers []string    `json:"transfers,omitempty"`
	Ref       ReferenceVO `json:"ref,omitempty"`
	Objects   []PointerVO `json:"objects"`
	HashAlgo  string      `json:"hash_algo"`
}

// ReferenceVO contains a git reference.
// https://github.com/git-lfs/git-lfs/blob/main/docs/api/batch.md#ref-property
type ReferenceVO struct {
	Name string `json:"name"`
}

// PointerVO contains LFS pointer data
type PointerVO struct {
	Oid  string `json:"oid"`
	Size int64  `json:"size"`
}

// BatchRespVO contains multiple object metadata Representation structures
// for use with the batch API.
// https://github.com/git-lfs/git-lfs/blob/main/docs/api/batch.md#successful-responses
type BatchRespVO struct {
	Transfer string         `json:"transfer,omitempty"`
	Objects  []ObjectRespVO `json:"objects"`
}

// ObjectRespVO is object metadata as seen by clients of the LFS server.
type ObjectRespVO struct {
	PointerVO
	Actions map[string]LinkVO `json:"actions,omitempty"`
	Error   *ObjectErrVO      `json:"error,omitempty"`
}

// LinkVO provides a structure with information about how to access a object.
type LinkVO struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt *time.Time        `json:"expires_at,omitempty"`
}

// ObjectErrVO defines the JSON structure returned to the client in case of an error.
type ObjectErrVO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LockVO represent a lock
// for use with the locks API.
type LockVO struct {
	Id       string       `json:"id"`
	Path     string       `json:"path"`
	LockedAt time.Time    `json:"locked_at"`
	Owner    *LockOwnerVO `json:"owner"`
}

// LockOwnerVO represent a lock owner
// for use with the locks API.
type LockOwnerVO struct {
	Name string `json:"name"`
}

type ErrVO struct {
	Message       string `json:"message"`
	Documentation string `json:"documentation_url,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
}

type PostLockReqVO struct {
	Path string      `json:"path"`
	Ref  ReferenceVO `json:"ref"`
}

type PostLockRespVO struct {
	Lock LockVO `json:"lock"`
}

type ListLockReqVO struct {
	Path    string `json:"path" form:"path"`
	Id      string `json:"id" form:"id"`
	Cursor  string `json:"cursor" form:"cursor"`
	Limit   int    `json:"limit" form:"limit"`
	RefSpec string `json:"refspec" form:"refspec"`
}

type ListLockVerifyReqVO struct {
	Cursor string      `json:"cursor"`
	Limit  int         `json:"limit"`
	Ref    ReferenceVO `json:"ref"`
}

type ListLockRespVO struct {
	Locks []LockVO `json:"locks"`
	Next  string   `json:"next_cursor,omitempty"`
}

type UnlockReqVO struct {
	Force bool `json:"force"`
}

type UnlockRespVO struct {
	Lock LockVO `json:"lock"`
}

type ListLockVerifyRespVO struct {
	Ours   []LockVO `json:"ours"`
	Theirs []LockVO `json:"theirs"`
	Next   string   `json:"next_cursor,omitempty"`
}
