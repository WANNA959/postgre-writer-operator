package postgre

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // no use on our end, just sets up the relevant drivers to be used by "database/sql" package
)

type PostgresDBClient struct {
	host         string
	port         int
	dbname       string
	user         string
	password     string
	DbConnection *sql.DB
}

var postgresDBClient *PostgresDBClient

func GetDbConnection() *sql.DB {
	return postgresDBClient.DbConnection
}

// If need be, setups a new DbConnection connection for a PostgresDBClient object (and returns it), else just returns the existing one functional one
// Not setting up a new connection blindly everything this method is called so as to maintain idempotency
func (pc *PostgresDBClient) SetupAndReturnDbConnection() (*sql.DB, error) {
	// if the .DbConnection attribute is nil, then we clearly need to setup a new connection as none exists presently
	// else if it exists but has some other issues (basically, some error occurred while "Ping"ing the DB with current connection), then setup a new connection
	// PS: expiry of connection won't be a concern because .Ping() not only checks the connection but resets the connections too automatically with new connections, if need be.
	// Ref: https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/database/sql/sql.go;l=868
	setupNewConnection := false
	if pc.DbConnection == nil {
		setupNewConnection = true
	} else if err := pc.DbConnection.Ping(); err != nil {
		setupNewConnection = true
	}
	fmt.Printf("SetupAndReturnDbConnection:%+v\n", setupNewConnection)

	if setupNewConnection {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			pc.host, pc.port, pc.user, pc.password, pc.dbname)
		newDbConnection, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			return nil, err
		}
		pc.DbConnection = newDbConnection
	}
	return pc.DbConnection, nil
}

// NewPostgresDBClient acts like a constructor to generating a new object of PostgresDBClient with a DbConnection as well
func NewPostgresDBClient(host string, port int, dbname, user, password string) (*PostgresDBClient, error) {

	postgresDBClient = &PostgresDBClient{
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
	}
	if _, err := postgresDBClient.SetupAndReturnDbConnection(); err != nil {
		return nil, err
	}
	return postgresDBClient, nil
}

// Close is used to gracefully close and wrap up the DbConnection associated with the PostgresDBClient object so as to avoid any associated memory leaks
func (pc *PostgresDBClient) Close() {
	pc.DbConnection.Close()
}
