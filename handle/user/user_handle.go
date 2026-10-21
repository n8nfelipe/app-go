package user_handle

import (
	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/controller"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/repositorie"
	"github.com/n8nfelipe/app-go.git/internal/usecase"
)

func GetRoutes(server *gin.Engine) {
	user := model.User{}
	userRepository := repositorie.NewUserRepository(&user)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userController := controller.NewUserController(userUseCase)

	server.POST("/login", userController.Login)
	server.POST("/users", userController.Create)
}
