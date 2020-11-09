package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddDipperMessage(ctx context.Context, r *model.DipperMessage) error {
	db := database.DB.Collection(table.DipperMssage)
	tn := time.Now()
	r.CreateTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteDipperMessage(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.DipperMssage)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateDipperMessage(ctx context.Context, r *model.DipperMessage) error {
	db := database.DB.Collection(table.DipperMssage)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetDipperMessage(ctx context.Context, guid string) *model.DipperMessage {
	db := database.DB.Collection(table.DipperMssage)
	user := new(model.DipperMessage)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadDipperMessage(ctx context.Context, opt *option) ([]*model.DipperMessage, error) {
	db := database.DB.Collection(table.DipperMssage)
	need := make(map[OptionKey]string)
	need[OptInstitutionId] = "institution_id"
	need[OptSenderGuid] = "sender_guid"
	need[OptReceiverGuid] = "receiver_guid"
	need[OptCreateTime] = "create_time"
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
