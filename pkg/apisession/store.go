package apisession

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/LeeZXin/zsf-utils/idutil"
	"strconv"
	"time"
	"zgit/standalone/modules/model/usermd"
)

const (
	SessionExpiry          = 2 * time.Hour
	RefreshSessionInterval = 10 * time.Minute
)

var (
	storeImpl = newMemStore()
)

type Session struct {
	SessionId string          `json:"sessionId"`
	UserInfo  usermd.UserInfo `json:"userInfo"`
	ExpireAt  int64           `json:"expireAt"`
}

type Store interface {
	GetBySessionId(string) (Session, bool, error)
	GetByAccount(string) (Session, bool, error)
	PutSession(Session) error
	DeleteByAccount(string) error
	DeleteBySessionId(string) error
	RefreshExpiry(string, int64) error
}

func GetStore() Store {
	return storeImpl
}

func GenSessionId() string {
	h := sha256.New()
	h.Write([]byte(idutil.RandomUuid() + strconv.FormatInt(time.Now().UnixNano(), 10)))
	return hex.EncodeToString(h.Sum(nil))
}
