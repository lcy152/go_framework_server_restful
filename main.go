package main

import (
	"log"
	"strconv"
	service "tumor_server/service"
)

func main() {
	config := ParseConfig("serverConfig.ini")
	InitLog(config.LogPath, config.LogLevel, "tumor_server.log", 1, 7)
	container := service.NewContainerInstance(config)
	createDefault(container)
	s := NewServer()
	log.Println("start server success")
	s.Run(":" + strconv.Itoa(config.Port))
}
