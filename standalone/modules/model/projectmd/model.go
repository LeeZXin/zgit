package projectmd

import (
	"encoding/json"
	"time"
	"zgit/pkg/perm"
)

const (
	ProjectTableName          = "project"
	ProjectUserTableName      = "project_user"
	ProjectUserGroupTableName = "project_user_group"
)

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
	Id        int64  `json:"id" xorm:"pk autoincr"`
	ProjectId string `json:"projectId"`
	Account   string `json:"account"`
	// 关联用户组
	GroupId string    `json:"groupId"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*ProjectUser) TableName() string {
	return ProjectUserTableName
}

type ProjectUserGroup struct {
	Id int64 `json:"id" xorm:"pk autoincr"`
	// 组id
	GroupId string `json:"groupId"`
	// 项目id
	ProjectId string `json:"projectId"`
	// 名称
	Name string `json:"name"`
	// 权限json内容
	Perm string `json:"perm"`
	// 是否是管理员用户组
	IsAdmin bool `json:"isAdmin"`
	// 创建时间
	Created time.Time `json:"created" xorm:"created"`
	// 更新时间
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*ProjectUserGroup) TableName() string {
	return ProjectUserGroupTableName
}

func (p *ProjectUserGroup) GetPermDetail() (perm.Detail, error) {
	ret := perm.Detail{}
	err := json.Unmarshal([]byte(p.Perm), &ret)
	return ret, err
}
