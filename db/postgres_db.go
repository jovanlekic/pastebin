package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB) (dbObj *PostgresDB) {
	dbObj = new(PostgresDB) // mora da se rezervise mem za obj
	dbObj.db = db
	return
}

func ConnectToPostgresDb(dbName, user, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=db port=5432 dbname=%s sslmode=disable", user, password, dbName)

	dbo, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = dbo.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to Postgres!->", dbName)

	driver, err := postgres.WithInstance(dbo, &postgres.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}
	fmt.Println("Success migration!")

	return dbo, nil
}

func DisconnectFromPostgresDb(client *sql.DB) {
	fmt.Print("Disconnected from Postgres!\n")
	client.Close()
}
