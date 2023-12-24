package clustermd

import "time"

const (
	ClusterNodeTableName = "cluster_node"
)

type NodeInfo struct {
	NodeId string `json:"nodeId"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type ClusterNode struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	NodeId   string    `json:"nodeId"`
	Name     string    `json:"name"`
	NodeDesc string    `json:"nodeDesc"`
	NodeHost string    `json:"nodeHost"`
	NodePort int       `json:"nodePort"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}

func (*ClusterNode) TableName() string {
	return ClusterNodeTableName
}

func (n *ClusterNode) ToNodeInfo() NodeInfo {
	return NodeInfo{
		NodeId: n.NodeId,
		Name:   n.Name,
		Host:   n.NodeHost,
		Port:   n.NodePort,
	}
}
