package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddInstitutionApplication(ctx context.Context, r *model.InstitutionApplication) error {
	db := database.DB.Collection(table.InstitutionApplication)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteInstitutionApplication(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.InstitutionApplication)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateInstitutionApplication(ctx context.Context, r *model.InstitutionApplication) error {
	db := database.DB.Collection(table.InstitutionApplication)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetInstitutionApplication(ctx context.Context, guid primitive.ObjectID) (*model.InstitutionApplication, error) {
	db := database.DB.Collection(table.InstitutionApplication)
	user := new(model.InstitutionApplication)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) LoadInstitutionApplication(ctx context.Context, opt *option) ([]*model.InstitutionApplication, int64, error) {
	db := database.DB.Collection(table.InstitutionApplication)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution_id"
	need[OptStatus] = "status"
	need[OptType] = "type"
	need[OptUser] = "user"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.InstitutionApplication
	for cur.Next(ctx) {
		r := new(model.InstitutionApplication)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
