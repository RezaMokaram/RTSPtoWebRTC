package main

import (
	"fmt"
	"webRTC/server"
	"webRTC/services"
)

func main() {
	config := services.NewConfigService().LoadConfig()
	streamService := services.NewStreamService()

	streamService.ServeStreams(config)

	e := server.NewServer()
	fmt.Println("Server build successfully ...")
	server.RunServer(e, config)

}
