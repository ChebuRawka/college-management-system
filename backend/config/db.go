package config

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
    connStr := "host=db port=5432 user=admin password=password dbname=college sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to the database")
    return db
}