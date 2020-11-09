package db

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"
	"tumor_server/model"

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
	database.createDefaultData()
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

func (db *Database) createDefaultData() {
	appCon, err := db.LoadAppConfig(context.TODO())
	if err == nil && len(appCon) == 0 {
		appConfig := &model.AppConfig{
			Guid:       "000001",
			AppVersion: "1.0.0",
			AppUrl:     "fs/AppVersion/tumor.apk",
		}
		err := db.AddAppConfig(context.TODO(), appConfig)
		if err != nil {
			log.Println("error: db auto create appconfig error")
		}
	}
	admin := db.GetUser(context.TODO(), model.AdminGuid)
	if admin == nil {
		data := []byte("datu2012")
		md5Ctx := md5.New()
		md5Ctx.Write(data)
		cipherStr := md5Ctx.Sum(nil)
		admin = &model.User{
			Guid:     model.AdminGuid,
			Name:     "admin",
			Password: hex.EncodeToString(cipherStr),
			Phone:    "15221536381",
			Token:    "eyJndWlkIjoiYWRtaW4iLCJsb2dpbl90aW1lIjoxNTg5OTY0Nzc4fQ==",
		}
		err := db.AddUser(context.TODO(), admin)
		if err != nil {
			log.Println("error: db auto create user error")
		}
	}
	ins := db.GetInstitution(context.TODO(), "datu")
	if ins == nil {
		ins = &model.Institution{
			Guid:    "datu",
			Name:    "shanghai datu",
			Manager: []string{admin.Guid},
			KeyCode: "ehpPUIlBTO6d5UFYI9KHkRKlpX",
		}
		err := db.AddInstitution(context.TODO(), ins)
		if err != nil {
			log.Println("error: db auto create institution error")
		}
	}
	ur := db.GetUserRouter(context.TODO(), "admin_router")
	if ur == nil {
		ur = &model.UserRouter{
			Guid:            "admin_router",
			InstitutionId:   ins.Guid,
			InstitutionName: ins.Name,
			UserGuid:        admin.Guid,
			Flag:            "user",
			IsCurrent:       true,
		}
		err := db.AddUserRouter(context.TODO(), ur)
		if err != nil {
			log.Println("error: db auto create user_router error")
		}
	}
}

func NewUUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}
