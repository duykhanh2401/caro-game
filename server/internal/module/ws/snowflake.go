package ws

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(1)
}

func GetRandomID() string {
	return node.Generate().String()
}
