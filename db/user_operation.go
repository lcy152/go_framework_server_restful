package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddUserOperation(ctx context.Context, r *model.UserOperation) error {
	db := database.DB.Collection(table.UserOperation)
	tn := time.Now()
	r.CreateTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteUserOperation(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.UserOperation)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserOperation(ctx context.Context, r *model.UserOperation) error {
	db := database.DB.Collection(table.UserOperation)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetUserOperation(ctx context.Context, guid string) *model.UserOperation {
	db := database.DB.Collection(table.UserOperation)
	user := new(model.UserOperation)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserOperation(ctx context.Context, opt *option) ([]*model.UserOperation, int64, error) {
	db := database.DB.Collection(table.UserOperation)
	need := make(map[OptionKey]string)
	need[OptUserGuid] = "user_guid"
	need[OptInstitutionId] = "institution_id"
	need[OptType] = "type"
	need[OptCreateTime] = "create_time"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.UserOperation
	for cur.Next(ctx) {
		r := new(model.UserOperation)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
