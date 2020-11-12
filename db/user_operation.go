package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddUserOperation(ctx context.Context, r *model.UserOperation) error {
	db := database.DB.Collection(table.UserOperation)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteUserOperation(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserOperation)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserOperation(ctx context.Context, r *model.UserOperation) error {
	db := database.DB.Collection(table.UserOperation)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetUserOperation(ctx context.Context, guid primitive.ObjectID) *model.UserOperation {
	db := database.DB.Collection(table.UserOperation)
	user := new(model.UserOperation)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserOperation(ctx context.Context, opt *option) ([]*model.UserOperation, int64, error) {
	db := database.DB.Collection(table.UserOperation)
	need := make(map[OptionKey]string)
	need[OptID] = "_id"
	need[OptUser] = "user_guid"
	need[OptInstitution] = "institution"
	need[OptType] = "type"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.UserOperation
	for cur.Next(ctx) {
		r := new(model.UserOperation)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
