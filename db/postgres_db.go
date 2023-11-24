package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectToPostgresDb() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	connStr := "user=jovanadragutinovic host=localhost port=5432 dbname=mydb sslmode=disable"

	dbo, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = dbo.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return dbo
}

func DisconnectFromPostgresDb(client *sql.DB) {
	fmt.Print("Disconnected!\n")
	client.Close()
}
