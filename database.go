package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const TimeFormat = "2006-01-02 15:04:05"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string // NB: Not particularly safe
	DBName   string
}

type GochatDB struct {
	DB *sql.DB
}

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

func (gcdb *GochatDB) LogConnection(message Message) error {
	cmd := `INSERT INTO users(username, remote_addr, registered) VALUES($1, $2, $3)`

	_, err := gcdb.DB.Exec(cmd,
		message.Sender.Username,
		message.Sender.Conn.RemoteAddr().String(),
		time.Now().Format(TimeFormat),
	)

	return err
}

func (gcdb *GochatDB) LogMessage(message Message) error {
	cmd := `INSERT INTO messages(message, sent) VALUES($1, $2)`

	_, err := gcdb.DB.Exec(cmd,
		message.Text,
		time.Now().Format(TimeFormat),
	)

	return err
}
