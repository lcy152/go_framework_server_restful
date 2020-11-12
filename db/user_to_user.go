package db

import (
	"context"
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (database *Database) AddUserToUser(ctx context.Context, user *model.UserToUser) error {
	db := database.DB.Collection(table.UserToUser)
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUserToUser(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToUser)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserToUser(ctx context.Context, user *model.UserToUser) error {
	db := database.DB.Collection(table.UserToUser)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUserToUser(ctx context.Context, guid primitive.ObjectID) *model.UserToUser {
	db := database.DB.Collection(table.UserToUser)
	user := new(model.UserToUser)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserToUserUser(ctx context.Context, userGUID primitive.ObjectID) ([]*model.UserToUser, error) {
	opt := NewOptions()
	opt.EQ[OptUser] = userGUID
	return database.LoadUserToUser(ctx, opt)
}

func (database *Database) LoadUserToUser(ctx context.Context, opt *option) ([]*model.UserToUser, error) {
	db := database.DB.Collection(table.UserToUser)
	need := make(map[OptionKey]string)
	need[OptUser] = "user"
	need[OptUserName] = "user_name"
	need[OptFriend] = "friend"
	need[OptFriendName] = "friend_name"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.UserToUser
	for cur.Next(ctx) {
		r := new(model.UserToUser)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
