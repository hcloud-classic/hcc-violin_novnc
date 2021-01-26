package client

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
	errors "innogrid.com/hcloud-classic/hcc_errors"

	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	rpcnovnc "innogrid.com/hcloud-classic/pb"
)

var novncconn grpc.ClientConn

func initNovnc() *errors.HccError {
	addr := config.ViolinNoVnc.ServerAddress + ":" + strconv.FormatInt(config.ViolinNoVnc.ServerPort, 10)
	logger.Logger.Println("Try connect to violin-novnc " + addr)
	novncconn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		// Change Middleware Name
		return errors.NewHccError(errors.ViolinNoVNCGraphQLInitFail, err.Error())
	}

	RC.novnc = rpcnovnc.NewNovncClient(novncconn)
	logger.Logger.Println("GRPC connected to violin-novnc")

	return nil
}

func cleanNovnc() {
	novncconn.Close()
}

func (rc *RpcClient) ControlVNC(reqData map[string]interface{}) (interface{}, *errors.HccError) {
	//req data mapping
	var req rpcnovnc.ReqControlVNC
	req.Vnc = &rpcnovnc.VNC{
		ServerUUID: reqData["server_uuid"].(string),
		Action:     reqData["action"].(string),
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.ViolinNoVnc.RequestTimeoutMs)*time.Millisecond)
	defer cancel()

	r, err := rc.novnc.ControlVNC(ctx, &req)
	if err != nil {
		return nil, errors.NewHccError(errors.ViolinNoVNCGraphQLReceiveError, "ControlVNC : "+err.Error())
	}
	return r, nil
}
