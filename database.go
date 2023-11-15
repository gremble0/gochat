package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string // NB: Not particularly safe
	DBName   string
}

func dbConnect(DBConf DBConfig) (*sql.DB, error)  {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		DBConf.Host, DBConf.Port, DBConf.User, DBConf.DBName)
	if DBConf.Password != "" {
		connStr += "password=" + DBConf.Password
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
