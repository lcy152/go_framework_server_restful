package db

import (
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func (database *Database) AddTask(ctx context.Context, r *model.Task) error {
	db := database.DB.Collection(table.Task)
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteTask(ctx context.Context, guid primitive.ObjectID) error {
	db := database.DB.Collection(table.Task)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateTask(ctx context.Context, r *model.Task) error {
	db := database.DB.Collection(table.Task)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.ID}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetTask(ctx context.Context, guid primitive.ObjectID) (*model.Task, error) {
	db := database.DB.Collection(table.Task)
	user := new(model.Task)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (database *Database) GetTaskByOrigin(ctx context.Context, originGuid string) *model.Task {
	db := database.DB.Collection(table.Task)
	user := new(model.Task)
	err := db.FindOne(ctx, bson.D{{"origin_guid", originGuid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func (database *Database) LoadTask(ctx context.Context, opt *option) ([]*model.Task, int64, error) {
	db := database.DB.Collection(table.Task)
	need := make(map[OptionKey]string)
	need[OptInstitution] = "institution_id"
	need[OptPatientName] = "ref_patient_name"
	need[OptPatientPid] = "ref_patient_pid"
	need[OptPatient] = "ref_patient_guid"
	need[OptPlan] = "ref_plan_guid"
	need[OptUser] = "execute_user_id"
	need[OptTaskState] = "task_state"
	need[OptTaskType] = "task_type"
	need[OptFxNumber] = "ref_fx_number"
	need[OptStudy] = "ref_study_guid"
	query, option := opt.toFind(need)
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.Task
	for cur.Next(ctx) {
		r := new(model.Task)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}
