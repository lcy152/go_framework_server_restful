package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddUserApplication(ctx context.Context, r *model.UserApplication) error {
	db := database.DB.Collection(table.UserApplication)
	tn := time.Now()
	r.CreateTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteUserApplication(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.UserApplication)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserApplication(ctx context.Context, r *model.UserApplication) error {
	db := database.DB.Collection(table.UserApplication)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetUserApplication(ctx context.Context, guid string) *model.UserApplication {
	db := database.DB.Collection(table.UserApplication)
	user := new(model.UserApplication)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserApplication(ctx context.Context, opt *option) ([]*model.UserApplication, int64, error) {
	db := database.DB.Collection(table.UserApplication)
	need := make(map[OptionKey]string)
	need[OptUserGuid] = "user_guid"
	need[OptInstitutionId] = "institution_id"
	need[OptStatus] = "status"
	need[OptType] = "type"
	need[OptCreateTime] = "create_time"
	need[OptCreator] = "creator"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.UserApplication
	for cur.Next(ctx) {
		r := new(model.UserApplication)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
