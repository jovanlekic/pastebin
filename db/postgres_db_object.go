package db

import (
	"context"
	"pastebin/db/models"
)

// CREATE
func (dbObj *PostgresDB) CreateObject(ctx context.Context, obj *models.Object) error {
	query := `
		INSERT INTO Object (paste_key, dev_key, message_id)
		VALUES ($1, $2, $3)
	`

	_, err := dbObj.db.ExecContext(ctx, query, obj.PasteKey, obj.DevKey, obj.MessageID)
	return err
}

// READ all objects with a certain devKey
func (dbObj *PostgresDB) ReadObjectsByDevKey(ctx context.Context, devKey string) ([]models.Object, error) {
	var objects []models.Object
	query := `
		SELECT paste_key, dev_key, message_id
		FROM Object
		WHERE dev_key = $1
	`

	rows, err := dbObj.db.QueryContext(ctx, query, devKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj models.Object
		if err := rows.Scan(&obj.PasteKey, &obj.DevKey, &obj.MessageID); err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

// READ Object
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

// READ
func (dbObj *PostgresDB) ReadObjectWithoutDevKey(ctx context.Context, pasteKey string) (*models.Object, error) {
	var obj models.Object
	query := `
		SELECT paste_key, dev_key, message_id
		FROM Object
		WHERE paste_key = $1
	`

	err := dbObj.db.QueryRowContext(ctx, query, pasteKey).Scan(&obj.PasteKey, &obj.DevKey, &obj.MessageID)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// UPDATE
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
func (dbObj *PostgresDB) DeleteObject(ctx context.Context, pasteKey, devKey string) error {
	query := `
		DELETE FROM Object
		WHERE paste_key = $1 AND dev_key = $2
	`

	_, err := dbObj.db.ExecContext(ctx, query, pasteKey, devKey)
	return err
}
