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

	// Set cluster consistency as LocalOne, the most economic one
	cluster.Consistency = gocql.LocalOne

	// Set cassandra username
	var username string
	if username = os.Getenv(USERNAME); username == "" {
		// If no username is provided, set it as default
		username = "cassandra"
	}

	// Set cassandra password
	var password string
	if password = os.Getenv(PASSWORD); password == "" {
		// If no password is provided, set it as default
		password = "cassandra"
	}

	// Add authentification parameters (if not defined,
	// set cassandra as username, and cassandra as password)
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
	// Create url table with id (random string), original_url,
	// author (user vanity), whether or not analysis is enabled.
	if err := Session.Query("CREATE TABLE IF NOT EXISTS harpocrate.url ( id TEXT, original_url TEXT, author TEXT, analytics BOOLEAN, PRIMARY KEY (id) ) WITH  compression = {'sstable_compression': 'LZ4Compressor'};").Exec(); err != nil {
		// If table haven't been created (got an error while creating it)
		// log error, and warn user
		log.Printf(
			"(CreateTables) Create table harpocrate.url got error: %v",
			err,
		)
	}
}
