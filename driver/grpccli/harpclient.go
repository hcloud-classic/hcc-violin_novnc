package grpccli

import (
	"context"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	"hcc/violin-novnc/action/grpc/rpcharp"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

var harpconn *grpc.ClientConn

func initHarp() error {
	var err error
	addr := config.Harp.Address + ":" + strconv.FormatInt(config.Harp.Port, 10)
	logger.Logger.Println("Try connect to harp " + addr)
	harpconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Logger.Fatalf("Connect Harp failed: %v", err)
		return err
	}

	RC.harp = rpcharp.NewHarpClient(harpconn)
	logger.Logger.Println("GRPC connection to harp created")

	return nil
}

func cleanHarp() {
	harpconn.Close()
}

func (rc *RpcClient) GetServerIP(srvUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	r, err := rc.harp.GetSubnetByServer(ctx, &rpcharp.ReqGetSubnetByServer{ServerUUID: srvUUID})
	if err != nil {
		return "", err
	}

	srvIP := strings.TrimRight(r.Subnet.NetworkIP, "0")
	return srvIP + "1", nil
}
