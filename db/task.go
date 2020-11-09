package db

import (
	"log"
	"time"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (database *Database) AddTask(ctx context.Context, r *model.Task) error {
	db := database.DB.Collection(table.Task)
	tn := time.Now()
	r.CreatedTime = tn
	_, error := db.InsertOne(ctx, r)
	return error
}

func (database *Database) DeleteTask(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.Task)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateTask(ctx context.Context, r *model.Task) error {
	db := database.DB.Collection(table.Task)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", r.Guid}}, bson.D{{"$set", r}})
	return error
}

func (database *Database) GetTask(ctx context.Context, guid string) *model.Task {
	db := database.DB.Collection(table.Task)
	user := new(model.Task)
	err := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if err != nil {
		return nil
	}
	return user
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
	need[OptInstitutionId] = "institution_id"
	need[OptPatientName] = "ref_patient_name"
	need[OptPatientPid] = "ref_patient_pid"
	need[OptPatientGuid] = "ref_patient_guid"
	need[OptPlanGuid] = "ref_plan_guid"
	need[OptUserID] = "execute_user_id"
	need[OptTaskState] = "task_state"
	need[OptTaskType] = "task_type"
	need[OptFxNumber] = "ref_fx_number"
	need[OptCreateTime] = "created_time"
	need[OptLastModTime] = "last_mod_time"
	need[OptStudyGuid] = "ref_study_guid"
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
