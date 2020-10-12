package client

import (
	"sync"

	"hcc/violin-novnc/action/grpc/pb/rpcharp"
)

type RpcClient struct {
	harp rpcharp.HarpClient
}

var RC = &RpcClient{}

func InitGRPCClient() {
	var wg sync.WaitGroup

	wg.Add(1)
	go initHarp(&wg)

	wg.Wait()
}

func CleanGRPCClient() {
	cleanHarp()
}
