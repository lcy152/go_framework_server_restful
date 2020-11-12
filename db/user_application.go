package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddUserApplication(ctx context.Context, r *model.UserApplication) error {
	db := database.DB.Collection(table.UserApplication)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteUserApplication(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserApplication)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserApplication(ctx context.Context, r *model.UserApplication) error {
	db := database.DB.Collection(table.UserApplication)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetUserApplication(ctx context.Context, guid primitive.ObjectID) *model.UserApplication {
	db := database.DB.Collection(table.UserApplication)
	user := new(model.UserApplication)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserApplication(ctx context.Context, opt *option) ([]*model.UserApplication, error) {
	db := database.DB.Collection(table.UserApplication)
	need := make(map[OptionKey]string)
	need[OptUser] = "user"
	need[OptFriend] = "friend"
	need[OptStatus] = "status"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.UserApplication
	for cur.Next(ctx) {
		r := new(model.UserApplication)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
