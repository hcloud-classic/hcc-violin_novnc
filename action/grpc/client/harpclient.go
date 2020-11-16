package client

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	"hcc/violin-novnc/action/grpc/errconv"
	"hcc/violin-novnc/action/grpc/pb/rpcharp"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/errors"
	"hcc/violin-novnc/lib/logger"
)

var harpconn *grpc.ClientConn

func initHarp(wg *sync.WaitGroup) *errors.HccError {
	var err error
	addr := config.Harp.Address + ":" + strconv.FormatInt(config.Harp.Port, 10)
	logger.Logger.Println("Try connect to harp " + addr)
	harpconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.NewHccError(errors.ViolinNoVNCGrpcConnectionFail, "harp : "+err.Error())
	}

	RC.harp = rpcharp.NewHarpClient(harpconn)
	logger.Logger.Println("GRPC connection to harp created")

	wg.Done()
	return nil
}

func cleanHarp() {
	harpconn.Close()
}

func (rc *RpcClient) GetServerIP(srvUUID string) (string, *errors.HccErrorStack) {
	var srvIP string
	var errStack *errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	res, err := rc.harp.GetSubnetByServer(ctx, &rpcharp.ReqGetSubnetByServer{ServerUUID: srvUUID})
	if err != nil {
		errStack = errors.NewHccErrorStack(errors.NewHccError(errors.ViolinNoVNCGrpcReceiveError, err.Error()))
		return "", errStack
	}
	if subnet := res.GetSubnet(); subnet != nil {
		srvIP = strings.TrimRight(subnet.GetNetworkIP(), "0") + "1"
	}
	if es := res.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(&es)
	}
	return srvIP, errStack
}
