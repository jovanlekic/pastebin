package db

import (
	"context"
	"database/sql"
	"fmt"
	"pastebin/db/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
	connStr := fmt.Sprintf("user=%s password=%s host=localhost port=5432 dbname=%s sslmode=disable", user, password, dbName)

	dbo, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = dbo.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to Postgres!")
	return dbo, nil
}

func DisconnectFromPostgresDb(client *sql.DB) {
	fmt.Print("Disconnected from Postgres!\n")
	client.Close()
}

// CREATE
func (dbObj *PostgresDB) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	query := `
		INSERT INTO Users (user_id, name, password, pasteNum, dev_key, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id
	`

	var createdUserID uuid.UUID
	err := dbObj.db.QueryRowContext(ctx, query, user.UserID, user.Name, user.Password, user.PasteNum, user.DevKey, user.Email).Scan(&createdUserID)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Println("User created successfully")
	return createdUserID, nil
}

func (dbObj *PostgresDB) CreateObject(ctx context.Context, obj *models.Object) error {
	query := `
		INSERT INTO Object (paste_key, dev_key, message_id)
		VALUES ($1, $2, $3)
	`

	_, err := dbObj.db.ExecContext(ctx, query, obj.PasteKey, obj.DevKey, obj.MessageID)
	return err
}

// READ
func (dbObj *PostgresDB) ReadUserById(ctx context.Context, userID uuid.UUID) (models.User, error) {
	var user models.User
	err := dbObj.db.QueryRowContext(ctx, "SELECT user_id, name, password, pasteNum, dev_key, email FROM Users WHERE user_id = $1", userID).
		Scan(&user.UserID, &user.Name, &user.Password, &user.PasteNum, &user.DevKey, &user.Email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (dbObj *PostgresDB) ReadUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := dbObj.db.QueryRowContext(ctx, "SELECT user_id, name, password, pasteNum, dev_key, email FROM Users WHERE name = $1", username).
		Scan(&user.UserID, &user.Name, &user.Password, &user.PasteNum, &user.DevKey, &user.Email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (dbObj *PostgresDB) ReadObject(ctx context.Context, pasteKey, devKey string) (*models.Object, error) {
	var obj models.Object
	query := `
		SELECT paste_key, dev_key, message_id
		FROM Object
		WHERE paste_key = $1 AND dev_key = $2
	`

	err := dbObj.db.QueryRowContext(ctx, query, pasteKey, devKey).Scan(&obj.PasteKey, &obj.DevKey, &obj.MessageID)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// UPDATE
func (dbObj *PostgresDB) UpdateUser(ctx context.Context, userID uuid.UUID, updatedUser models.User) error {
	_, err := dbObj.db.ExecContext(ctx, "UPDATE Users SET name=$1, password=$2, pasteNum=$3, dev_key=$4, email=$5 WHERE user_id=$6",
		updatedUser.Name, updatedUser.Password, updatedUser.PasteNum, updatedUser.DevKey, updatedUser.Email, userID)
	if err != nil {
		return err
	}
	fmt.Println("User updated successfully")
	return nil
}

func (dbObj *PostgresDB) UpdateObject(ctx context.Context, obj *models.Object) error {
	query := `
		UPDATE Object
		SET message_id = $1
		WHERE paste_key = $2 AND dev_key = $3
	`

	_, err := dbObj.db.ExecContext(ctx, query, obj.MessageID, obj.PasteKey, obj.DevKey)
	return err
}

// DELETE
func (dbObj *PostgresDB) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := dbObj.db.ExecContext(ctx, "DELETE FROM Users WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	fmt.Println("User deleted successfully")
	return nil
}

func (dbObj *PostgresDB) DeleteObject(ctx context.Context, pasteKey, devKey string) error {
	query := `
		DELETE FROM Object
		WHERE paste_key = $1 AND dev_key = $2
	`

	_, err := dbObj.db.ExecContext(ctx, query, pasteKey, devKey)
	return err
}
