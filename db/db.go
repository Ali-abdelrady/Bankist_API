package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSource string) {
    var err error
    DB, err = sql.Open("postgres", dataSource)
    if err != nil {
        log.Fatal(err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Database not reachable:", err)
    }
    log.Println("Connected to the database.")
}
