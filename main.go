package main

import (
	"log"
	"strconv"
	service "tumor_server/service"
)

func main() {
	config := ParseConfig("serverConfig.ini")
	service.NewContainerInstance(config)
	s := NewServer()
	log.Println("start server success")
	s.Run(":" + strconv.Itoa(config.Port))
}
