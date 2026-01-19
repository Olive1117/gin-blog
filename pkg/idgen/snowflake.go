package idgen

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *snowflake.Node

func init() {
	nodeIDStr := os.Getenv("NODE_ID")
	if nodeIDStr == "" {
		node = initSnowflake(0)
		return
	}
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		if logger.L != nil {
			logger.L.Error("环境变量 NODE_ID 格式错误", zap.Error(err))
		} else {
			fmt.Printf("Warning: NODE_ID [%s] invalid, fallback to node 0\n", nodeIDStr)
		}
		node = initSnowflake(0)
	} else {
		node = initSnowflake(nodeID)
	}
}

func initSnowflake(id int64) *snowflake.Node {
	entity, err := snowflake.NewNode(id)
	if err != nil {
		panic("critical: snowflake node 0 failed to start:" + err.Error())
	}
	return entity
}

func NextID() int64 {
	return node.Generate().Int64()
}
