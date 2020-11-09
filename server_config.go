package main

import (
	"fmt"
	"os"
	service "tumor_server/service"

	"gopkg.in/ini.v1"
)

func ParseConfig(file string) *service.ServerConfig {
	cfg, err := ini.Load(file)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	pwd, _ := os.Getwd()

	config := &service.ServerConfig{}
	config.Port = cfg.Section("config").Key("HttpPort").MustInt(2324)
	config.RedisIP = cfg.Section("config").Key("RedisIP").MustString("127.0.0.1")
	config.RedisPort = cfg.Section("config").Key("RedisPort").MustInt(6379)
	config.ExpireTime = cfg.Section("config").Key("ExpireTime").MustInt(3600)
	config.ShortMessageInvalidTime = cfg.Section("config").Key("ShortMessageInvalidTime").MustInt(300)
	config.ShortMessageSpaceTime = cfg.Section("config").Key("ShortMessageSpaceTime").MustInt(60)
	config.ShortMessageExpireTime = cfg.Section("config").Key("ShortMessageExpireTime").MustInt(36000)
	config.ShortMessageLimitedCount = cfg.Section("config").Key("ShortMessageLimitedCount").MustInt(5)
	config.BlackListTime = cfg.Section("config").Key("BlackListTime").MustInt(60)
	config.BlackListLimiteTime = cfg.Section("config").Key("BlackListLimiteTime").MustInt(5)
	config.BlackListLimitCount = cfg.Section("config").Key("BlackListLimitCount").MustInt(10)
	config.LogPath = cfg.Section("config").Key("LOGPATH").MustString(pwd)
	config.LogLevel = cfg.Section("config").Key("LogLevel").MustInt(7)
	config.DatabaseUrl = cfg.Section("config").Key("DatabaseIP").MustString("127.0.0.1")
	config.RabbitMQIP = cfg.Section("config").Key("RabbitMQIP").MustString("127.0.0.1:5672")
	config.RabbitMQVHost = cfg.Section("config").Key("RabbitMQVHost").MustString("tumor")
	config.RabbitMQVUserName = cfg.Section("config").Key("RabbitMQVUserName").MustString("")
	config.RabbitMQVPassword = cfg.Section("config").Key("RabbitMQVPassword").MustString("")
	return config
}
