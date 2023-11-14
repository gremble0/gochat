package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string // NB: Not particularly safe
	DBName   string
}

func dbConnect(DBConf DBConfig) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBConf.Host, DBConf.Port, DBConf.User, DBConf.Password, DBConf.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to database\n")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
