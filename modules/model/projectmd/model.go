package projectmd

import "time"

const (
	ProjectTableName = "project"
)

type Project struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	ProjectId   string    `json:"projectId"`
	Name        string    `json:"name"`
	ProjectDesc string    `json:"projectDesc"`
	CorpId      string    `json:"corpId"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*Project) TableName() string {
	return ProjectTableName
}
