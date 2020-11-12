package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddSingleChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.SingleChat)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteSingleChat(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.SingleChat)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateSingleChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.SingleChat)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetSingleChat(ctx context.Context, guid primitive.ObjectID) *model.ChatMessage {
	db := database.DB.Collection(table.SingleChat)
	user := new(model.ChatMessage)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadSingleChat(ctx context.Context, opt *option) ([]*model.ChatMessage, error) {
	db := database.DB.Collection(table.SingleChat)
	need := make(map[OptionKey]string)
	need[OptID] = "_id"
	need[OptInstitution] = "institution"
	need[OptSender] = "sender"
	need[OptReceiver] = "receiver"
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
