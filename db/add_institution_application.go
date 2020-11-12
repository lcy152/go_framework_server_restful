package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddAddInstitutionApplication(ctx context.Context, r *model.AddInstitutionApplication) error {
	db := database.DB.Collection(table.AddInstitutionApplication)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteAddInstitutionApplication(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.AddInstitutionApplication)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateAddInstitutionApplication(ctx context.Context, r *model.AddInstitutionApplication) error {
	db := database.DB.Collection(table.AddInstitutionApplication)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetAddInstitutionApplication(ctx context.Context, guid primitive.ObjectID) (*model.AddInstitutionApplication, error) {
	db := database.DB.Collection(table.AddInstitutionApplication)
	user := new(model.AddInstitutionApplication)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) LoadAddInstitutionApplication(ctx context.Context, opt *option) ([]*model.AddInstitutionApplication, int64, error) {
	db := database.DB.Collection(table.AddInstitutionApplication)
	need := make(map[OptionKey]string)
	need[OptStatus] = "status"
	need[OptUser] = "user"
	need[OptUserName] = "user_name"
	need[OptInstitutionName] = "institution.name"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.AddInstitutionApplication
	for cur.Next(ctx) {
		r := new(model.AddInstitutionApplication)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
