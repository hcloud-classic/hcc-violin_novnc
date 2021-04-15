package client

import (
	"sync"

	rpcharp "innogrid.com/hcloud-classic/pb"
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
