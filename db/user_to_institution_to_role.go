package db

import (
	"context"
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (database *Database) AddUserToInstitutionToRole(ctx context.Context, user *model.UserToInstitutionToRole) error {
	db := database.DB.Collection(table.UserToInstitutionToRole)
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUserToInstitutionToRole(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToInstitutionToRole)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserToInstitutionToRole(ctx context.Context, user *model.UserToInstitutionToRole) error {
	db := database.DB.Collection(table.UserToInstitutionToRole)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUserToInstitutionToRole(ctx context.Context, guid primitive.ObjectID) *model.UserToInstitutionToRole {
	db := database.DB.Collection(table.UserToInstitutionToRole)
	user := new(model.UserToInstitutionToRole)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUserToInstitutionToRoleUserInstitution(ctx context.Context, institutionId, userGuid primitive.ObjectID) ([]*model.UserToInstitutionToRole, error) {
	opt := NewOptions()
	opt.EQ[OptInstitution] = institutionId
	opt.EQ[OptUser] = userGuid
	return database.LoadUserToInstitutionToRole(ctx, opt)
}

func (database *Database) LoadUserToInstitutionToRole(ctx context.Context, opt *option) ([]*model.UserToInstitutionToRole, error) {
	db := database.DB.Collection(table.UserToInstitutionToRole)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution"
	need[OptUser] = "user"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var list []*model.UserToInstitutionToRole
	for cur.Next(ctx) {
		r := new(model.UserToInstitutionToRole)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
