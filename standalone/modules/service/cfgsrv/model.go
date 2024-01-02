package cfgsrv

import "encoding/json"

type SysCfg struct {
	// 禁用用户注册功能
	DisableSelfRegisterUser bool `json:"disableSelfRegisterUser"`
	// 允许用户自行创建项目
	AllowUserCreateProject bool `json:"allowUserCreateProject"`
}

func (c *SysCfg) Key() string {
	return "sys_cfg"
}

func (c *SysCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *SysCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}
