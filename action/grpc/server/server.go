package server

import (
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hcc/violin-novnc/action/grpc/errconv"
	rpcnovnc "hcc/violin-novnc/action/grpc/pb/rpcviolin_novnc"

	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/errors"
	"hcc/violin-novnc/lib/logger"
)

type server struct {
	rpcnovnc.UnimplementedNovncServer
}

var srv *grpc.Server

/*
func (s *server) CreateVNC(ctx context.Context, in *rpcnovnc.ReqNoVNC) (*rpcnovnc.ResNoVNC, error) {
	driver.RunnerGRPC(in.Vncs)
	return &rpcnovnc.ResNoVNC{Vncs: in.Vncs}, nil
}
*/

func (s *server) ControlVNC(ctx context.Context, in *rpcnovnc.ReqControlVNC) (*rpcnovnc.ResControlVNC, error) {
	var port string
	var errStack *errors.HccErrorStack = nil
	var result rpcnovnc.ResControlVNC

	vnc := in.Vnc

	switch vnc.Action {
	case "CREATE":
		port, errStack = driver.VNCD.Create(vnc.ServerUUID)
		if errStack != nil {
			result.HccErrorStack = errconv.HccStackToGrpc(errStack)
			return &result, nil
		}
	case "DELETE":
		errStack = driver.VNCD.Delete(vnc.ServerUUID)
		if errStack != nil {
			result.HccErrorStack = errconv.HccStackToGrpc(errStack)
			return &result, nil
		}
		port = "Success"
	case "UPDATE":
	case "INFO":
	default:
		logger.Logger.Println("Undefined Action: " + vnc.Action)
		errStack = errors.NewHccErrorStack(errors.NewHccError(errors.ViolinNoVNCGrpcOperationFail, "Undefined Action("+vnc.Action+")"))
		result.HccErrorStack = errconv.HccStackToGrpc(errStack)
	}

	result.Port = port

	return &result, nil
}

func InitGRPCServer() {

	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(config.HTTP.Port, 10))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(config.HTTP.Port, 10) + "...")

	srv = grpc.NewServer()
	rpcnovnc.RegisterNovncServer(srv, &server{})
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}

func CleanGRPCServer() {
	srv.Stop()
}
