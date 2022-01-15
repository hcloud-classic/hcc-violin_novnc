package client

import (
	"innogrid.com/hcloud-classic/pb"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	horn pb.HornClient
	harp pb.HarpClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {
	err := initHorn()
	if err != nil {
		return err
	}

	err = initHarp()
	if err != nil {
		return err
	}
	checkHarp()

	return nil
}

// End : Close connections of gRPC clients
func End() {
	closeHarp()
	closeHorn()
}
