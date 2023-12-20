package clustermd

import (
	"encoding/json"
	"github.com/LeeZXin/zsf/logger"
	"os"
	"path/filepath"
	"zgit/setting"
)

var (
	nodes = make([]NodeInfo, 0)
)

func InitNodesConfig() {
	content, err := os.ReadFile(filepath.Join(setting.ResourcesDir(), "nodes.json"))
	if err != nil {
		logger.Logger.Error("read nodes.json failed")
		return
	}
	err = json.Unmarshal(content, &nodes)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	logger.Logger.Info("load nodes.json")
	logger.Logger.Info(string(content))
}

type NodeInfo struct {
	NodeId string `json:"nodeId"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}
