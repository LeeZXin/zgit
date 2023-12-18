package lfsmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"time"
)

var (
	store = hashmap.NewConcurrentHashMap[string, MetaObject]()
)

type MetaObject struct {
	Id      int64     `json:"id"`
	RepoId  string    `json:"repoId"`
	Oid     string    `json:"oid"`
	Size    int64     `json:"size"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*MetaObject) TableName() string {
	return "z_lfs_meta"
}

func GetMetaObjectByOid(ctx context.Context, oid string) (MetaObject, bool, error) {
	object, b := store.Get(oid)
	return object, b, nil
}

func InsertMetaObject(object MetaObject) error {
	store.Put(object.Oid, object)
	return nil
}
