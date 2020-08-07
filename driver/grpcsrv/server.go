package grpcsrv

import (
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	rpcnovnc "hcc/violin-novnc/action/grpc/rpcviolin_novnc"

	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

type server struct {
	rpcnovnc.UnimplementedNovncServer
}

/*
func (s *server) CreateVNC(ctx context.Context, in *rpcnovnc.ReqNoVNC) (*rpcnovnc.ResNoVNC, error) {
	driver.RunnerGRPC(in.Vncs)
	return &rpcnovnc.ResNoVNC{Vncs: in.Vncs}, nil
}
*/

func (s *server) ControlVNC(ctx context.Context, in *rpcnovnc.ReqNoVNC) (*rpcnovnc.ResNoVNC, error) {
	for _, vnc := range in.Vncs {
		switch vnc.Action {
		case "CREATE":
			driver.VNCM.Create(vnc.Token, vnc.ServerUUID)
		case "DELETE":

		case "UPDATE":
			continue
		case "INFO":
			continue
		default:
			logger.Logger.Println("Undefined Action: " + vnc.Action)
		}

	}
	return &rpcnovnc.ResNoVNC{Vncs: in.Vncs}, nil
}

func InitGRPCServer() error {

	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(config.HTTP.Port, 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(config.HTTP.Port, 10) + "...")

	s := grpc.NewServer()
	rpcnovnc.RegisterNovncServer(s, &server{})
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func CleanGRPC() {

}
