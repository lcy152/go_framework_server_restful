package db

import (
	"context"
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (database *Database) AddUserToUserGroup(ctx context.Context, user *model.UserToUserGroup) error {
	db := database.DB.Collection(table.UserToUserGroup)
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUserToUserGroup(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToUserGroup)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserToUserGroup(ctx context.Context, user *model.UserToUserGroup) error {
	db := database.DB.Collection(table.UserToUserGroup)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUserToUserGroup(ctx context.Context, guid primitive.ObjectID) *model.UserToUserGroup {
	db := database.DB.Collection(table.UserToUserGroup)
	user := new(model.UserToUserGroup)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserToUserGroupUser(ctx context.Context, userGuid primitive.ObjectID) ([]*model.UserToUserGroup, error) {
	opt := NewOptions()
	opt.EQ[OptUser] = userGuid
	return database.LoadUserToUserGroup(ctx, opt)
}

func (database *Database) LoadUserToUserGroupGroup(ctx context.Context, groupGuid primitive.ObjectID) ([]*model.UserToUserGroup, error) {
	opt := NewOptions()
	opt.EQ[OptGroup] = groupGuid
	return database.LoadUserToUserGroup(ctx, opt)
}

func (database *Database) LoadUserToUserGroupInstitution(ctx context.Context, institutionId primitive.ObjectID) ([]*model.UserToUserGroup, error) {
	opt := NewOptions()
	opt.EQ[OptInstitution] = institutionId
	return database.LoadUserToUserGroup(ctx, opt)
}

func (database *Database) DeleteUserToUserGroupGroup(ctx context.Context, groupGuid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToUserGroup)
	_, error := db.DeleteMany(ctx, bson.D{{"group", groupGuid}})
	return error
}

func (database *Database) LoadUserToUserGroup(ctx context.Context, opt *option) ([]*model.UserToUserGroup, error) {
	db := database.DB.Collection(table.UserToUserGroup)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution"
	need[OptUser] = "user"
	need[OptGroup] = "group"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.UserToUserGroup
	for cur.Next(ctx) {
		r := new(model.UserToUserGroup)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
