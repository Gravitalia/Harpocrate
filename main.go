package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Gravitalia/Harpocrate/proto"
	"google.golang.org/grpc"
)

// server struct defines basic Harpocrate server
type server struct {
	proto.UnimplementedHarpocrateServer
}

// Upload defines route to reduce URL size, checks
func (s *server) Upload(
	ctx context.Context,
	in *proto.ReduceRequest,
) (*proto.ReduceReponse, error) {
	return &proto.ReduceReponse{
		Id: "1",
	}, nil
}

func main() {
	// Get port from environnement
	port := os.Getenv("PORT")
	if port == "" {
		// If port is not in environnement, set it to 5000
		port = "5000"
	}

	// Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Log used port in console
	log.Printf("Listening on port: %s\n", port)

	// Set maximum message size
	var opts []grpc.ServerOption
	maxMsgSize := 50 * 1024 // 50 KiB
	opts = append(
		opts,
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	// Create gRPC server
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterHarpocrateServer(grpcServer, &server{})

	// Listen gRPC requests
	grpcServer.Serve(lis)
}
