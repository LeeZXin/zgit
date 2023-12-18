package lfsmd

import (
	"context"
	"time"
)

type LfsLock struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  string    `json:"repoId"`
	OwnerId string    `json:"ownerId"`
	Path    string    `json:"path" xorm:"TEXT"`
	Created time.Time `json:"created" xorm:"created"`
}

func (l LfsLock) TableName() string {
	return "z_lfs_lock"
}

func InsertLock(ctx context.Context, lock *LfsLock) error {
	return nil
}

func FindLockByPath(ctx context.Context, path string) (LfsLock, bool, error) {
	return LfsLock{}, true, nil
}
