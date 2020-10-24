package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB = nil

func dbConnect() {
	if dbConn == nil {
		conn, err := sql.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME")+"?parseTime=true&loc="+os.Getenv("DB_LOCATION"))

		if err == nil {
			conn.SetConnMaxLifetime(time.Minute * 3)
			conn.SetMaxOpenConns(10)
			conn.SetMaxIdleConns(10)

			dbConn = conn
		} else {
			log.Println(err)
		}
	}
}
