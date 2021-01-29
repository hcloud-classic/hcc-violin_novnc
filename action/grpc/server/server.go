package server

import (
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	errors "innogrid.com/hcloud-classic/hcc_errors"
	rpcnovnc "innogrid.com/hcloud-classic/pb"

	"hcc/violin-novnc/action/grpc/errconv"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/model"
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
	var vncInfo model.Vnc
	var errStack *errors.HccErrorStack = nil
	var result rpcnovnc.ResControlVNC

	vnc := in.GetVnc()
	vncInfo.ServerUUID = vnc.GetServerUUID()

	switch vnc.GetAction() {
	case "CREATE":
		errStack = driver.VNCD.Create(&vncInfo)
		if errStack != nil {
			result.HccErrorStack = errconv.HccStackToGrpc(errStack)
			return &result, nil
		}

	case "DELETE":
		errStack = driver.VNCD.Delete(&vncInfo)
		if errStack != nil {
			result.HccErrorStack = errconv.HccStackToGrpc(errStack)
			return &result, nil
		}
		vncInfo.WebSocket = "Success"

	case "UPDATE":
	case "INFO":
	default:
		logger.Logger.Println("Undefined Action: " + vnc.GetAction())
		errStack = errors.NewHccErrorStack(errors.NewHccError(
			errors.ViolinNoVNCGrpcOperationFail,
			"Undefined Action("+vnc.GetAction()+")"))
		result.HccErrorStack = errconv.HccStackToGrpc(errStack)
	}

	result.Port = vncInfo.WebSocket

	return &result, nil
}

func InitGRPCServer() {

	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(config.Grpc.Port, 10))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(config.Grpc.Port, 10) + "...")

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
