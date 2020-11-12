package db

import (
	"context"
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (database *Database) AddUserToInstitution(ctx context.Context, user *model.UserToInstitution) error {
	db := database.DB.Collection(table.UserToInstitution)
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUserToInstitution(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToInstitution)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUserToInstitution(ctx context.Context, user *model.UserToInstitution) error {
	db := database.DB.Collection(table.UserToInstitution)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUserToInstitution(ctx context.Context, guid primitive.ObjectID) (*model.UserToInstitution, error) {
	db := database.DB.Collection(table.UserToInstitution)
	user := new(model.UserToInstitution)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) GetUserToInstitutionUserType(ctx context.Context, institutionId, userGuid primitive.ObjectID, userType string) (*model.UserToInstitution, error) {
	db := database.DB.Collection(table.UserToInstitution)
	uti := new(model.UserToInstitution)
	err := db.FindOne(ctx, bson.D{{"institution", institutionId}, {"user", userGuid}, {"type", userType}}).Decode(uti)
	if err != nil {
		return nil, err
	}
	return uti, nil
}

func (database *Database) GetUserToInstitutionDipperUser(ctx context.Context, institutionId primitive.ObjectID, userGuid string) (*model.UserToInstitution, error) {
	db := database.DB.Collection(table.UserToInstitution)
	uti := new(model.UserToInstitution)
	err := db.FindOne(ctx, bson.D{{"institution", institutionId}, {"dipper_user", userGuid}}).Decode(uti)
	if err != nil {
		return nil, err
	}
	return uti, nil
}

func (database *Database) LoadUserToInstitutionInstitutionUser(ctx context.Context, institutionId, userGuid primitive.ObjectID) ([]*model.UserToInstitution, error) {
	opt := NewOptions()
	opt.EQ[OptInstitution] = institutionId
	opt.EQ[OptUser] = userGuid
	return database.LoadUserToInstitution(ctx, opt)
}

func (database *Database) LoadUserToInstitutionInstitution(ctx context.Context, institutionId primitive.ObjectID) ([]*model.UserToInstitution, error) {
	opt := NewOptions()
	opt.EQ[OptInstitution] = institutionId
	return database.LoadUserToInstitution(ctx, opt)
}

func (database *Database) DeleteUserToInstitutionInstitution(ctx context.Context, institutionId primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToInstitution)
	_, error := db.DeleteMany(ctx, bson.D{{"institution", institutionId}})
	return error
}

func (database *Database) DeleteUserToInstitutionInstitutionUser(ctx context.Context, institutionId, userGuid primitive.ObjectID) error {
	db := database.DB.Collection(table.UserToInstitution)
	_, error := db.DeleteOne(ctx, bson.D{{"institution", institutionId}, {"user", userGuid}})
	return error
}

func (database *Database) GetUserToInstitutionCurrent(ctx context.Context, userGuid primitive.ObjectID) (*model.UserToInstitution, error) {
	db := database.DB.Collection(table.UserToInstitution)
	UserToInstitution := new(model.UserToInstitution)
	err := db.FindOne(ctx, bson.D{{"user", userGuid}, {"current", true}}).Decode(UserToInstitution)
	if err != nil {
		return nil, err
	}
	return UserToInstitution, nil
}

func (database *Database) LoadUserToInstitution(ctx context.Context, opt *option) ([]*model.UserToInstitution, error) {
	db := database.DB.Collection(table.UserToInstitution)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution"
	need[OptCurrent] = "current"
	need[OptManager] = "manager"
	need[OptType] = "type"
	need[OptFlag] = "flag"
	need[OptName] = "name"
	need[OptDipperUser] = "dipper_user"
	need[OptUser] = "user"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.UserToInstitution
	for cur.Next(ctx) {
		r := new(model.UserToInstitution)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
