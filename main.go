package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/Gravitalia/Harpocrate/database"
	"github.com/Gravitalia/Harpocrate/helpers"
	"github.com/Gravitalia/Harpocrate/model"
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
	id, err := helpers.GenerateRandomString(8)
	if err != nil {
		log.Printf("Cannot create random string: %v", err)

		return &proto.ReduceReponse{
			Id: "",
		}, nil
	}

	// Check if URL is valid
	if !regexp.MustCompile(`(?m)^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`).MatchString(in.Url) {
		return &proto.ReduceReponse{
			Id: "",
		}, nil
	}

	// Get HTML content, and check if page exists
	_, err = helpers.GetPageHTML(in.Url)
	if err != nil {
		return &proto.ReduceReponse{
			Id: "",
		}, nil
	}

	database.Session.Query(
		"INSERT INTO harpocrate.url ( id, original_url, author, analytics, phishing ) VALUES (?, ?, ?, ?, ?);",
		model.URL{
			Id:          id,
			OriginalUrl: in.Url,
			Author:      "",
			Analytics:   in.Opt.Number() == 1 || in.Opt.Number() == 3,
			Phishing:    -1.456,
		},
	)

	return &proto.ReduceReponse{
		Id: id,
	}, nil
}

func main() {
	fmt.Println(helpers.GenerateRandomString(8))
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
