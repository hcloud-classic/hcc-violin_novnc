package server

import (
	"net"
	"strconv"

	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"innogrid.com/hcloud-classic/pb"
)

type server struct {
	pb.UnimplementedNovncServer
}

var srv *grpc.Server

func InitGRPCServer() {
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(config.Grpc.Port, 10))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(config.Grpc.Port, 10) + "...")

	srv = grpc.NewServer()
	pb.RegisterNovncServer(srv, &server{})
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}

func CleanGRPCServer() {
	srv.Stop()
}
