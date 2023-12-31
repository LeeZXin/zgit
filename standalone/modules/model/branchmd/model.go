package branchmd

import "time"

const (
	ProtectedBranchTableName = "protected_branch"
)

type ProtectedBranch struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Branch  string    `json:"branch"`
	RepoId  string    `json:"repoId"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*ProtectedBranch) TableName() string {
	return ProtectedBranchTableName
}
