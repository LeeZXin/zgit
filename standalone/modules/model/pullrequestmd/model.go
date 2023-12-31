package pullrequestmd

import (
	"time"
	"zgit/pkg/i18n"
)

const (
	PullRequestTableName = "pull_request"
)

type PrStatus int

const (
	PrOpenStatus PrStatus = iota
	PrClosedStatus
	PrMergedStatus
)

func (s PrStatus) Int() int {
	return int(s)
}

func (s PrStatus) Readable() string {
	switch s {
	case PrOpenStatus:
		return i18n.GetByKey(i18n.PullRequestOpenStatus)
	case PrClosedStatus:
		return i18n.GetByKey(i18n.PullRequestClosedStatus)
	case PrMergedStatus:
		return i18n.GetByKey(i18n.PullRequestMergedStatus)
	default:
		return i18n.GetByKey(i18n.PullRequestUnknownStatus)
	}
}

type PullRequest struct {
	Id             int64     `json:"id" xorm:"pk autoincr"`
	PrId           string    `json:"prId"`
	RepoId         string    `json:"repoId"`
	Target         string    `json:"target"`
	TargetCommitId string    `json:"targetCommitId"`
	Head           string    `json:"head"`
	HeadCommitId   string    `json:"headCommitId"`
	PrStatus       int       `json:"prStatus"`
	CreateBy       string    `json:"createBy"`
	Created        time.Time `json:"created" xorm:"created"`
	Updated        time.Time `json:"updated" xorm:"updated"`
}

func (*PullRequest) TableName() string {
	return PullRequestTableName
}
