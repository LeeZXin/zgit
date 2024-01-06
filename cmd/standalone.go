package cmd

import (
	"errors"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
	"regexp"
	"zgit/pkg/git"
	"zgit/setting"
	"zgit/standalone/modules/api/branchapi"
	"zgit/standalone/modules/api/cfgapi"
	"zgit/standalone/modules/api/hookapi"
	"zgit/standalone/modules/api/lfsapi"
	"zgit/standalone/modules/api/projectapi"
	"zgit/standalone/modules/api/pullrequestapi"
	"zgit/standalone/modules/api/repoapi"
	"zgit/standalone/modules/api/sshkeyapi"
	"zgit/standalone/modules/api/userapi"
	"zgit/standalone/modules/service/cfgsrv"
	"zgit/standalone/sshserv"
)

var Standalone = &cli.Command{
	Name:        "standalone",
	Usage:       "This command starts zgit standalone server",
	Description: "zgit",
	Action:      runStandalone,
}

var (
	validCorpIdPattern = regexp.MustCompile("^\\w{1,32}$")
)

func runStandalone(*cli.Context) error {
	// 检查corpId配置
	if !validCorpIdPattern.MatchString(setting.StandaloneCorpId()) {
		return errors.New("invalid standalone corpId config")
	}
	logger.Logger.Info("zgit works on standalone mode")
	// 初始化系统配置
	cfgsrv.InitSysCfg()
	// 初始化ssh服务
	sshserv.InitSsh()
	//
	git.InitGit()
	// 初始化api
	lfsapi.InitApi()
	// webhook
	hookapi.InitApi()
	// 用户
	userapi.InitApi()
	// 仓库api
	repoapi.InitApi()
	// ssh公钥
	sshkeyapi.InitApi()
	// 项目
	projectapi.InitApi()
	// 合并请求
	pullrequestapi.InitApi()
	// 分支
	branchapi.InitApi()
	// 系统配置api
	cfgapi.InitApi()
	starter.Run()
	return nil
}
