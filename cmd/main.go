package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/config"
	client_handle "github.com/n8nfelipe/app-go.git/handle/client"
	user_handle "github.com/n8nfelipe/app-go.git/handle/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Falha ao carregar configuração: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET não definido (configure no .env ou no ambiente)")
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	client_handle.GetRoutes(server)
	user_handle.GetRoutes(server)

	addr := ":" + cfg.Server.Port
	log.Printf("Servidor escutando em %s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatalf("Falha ao iniciar servidor: %v", err)
	}
}
