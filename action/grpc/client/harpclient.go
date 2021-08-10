package client

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	errors "innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"

	"hcc/violin-novnc/action/grpc/errconv"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

var harpConn *grpc.ClientConn

func initHarp() error {
	var err error

	addr := config.Harp.ServerAddress + ":" + strconv.FormatInt(config.Harp.ServerPort, 10)
	harpConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.harp = pb.NewHarpClient(harpConn)
	logger.Logger.Println("gRPC harp client ready")

	return nil
}

func closeHarp() {
	_ = harpConn.Close()
}

func pingHarp() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Harp.ServerAddress,
		strconv.FormatInt(config.Harp.ServerPort, 10)),
		time.Duration(config.Grpc.ClientPingTimeoutMs)*time.Millisecond)
	if err != nil {
		return false
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
		return true
	}

	return false
}

func checkHarp() {
	ticker := time.NewTicker(time.Duration(config.Grpc.ClientPingIntervalMs) * time.Millisecond)
	go func() {
		connOk := true
		for range ticker.C {
			pingOk := pingHarp()
			if pingOk {
				if !connOk {
					logger.Logger.Println("checkHarp(): Ping Ok! Resetting connection...")
					closeHarp()
					err := initHarp()
					if err != nil {
						logger.Logger.Println("checkHarp(): " + err.Error())
						continue
					}
					connOk = true
				}
			} else {
				if connOk {
					logger.Logger.Println("checkHarp(): Harp module seems dead. Pinging...")
				}
				connOk = false
			}
		}
	}()
}

// GetServerIP : Get the server's Leader Node IP address
func (rc *RPCClient) GetServerIP(srvUUID string) (string, *errors.HccErrorStack) {
	var srvIP string
	var errStack *errors.HccErrorStack

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	res, err := rc.harp.GetSubnetByServer(ctx, &pb.ReqGetSubnetByServer{ServerUUID: srvUUID})
	if err != nil {
		errStack = errors.NewHccErrorStack(errors.NewHccError(errors.ViolinNoVNCGrpcReceiveError, err.Error()))
		return "", errStack
	}
	if subnet := res.GetSubnet(); subnet != nil {
		srvIP = strings.TrimRight(subnet.GetNetworkIP(), "0") + "1"
	}
	if es := res.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(es)
	}
	return srvIP, errStack
}

// CreatePortForwarding : Create the AdaptiveIP Port Forwarding
func (rc *RPCClient) CreatePortForwarding(in *pb.ReqCreatePortForwarding) (*pb.ResCreatePortForwarding, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreatePortForwarding, err := rc.harp.CreatePortForwarding(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreatePortForwarding, nil
}

// DeletePortForwarding : Delete the AdaptiveIP Port Forwarding
func (rc *RPCClient) DeletePortForwarding(in *pb.ReqDeletePortForwarding) (*pb.ResDeletePortForwarding, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeletePortForwarding, err := rc.harp.DeletePortForwarding(ctx, in)
	if err != nil {
		return nil, err
	}

	return resDeletePortForwarding, nil
}
