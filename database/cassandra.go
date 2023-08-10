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
	PASSWORD string = "CASSANDRA_PASSWOWRD"
)

var Session *gocql.Session

// CreateSession returns an error. It allows to create a new session with basic authentification data (username and password) and then, set it to global variable Session.
func CreateSession() error {
	// Get adress from environnement
	address := os.Getenv(HOST)
	if address == "" {
		// Log that no host/address has been provided
		log.Println(
			"No CASSANDRA_HOST environment variable found. Using default address: 127.0.0.1",
		)
		// Set basic address if no one is given
		address = "127.0.0.1"
	}

	// Create new cluster with custom address
	cluster := gocql.NewCluster(address)

	// Set cluster consistency as LocalOne, the most economic one
	cluster.Consistency = gocql.LocalOne

	// Set cassandra username
	var username string
	// Get username from environnement
	if username = os.Getenv(USERNAME); username == "" {
		// If no username is provided, set it as default
		username = "cassandra"
	}

	// Set cassandra password
	var password string
	// Get password from environnement
	if password = os.Getenv(PASSWORD); password == "" {
		// If no password is provided, set it as default
		password = "cassandra"
	}

	// Add authentification parameters (if not defined, set cassandra as username, and cassandra as password)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	// Create session. If session do not works, exit
	if connection, err := cluster.CreateSession(); err != nil {
		log.Fatalf(
			"(CreateSession) Cannot create new Cassandra session: %v",
			err,
		)
		return err
	} else {
		Session = connection
	}

	return nil
}

// CreateTables allows to create tables of the keyspace. It returns nothing, but performs optimized action.
func CreateTables() {
	if err := Session.Query("CREATE TABLE IF NOT EXISTS harpocrate.users ( id TEXT, original_url TEXT, author TEXT, analytics BOOLEAN, PRIMARY KEY (id) ) WITH  compression = {'sstable_compression': 'LZ4Compressor'};").Exec(); err != nil {
		log.Printf(
			"(CreateTables) Create table harpocrate.users got error: %v",
			err,
		)
	}
}
