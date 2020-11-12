package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddDipperMessage(ctx context.Context, r *model.DipperMessage) error {
	db := database.DB.Collection(table.DipperMssage)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteDipperMessage(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.DipperMssage)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateDipperMessage(ctx context.Context, r *model.DipperMessage) error {
	db := database.DB.Collection(table.DipperMssage)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetDipperMessage(ctx context.Context, guid primitive.ObjectID) (*model.DipperMessage, error) {
	db := database.DB.Collection(table.DipperMssage)
	user := new(model.DipperMessage)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) LoadDipperMessage(ctx context.Context, opt *option) ([]*model.DipperMessage, error) {
	db := database.DB.Collection(table.DipperMssage)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "sender_institution"
	need[OptSender] = "sender"
	need[OptReceiver] = "receiver"
	need[OptData] = "data"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.DipperMessage
	for cur.Next(ctx) {
		r := new(model.DipperMessage)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
