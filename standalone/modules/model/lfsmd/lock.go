package lfsmd

import (
	"time"
)

const (
	LfsLockTableName = "lfs_lock"
)

type LfsLock struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  string    `json:"repoId"`
	Owner   string    `json:"owner"`
	Path    string    `json:"path" xorm:"TEXT"`
	Created time.Time `json:"created" xorm:"created"`
}

func (l LfsLock) TableName() string {
	return LfsLockTableName
}
