package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddRoleToAuthority(ctx context.Context, RoleToAuthority *model.RoleToAuthority) error {
	db := database.DB.Collection(table.RoleToAuthority)
	_, error := db.InsertOne(ctx, RoleToAuthority)
	return error
}

func (database *Database) DeleteRoleToAuthority(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.RoleToAuthority)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateRoleToAuthority(ctx context.Context, RoleToAuthority *model.RoleToAuthority) error {
	db := database.DB.Collection(table.RoleToAuthority)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", RoleToAuthority.ID}}, bson.D{{"$set", RoleToAuthority}})
	return error
}

func (database *Database) GetRoleToAuthority(ctx context.Context, guid primitive.ObjectID) (*model.RoleToAuthority, error) {
	db := database.DB.Collection(table.RoleToAuthority)
	RoleToAuthority := new(model.RoleToAuthority)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(RoleToAuthority)
	if res != nil {
		return nil, res
	}
	return RoleToAuthority, nil
}

func (database *Database) GetRoleToAuthorityKey(ctx context.Context, role, authority primitive.ObjectID) *model.RoleToAuthority {
	db := database.DB.Collection(table.RoleToAuthority)
	RoleToAuthority := new(model.RoleToAuthority)
	res := db.FindOne(ctx, bson.D{{"role", role}, {"authority", authority}}).Decode(RoleToAuthority)
	if res != nil {
		return nil
	}
	return RoleToAuthority
}

func (database *Database) DeleteRoleToAuthorityRole(ctx context.Context, role primitive.ObjectID) error {
	db := database.DB.Collection(table.RoleToAuthority)
	_, error := db.DeleteMany(ctx, bson.D{{"role", role}})
	return error
}

func (database *Database) DeleteRoleToAuthorityAuthority(ctx context.Context, authority primitive.ObjectID) error {
	db := database.DB.Collection(table.RoleToAuthority)
	_, error := db.DeleteMany(ctx, bson.D{{"authority", authority}})
	return error
}

func (database *Database) LoadRoleToAuthorityRole(ctx context.Context, role primitive.ObjectID) ([]*model.RoleToAuthority, error) {
	opt := NewOptions()
	opt.EQ[OptRole] = role
	return database.LoadRoleToAuthority(ctx, opt)
}

func (database *Database) LoadRoleToAuthorityAuthority(ctx context.Context, authority primitive.ObjectID) ([]*model.RoleToAuthority, error) {
	opt := NewOptions()
	opt.EQ[OptAuthority] = authority
	return database.LoadRoleToAuthority(ctx, opt)
}

func (database *Database) LoadRoleToAuthority(ctx context.Context, opt *option) ([]*model.RoleToAuthority, error) {
	db := database.DB.Collection(table.RoleToAuthority)
	need := make(map[OptionKey]string)
	need[OptRole] = "role"
	need[OptAuthority] = "authority"
	query, option := opt.toFind(need)
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, err
	}
	var list []*model.RoleToAuthority
	for cur.Next(ctx) {
		r := new(model.RoleToAuthority)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
