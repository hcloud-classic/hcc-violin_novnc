package grpccli

import (
	"hcc/violin-novnc/action/grpc/rpcharp"
)

type RpcClient struct {
	harp rpcharp.HarpClient
}

var RC = &RpcClient{}

func InitGRPCClient() {
	go initHarp()
}
