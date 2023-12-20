package clustermd

import (
	"context"
)

func GetByNodeId(ctx context.Context, id string) (NodeInfo, bool, error) {
	for _, node := range nodes {
		if node.NodeId == id {
			return node, true, nil
		}
	}
	return NodeInfo{}, false, nil
}
