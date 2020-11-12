package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddAuthority(ctx context.Context, Authority *model.Authority) error {
	db := database.DB.Collection(table.Authority)
	_, error := db.InsertOne(ctx, Authority)
	return error
}

func (database *Database) DeleteAuthority(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.Authority)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateAuthority(ctx context.Context, Authority *model.Authority) error {
	db := database.DB.Collection(table.Authority)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", Authority.ID}}, bson.D{{"$set", Authority}})
	return error
}

func (database *Database) GetAuthority(ctx context.Context, guid primitive.ObjectID) (*model.Authority, error) {
	db := database.DB.Collection(table.Authority)
	Authority := new(model.Authority)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(Authority)
	if err != nil {
		return nil, err
	}
	return Authority, nil
}

func (database *Database) GetAuthorityNumber(ctx context.Context, number string) (*model.Authority, error) {
	db := database.DB.Collection(table.Authority)
	Authority := new(model.Authority)
	err := db.FindOne(ctx, bson.D{{"number", number}}).Decode(Authority)
	if err != nil {
		return nil, err
	}
	return Authority, nil
}

func (database *Database) LoadAuthority(ctx context.Context, opt *option) ([]*model.Authority, error) {
	db := database.DB.Collection(table.Authority)
	need := make(map[OptionKey]string)
	need[OptName] = "name"
	need[OptNumber] = "number"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.Authority
	for cur.Next(ctx) {
		r := new(model.Authority)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
