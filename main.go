package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Gravitalia/Harpocrate/database"
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
	// Init database session, if session not initizalied,
	// Exit with code 1, and print error message
	if err := database.CreateSession(); err != nil {
		log.Fatalf(
			"Cannot create new Cassandra session: %v",
			err,
		)
	}

	database.CreateTables()

	// Get port from environnement, if no one is set, take 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

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
	grpcServer.Serve(lis)
}
