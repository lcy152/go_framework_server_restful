package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddInstitution(ctx context.Context, r *model.Institution) error {
	db := database.DB.Collection(table.Institution)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteInstitution(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.Institution)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateInstitution(ctx context.Context, r *model.Institution) error {
	db := database.DB.Collection(table.Institution)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetInstitution(ctx context.Context, guid primitive.ObjectID) (*model.Institution, error) {
	db := database.DB.Collection(table.Institution)
	user := new(model.Institution)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) GetInstitutionCode(ctx context.Context, code string) (*model.Institution, error) {
	db := database.DB.Collection(table.Institution)
	user := new(model.Institution)
	err := db.FindOne(ctx, bson.D{{"code", code}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) LoadInstitution(ctx context.Context, opt *option) ([]*model.Institution, int64, error) {
	db := database.DB.Collection(table.Institution)
	need := make(map[OptionKey]string)
	need[OptName] = "name"
	need[OptID] = "_id"
	need[OptAddress] = "address"
	need[OptCode] = "code"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.Institution
	for cur.Next(ctx) {
		r := new(model.Institution)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
