package db

import (
	"log"
	"time"

	"github.com/fatih/structs"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type Database struct {
	connectUrl string
	Client     *mongo.Client
	DB         *mongo.Database
}

func NewDatabase(url string) *Database {
	database := new(Database)
	database.connectUrl = url

	//uri := "mongodb://datu_super_root:c74c112dc3130e35e9ac88c90d214555__strong@localhost:27127/datu_data"
	opts := options.Client()
	opts.SetDirect(true)
	opts.ApplyURI(url)
	client, err := mongo.NewClient(opts)
	database.Client = client
	database.DB = client.Database("tumor_data")
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	database.ensureIndex(context.TODO())
	database.CreateEmptyTable()
	//defer cancel()
	return database
}

func (db *Database) ensureIndex(ctx context.Context) error {

	ensure := func(table string, ks bson.D) error {
		idxs := db.DB.Collection(table).Indexes()
		cur, err := idxs.List(ctx)
		if err != nil {
			return err
		}

		found := 0
		for cur.Next(ctx) {
			d := bson.M{}
			if err = cur.Decode(d); err != nil {
				return err
			}

			v := d["key"].(bson.M)
			bsonM := ks.Map()
			for k := range bsonM {
				if _, ok := v[k]; ok {
					found = found + 1
				}
			}
		}
		if found == len(ks) {
			return nil
		}

		idm := mongo.IndexModel{
			Keys:    ks,
			Options: options.Index().SetUnique(true),
		}
		_, err = idxs.CreateOne(ctx, idm)
		return err
	}
	err := ensure(table.User, bson.D{{"phone", 1}})
	if err != nil {
		return err
	}
	err = ensure(table.User, bson.D{{"token", 1}})
	return err
}

func (db *Database) CreateEmptyTable() {
	collectionList, err := db.DB.ListCollectionNames(context.TODO(), bson.D{{}})
	CreateTableFunc := func(name string) {
		c := db.DB.Collection(name)
		i := uuid.Must(uuid.NewV4(), nil).String()
		_, err := c.InsertOne(context.TODO(), bson.D{{"_id", i}})
		if err == nil {
			_, _ = c.DeleteOne(context.TODO(), bson.D{{"_id", i}})
		}
	}
	if err == nil {
		colMap := make(map[string]bool)
		for _, v := range collectionList {
			colMap[v] = true
		}
		tableNames := structs.Values(table)
		for _, v := range tableNames {
			name := v.(string)
			if !colMap[name] {
				CreateTableFunc(name)
			}
		}
	}

}
