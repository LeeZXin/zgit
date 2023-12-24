package projectmd

import (
	"github.com/LeeZXin/zsf-utils/idutil"
	"time"
)

const (
	ProjectTableName     = "project"
	ProjectUserTableName = "project_user"
)

func GenProjectId() string {
	return idutil.RandomUuid()
}

type Project struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	ProjectId   string    `json:"projectId"`
	Name        string    `json:"name"`
	ProjectDesc string    `json:"projectDesc"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*Project) TableName() string {
	return ProjectTableName
}

type ProjectUser struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	ProjectId string    `json:"projectId"`
	Account   string    `json:"account"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*ProjectUser) TableName() string {
	return ProjectUserTableName
}
