package server

import (
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	rpcnovnc "hcc/violin-novnc/action/grpc/pb/rpcviolin_novnc"

	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
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
	var err error
	vnc := in.Vnc

	switch vnc.Action {
	case "CREATE":
		port, err = driver.VNCD.Create(vnc.ServerUUID)
		if err != nil {
			return nil, err
		}
	case "DELETE":
		err = driver.VNCD.Delete(vnc.ServerUUID)
		if err != nil {
			return nil, err
		}
		port = "Success"
	case "UPDATE":
	case "INFO":
	default:
		logger.Logger.Println("Undefined Action: " + vnc.Action)
	}

	return &rpcnovnc.ResControlVNC{Port: port}, nil
}

func InitGRPCServer() error {

	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(config.HTTP.Port, 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(config.HTTP.Port, 10) + "...")

	srv = grpc.NewServer()
	rpcnovnc.RegisterNovncServer(srv, &server{})
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func CleanGRPCServer() {
	srv.Stop()
}
