package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	ctx       context.Context
	dbName    string
	tableName string
	db        *mongo.Collection
}

func NewDB(client *mongo.Client, dbName string, tableName string) (dbObj *DB) {
	dbObj = new(DB) // mora da se rezervise mem za obj
	dbObj.ctx = context.TODO()
	dbObj.dbName = dbName
	dbObj.tableName = tableName
	dbObj.db = client.Database(dbName).Collection(tableName)
	return
}

func (dbObj *DB) FindOne(filter interface{}, opts ...*options.FindOneOptions) (result bson.M, err error) {
	err = dbObj.db.FindOne(dbObj.ctx, filter, opts...).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Print("No document was found with this filter!\n", err)
		return
	}
	if err != nil {
		panic(err)
	}
	return
}
