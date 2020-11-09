package db

import (
	"context"
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (database *Database) AddUserRouter(ctx context.Context, user *model.UserRouter) error {
	db := database.DB.Collection(table.UserRouter)
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUserRouter(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.UserRouter)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserRouter(ctx context.Context, user *model.UserRouter) error {
	db := database.DB.Collection(table.UserRouter)
	tn := time.Now()
	user.LastModTime = tn
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.Guid}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUserRouter(ctx context.Context, guid string) *model.UserRouter {
	db := database.DB.Collection(table.UserRouter)
	user := new(model.UserRouter)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) GetUserRouterByTumorUser(ctx context.Context, institutionId, guid string) *model.UserRouter {
	db := database.DB.Collection(table.UserRouter)
	UserRouter := new(model.UserRouter)
	err := db.FindOne(ctx, bson.D{{"institution_id", institutionId}, {"user_guid", guid}}).Decode(UserRouter)
	if err != nil {
		return nil
	}
	return UserRouter
}

func (database *Database) GetUserRouterByDipperUser(ctx context.Context, institutionId, guid string) *model.UserRouter {
	db := database.DB.Collection(table.UserRouter)
	UserRouter := new(model.UserRouter)
	err := db.FindOne(ctx, bson.D{{"institution_id", institutionId}, {"dipper_user", guid}}).Decode(UserRouter)
	if err != nil {
		return nil
	}
	return UserRouter
}

func (database *Database) LoadUserRouter(ctx context.Context, opt *option) ([]*model.UserRouter) {
	db := database.DB.Collection(table.UserRouter)
	need := make(map[OptionKey]string)
	need[OptInstitutionId] = "institution_id"
	need[OptCurrent] = "is_current"
	need[OptRole] = "role_id_list"
	need[OptType] = "type"
	need[OptFlag] = "flag"
	need[OptName] = "name"
	need[OptDipperUser] = "dipper_user"
	need[OptUserGuid] = "user_guid"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		log.Println(err)
		return nil
	}
	var list []*model.UserRouter
	for cur.Next(ctx) {
		r := new(model.UserRouter)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list
}
