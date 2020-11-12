package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddGroupChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.GroupChat)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteGroupChat(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.GroupChat)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateGroupChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.GroupChat)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetGroupChat(ctx context.Context, guid primitive.ObjectID) (*model.ChatMessage, error) {
	db := database.DB.Collection(table.GroupChat)
	user := new(model.ChatMessage)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) LoadGroupChat(ctx context.Context, opt *option) ([]*model.ChatMessage, error) {
	db := database.DB.Collection(table.GroupChat)
	need := make(map[OptionKey]string)
	need[OptID] = "_id"
	need[OptSender] = "sender"
	need[OptGroup] = "group"
	need[OptType] = "type"
	need[OptData] = "data"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.ChatMessage
	for cur.Next(ctx) {
		r := new(model.ChatMessage)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
