package client

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"

	rpcnovnc "hcc/violin-novnc/action/grpc/pb/rpcviolin_novnc"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

var novncconn grpc.ClientConn

func initNovnc() error {
	addr := config.ViolinNoVnc.ServerAddress + ":" + strconv.FormatInt(config.ViolinNoVnc.ServerPort, 10)
	logger.Logger.Println("Try connect to violin-novnc " + addr)
	novncconn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Logger.Fatalf("Connect Violin-Novnc failed: %v", err)
		return err
	}

	RC.novnc = rpcnovnc.NewNovncClient(novncconn)
	logger.Logger.Println("GRPC connected to violin-novnc")

	return nil
}

func cleanNovnc() {
	novncconn.Close()
}

func (rc *RpcClient) ControlVNC(reqData map[string]interface{}) (interface{}, error) {
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
		return nil, err
	}
	logger.Logger.Println(r)

	return r, nil
}
