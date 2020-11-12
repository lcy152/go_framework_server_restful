package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"tumor_server/db"
	"tumor_server/model"
	"tumor_server/service"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createDefault(sc *service.Container) {
	createAppConfig(sc)
	createUser(sc)
	createAuthority(sc)
}

func createAppConfig(sc *service.Container) {
	file, err := os.Open("config/app_config.json")
	if err != nil {
		logrus.Println(err)
		return
	}
	defer file.Close()
	appConfigList := []model.AppConfig{}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Println(err)
		return
	}
	err = json.Unmarshal(b, &appConfigList)
	if err != nil {
		logrus.Println(err)
		return
	}
	for _, appConfig := range appConfigList {
		opt := db.NewOptions()
		opt.EQ[db.OptName] = appConfig.Name
		list, _ := sc.DB.LoadAppConfig(context.TODO(), opt)
		if len(list) == 0 {
			appConfig.ID = primitive.NewObjectID()
			err = sc.DB.AddAppConfig(context.TODO(), &appConfig)
			if err != nil {
				logrus.Println(err)
				return
			}
		}
	}
}

func createUser(sc *service.Container) {
	file, err := os.Open("config/user.json")
	if err != nil {
		logrus.Println(err)
		return
	}
	defer file.Close()
	dataList := []model.User{}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Println(err)
		return
	}
	err = json.Unmarshal(b, &dataList)
	if err != nil {
		logrus.Println(err)
		return
	}
	for _, data := range dataList {
		_, err := sc.DB.GetUserByPhone(context.TODO(), data.Phone)
		if err != nil {
			data.ID = primitive.NewObjectID()
			err = sc.DB.AddUser(context.TODO(), &data)
			if err != nil {
				logrus.Println(err)
				return
			}
		}
	}
}

func createAuthority(sc *service.Container) {
	file, err := os.Open("config/authority.json")
	if err != nil {
		logrus.Println(err)
		return
	}
	defer file.Close()
	dataList := []model.Authority{}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Println(err)
		return
	}
	err = json.Unmarshal(b, &dataList)
	if err != nil {
		logrus.Println(err)
		return
	}
	for _, data := range dataList {
		_, err := sc.DB.GetAuthorityNumber(context.TODO(), data.Number)
		if err != nil {
			data.ID = primitive.NewObjectID()
			err = sc.DB.AddAuthority(context.TODO(), &data)
			if err != nil {
				logrus.Println(err)
				return
			}
		}
	}
}
