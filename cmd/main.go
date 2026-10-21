package main

import (
	"log"

	"github.com/gin-gonic/gin"
	client_handle "github.com/n8nfelipe/app-go.git/handle/client"
	user_handle "github.com/n8nfelipe/app-go.git/handle/user"
)

func main() {

	server := gin.Default()
	client_handle.GetRoutes(server)
	user_handle.GetRoutes(server)

	err := server.Run(":8081")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
