package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddRole(ctx context.Context, Role *model.Role) error {
	db := database.DB.Collection(table.Role)
	_, error := db.InsertOne(ctx, Role)
	return error
}

func (database *Database) DeleteRole(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.Role)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateRole(ctx context.Context, Role *model.Role) error {
	db := database.DB.Collection(table.Role)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", Role.ID}}, bson.D{{"$set", Role}})
	return error
}

func (database *Database) GetRole(ctx context.Context, guid primitive.ObjectID) (*model.Role, error) {
	db := database.DB.Collection(table.Role)
	Role := new(model.Role)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(Role)
	if err != nil {
		return nil, err
	}
	return Role, nil
}

func (database *Database) LoadRole(ctx context.Context, opt *option) ([]*model.Role, int64, error) {
	db := database.DB.Collection(table.Role)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution"
	need[OptName] = "name"
	need[OptStatus] = "status"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.Role
	for cur.Next(ctx) {
		r := new(model.Role)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
