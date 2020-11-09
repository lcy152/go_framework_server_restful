package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddAppConfig(ctx context.Context, r *model.AppConfig) error {
	db := database.DB.Collection(table.AppConfig)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteAppConfig(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.AppConfig)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateAppConfig(ctx context.Context, r *model.AppConfig) error {
	db := database.DB.Collection(table.AppConfig)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetAppConfig(ctx context.Context, guid string) *model.AppConfig {
	db := database.DB.Collection(table.AppConfig)
	r := new(model.AppConfig)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(r)
	if err != nil {
		return nil
	}
	return r
}

func (database *Database) LoadAppConfig(ctx context.Context) ([]*model.AppConfig, error) {
	db := database.DB.Collection(table.AppConfig)
	cur, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	list := []*model.AppConfig{}
	for cur.Next(ctx) {
		r := new(model.AppConfig)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
