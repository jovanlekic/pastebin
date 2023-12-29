package db

import "context"

// Function to mark a key as used
func (dbObj *PostgresDB) MarkKeyAsUsed(ctx context.Context, key string) error {
	query := `
        UPDATE Keys
        SET used = true
        WHERE key = $1
    `

	_, err := dbObj.db.ExecContext(ctx, query, key)
	return err
}

// Function to check if a key is marked as used
func (dbObj *PostgresDB) IsKeyUsed(ctx context.Context, key string) (bool, error) {
	query := `
        SELECT used
        FROM Keys
        WHERE key = $1
    `

	var isUsed bool
	err := dbObj.db.QueryRowContext(ctx, query, key).Scan(&isUsed)
	if err != nil {
		return false, err
	}

	return isUsed, nil
}

// Function to get the first unused key
func (dbObj *PostgresDB) GetFirstUnusedKey(ctx context.Context) (string, error) {
	query := `
        SELECT key
        FROM Keys
        WHERE used = false
        LIMIT 1
    `

	var key string
	err := dbObj.db.QueryRowContext(ctx, query).Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// Function to get the first unused key and mark it as used
func (dbObj *PostgresDB) GetAndMarkFirstUnusedKey(ctx context.Context) (string, error) {
	key, err := dbObj.GetFirstUnusedKey(ctx)
	if err != nil {
		return "", err
	}

	// Mark the key as used
	err = dbObj.MarkKeyAsUsed(ctx, key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// Function to insert a key into the Keys table
func (dbObj *PostgresDB) InsertKeyIntoKeys(ctx context.Context, key string) error {
	query := `
        INSERT INTO Keys (id, key, used)
        VALUES (uuid_generate_v4(), $1, false)
    `

	_, err := dbObj.db.ExecContext(ctx, query, key)
	if err != nil {
		return err
	}
	return nil
}

// Function to fill the Keys table with all possible 6-character combinations of letters 'a', 'b', and 'c'
func (dbObj *PostgresDB) FillKeysTable(ctx context.Context) error {
	characters := "abc"

	tx, err := dbObj.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = dbObj.InsertMigration(ctx, "FillKeysTable")
	if err != nil {
		// Roll back the transaction if migration insertion fails
		_ = tx.Rollback()
		return err
	}

	defer func() {
		if err != nil {
			// Roll back the transaction if an error occurs during key insertion
			_ = tx.Rollback()
		}
	}()

	for i := 0; i < len(characters); i++ {
		for j := 0; j < len(characters); j++ {
			for k := 0; k < len(characters); k++ {
				for l := 0; l < len(characters); l++ {
					for m := 0; m < len(characters); m++ {
						for n := 0; n < len(characters); n++ {
							word := string(characters[i]) + string(characters[j]) + string(characters[k]) +
								string(characters[l]) + string(characters[m]) + string(characters[n])

							err := dbObj.InsertKeyIntoKeys(ctx, word)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	// Commit the transaction if all key insertions are successful
	err = tx.Commit()
	if err != nil {
		// Clear the migration if committing fails
		_ = dbObj.DeleteMigration(ctx, "FillKeysTable")
		return err
	}

	return nil
}
