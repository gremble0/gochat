package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const TimeFormat = "2006-01-02 15:04:05"

// DBConfig contains the necessary configuration for connecting to a database
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string // NB: Not particularly safe
	DBName   string
}

// GochatDB is a wrapper struct for easier expansion of the sql.DB struct
type GochatDB struct {
	DB *sql.DB
}

// dbConnect connects to a database described by a database config.
// Returns error if unable to open and ping the database
func dbConnect(dbc DBConfig) (*GochatDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.DBName)
	if dbc.Password != "" {
		connStr += "password=" + dbc.Password
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &GochatDB{DB: db}, nil
}

// LogConnection inserts a new connection into the database
func (gcdb *GochatDB) LogConnection(message Message) error {
	cmd := `INSERT INTO users(username, remote_addr, registered) VALUES($1, $2, $3)`

	_, err := gcdb.DB.Exec(cmd,
		message.Sender.Username,
		message.Sender.Conn.RemoteAddr().String(),
		time.Now().Format(TimeFormat),
	)

	return err
}

// LogMessage inserts a new message from an existing connection into the databse
func (gcdb *GochatDB) LogMessage(message Message) error {
	cmd := `INSERT INTO messages(message, sender, sender_addr, sent) VALUES($1, $2, $3, $4)`

	_, err := gcdb.DB.Exec(cmd,
		message.Text,
		message.Sender.Username,
		message.Sender.Conn.RemoteAddr().String(),
		time.Now().Format(TimeFormat),
	)

	return err
}
