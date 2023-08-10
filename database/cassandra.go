package database

import (
	"log"
	"os"

	"github.com/gocql/gocql"
)

// Set environnement variable to avoid magic strings
const (
	HOST     string = "CASSANDRA_HOST"
	USERNAME string = "CASSANDRA_USERNAME"
	PASSWORD string = "CASSANDRA_PASSWORD"
)

var Session *gocql.Session

// CreateSession allows to create a new session with basicauthentification data (username and password)
// and then, set it to global variable Session. It returns an error in case of mis-connection.
func CreateSession() error {
	address := os.Getenv(HOST)
	if address == "" {
		// Set basic address if no one is given and log it
		log.Println(
			"No CASSANDRA_HOST environment variable found. Using default address: 127.0.0.1",
		)
		address = "127.0.0.1"
	}

	cluster := gocql.NewCluster(address)

	// Set cluster consistency as LocalOne,
	// the most economic one because it only needs to be
	// accepted by one replica node in the local datacenter
	cluster.Consistency = gocql.LocalOne

	// Obtains Cassandra credentials from environment variables.
	// In case they are absent, we use `cassandra` as both login and password.
	var username string
	if username = os.Getenv(USERNAME); username == "" {
		username = "cassandra"
	}

	var password string
	if password = os.Getenv(PASSWORD); password == "" {
		password = "cassandra"
	}

	// Add authentification parameters
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	// Create session. If session do not works, exit and close server
	if connection, err := cluster.CreateSession(); err != nil {
		return err
	} else {
		Session = connection
	}

	return nil
}

// CreateTables allows to create tables of the keyspace.
// It returns nothing, but performs optimized action.
func CreateTables() {
	if err := Session.Query("CREATE TABLE IF NOT EXISTS harpocrate.url ( id TEXT, original_url TEXT, author TEXT, analytics BOOLEAN, PRIMARY KEY (id) ) WITH  compression = {'sstable_compression': 'LZ4Compressor'};").Exec(); err != nil {
		log.Printf(
			"(CreateTables) Create table harpocrate.url got error: %v",
			err,
		)
	}
}
