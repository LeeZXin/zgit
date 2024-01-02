package syscfgmd

import "time"

const (
	SysCfgTableName = "sys_cfg"
)

type SysCfg struct {
	Id      int64     `xorm:"pk autoincr"`
	CfgKey  string    `json:"cfgKey"`
	Content string    `json:"content"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*SysCfg) TableName() string {
	return SysCfgTableName
}
