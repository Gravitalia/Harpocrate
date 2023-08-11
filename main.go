package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

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

// Reduce defines route to reduce URL size, checks
func (s *server) Reduce(
	ctx context.Context,
	in *proto.ReduceRequest,
) (*proto.ReduceReponse, error) {
	id, err := helpers.GenerateRandomString(8)
	fmt.Println(id)
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
	content, err := helpers.GetPageHTML(in.Url)
	if err != nil {
		return &proto.ReduceReponse{
			Id: "",
		}, nil
	}

	database.Session.Query(
		"INSERT INTO harpocrate.url ( id, original_url, author, analytics, phishing ) VALUES (?, ?, ?, ?, ?);",
		id, in.Url, "", in.Opt.Number() == 1 || in.Opt.Number() == 3, helpers.CheckHTML(content),
	)

	return &proto.ReduceReponse{
		Id: id,
	}, nil
}

// generateHTMLPage parse templating file and then, return an html
// containing URL
func generateHTMLPage(w http.ResponseWriter, tmplFile string, data interface{}) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	// Create multiplexer
	httpMux := http.NewServeMux()

	go func() {
		if err := http.Serve(lis, httpMux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			return
		}

		var data model.Basics
		database.Session.Query(
			"SELECT original_url, phishing FROM harpocrate.url WHERE id = ?;",
			id,
		).Scan(&data.OriginalUrl, &data.Phishing)

		if data.Phishing > 0.5 {
			http.Redirect(w, r, data.OriginalUrl, http.StatusPermanentRedirect)
		} else {
			generateHTMLPage(w, "template.html", data)
		}
	})

	// Create gRPC server
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterHarpocrateServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
