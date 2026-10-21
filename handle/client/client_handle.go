package client_handle

import (
	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/controller"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/repositorie"
	"github.com/n8nfelipe/app-go.git/internal/usecase"
	"github.com/n8nfelipe/app-go.git/middleware"
)

func GetRoutes(server *gin.Engine) {
	client := model.Client{}
	clientRepository := repositorie.NewClientRepository(&client)
	clientUseCase := usecase.NewClientUseCase(clientRepository)
	clientController := controller.NewClientController(clientUseCase)

	protected := server.Group("/clients")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/:id", clientController.GetClientByID)
		protected.GET("/", clientController.GetClients)
		protected.PUT("/:id", clientController.UpdateClient)
		protected.POST("/", clientController.CreateClient)
	}

}
