package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddUserGroup(ctx context.Context, r *model.UserGroup) error {
	db := database.DB.Collection(table.UserGroup)
	tn := time.Now()
	r.CreateTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteUserGroup(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.UserGroup)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserGroup(ctx context.Context, r *model.UserGroup) error {
	db := database.DB.Collection(table.UserGroup)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetUserGroup(ctx context.Context, guid string) *model.UserGroup {
	db := database.DB.Collection(table.UserGroup)
	user := new(model.UserGroup)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserGroup(ctx context.Context, opt *option) ([]*model.UserGroup, int64, error) {
	db := database.DB.Collection(table.UserGroup)
	need := make(map[OptionKey]string)
	need[OptManager] = "manager._id"
	need[OptMember] = "member._id"
	need[OptCreateTime] = "create_time"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.UserGroup
	for cur.Next(ctx) {
		r := new(model.UserGroup)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
