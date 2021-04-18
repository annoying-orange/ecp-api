package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/annoying-orange/ecp-api/graph"
)

const (
	defaultPort          = "8080"
	defaultMysqlHost     = "127.0.0.1:3306"
	defaultMysqlUser     = "root"
	defaultMysqlPassword = "passw0rd!"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = defaultMysqlHost
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = defaultMysqlUser
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		mysqlPassword = defaultMysqlPassword
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/ecp", mysqlUser, mysqlPassword, mysqlHost)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	log.Printf("Connected to %s", dsn)

	// Start service
	s, err := graph.NewResolver(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\nconnect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(s.Serve("/", port))
}
