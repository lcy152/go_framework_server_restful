package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func (database *Database) AddSingleChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.SingleChat)
	tn := time.Now()
	r.CreateTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteSingleChat(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.SingleChat)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateSingleChat(ctx context.Context, r *model.ChatMessage) error {
	db := database.DB.Collection(table.SingleChat)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetSingleChat(ctx context.Context, guid string) *model.ChatMessage {
	db := database.DB.Collection(table.SingleChat)
	user := new(model.ChatMessage)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadSingleChat(ctx context.Context, opt *option, user1, user2, search string) ([]*model.ChatMessage, error) {
	db := database.DB.Collection(table.SingleChat)
	option := &options.FindOptions{}
	userList := []string{user1, user2}
	ascend := 1
	if !opt.Ascend {
		ascend = -1
	}
	query := bson.M{"sender_guid": bson.M{"$in": userList}, "receiver_guid": bson.M{"$in": userList}, "data": bson.M{"$regex": search, "$options": "$i"}}
	option.Sort = bson.M{"create_time": ascend}
	skip := opt.PageIndex * opt.PageSize
	option.Skip = &skip
	limit := opt.PageSize
	option.Limit = &limit
	cur, err := db.Find(ctx, query, option)
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
