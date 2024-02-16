package kgs

import (
	"context"
	"database/sql"
	"pastebin/db"
)


type KGS interface {
	Check(key string) (string, error)
}

type kgs struct {
	ctx context.Context
	db  *db.PostgresDB
}


func initKgs(db *db.PostgresDB) error {
	doneMigration, err := db.CheckMigration(context.Background(), "FillKeysTable")
	if err != nil {
		return err
	}
	if !doneMigration {
		err := db.FillKeysTable(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func GetInstance(conn *sql.DB) *kgs {
	var instance *kgs
	instance = new(kgs)
	instance.ctx = context.Background()
	instance.db = db.NewPostgresDB(conn)
	initKgs(instance.db)
	
	return instance
}

func (k *kgs) Check(key string) (string, error) {
	if key == "" {
		res, err := k.db.GetAndMarkFirstUnusedKey(k.ctx)
		if err != nil {
			return "", err
		}
		return res, nil
	} else {
		isUsed, err := k.db.IsKeyUsed(k.ctx, key)
		if err != nil {
			return "", err
		}
		if isUsed {
			res, err := k.db.GetAndMarkFirstUnusedKey(k.ctx)
			if err != nil {
				return "", err
			}
			return res, nil
		}
		return key, nil
	}
}
