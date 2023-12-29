package db

import (
	"context"
	"fmt"
)

// Function to insert a migration record
func (dbObj *PostgresDB) InsertMigration(ctx context.Context, migrationName string) error {
	query := `
        INSERT INTO Migration (id, migration)
        VALUES (uuid_generate_v4(), $1)
    `

	_, err := dbObj.db.ExecContext(ctx, query, migrationName)
	if err != nil {
		return err
	}

	fmt.Printf("Migration '%s' inserted successfully\n", migrationName)
	return nil
}

// Function to delete a migration record
func (dbObj *PostgresDB) DeleteMigration(ctx context.Context, migrationName string) error {
	query := `
        DELETE FROM Migration
        WHERE migration = $1
    `

	_, err := dbObj.db.ExecContext(ctx, query, migrationName)
	if err != nil {
		return err
	}

	fmt.Printf("Migration '%s' deleted successfully\n", migrationName)
	return nil
}

func (dbObj *PostgresDB) CheckMigration(ctx context.Context, migrationName string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM Migration
            WHERE migration = $1
        )
    `

	var exists bool
	err := dbObj.db.QueryRowContext(ctx, query, migrationName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
