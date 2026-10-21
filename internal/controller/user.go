package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/usecase"
	"github.com/n8nfelipe/app-go.git/middleware"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
}

func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		UserUseCase: userUseCase,
	}
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (uc *UserController) Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(400, gin.H{"error": "Dados de login inválidos"})
		return
	}

	username := loginData.Username
	password := loginData.Password
	user, err := uc.UserUseCase.GetByUsernameAndPassword(username, password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Credenciais inválidas"})
		return
	}
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Erro ao gerar token"})
		return
	}
	ctx.JSON(200, gin.H{"token": token, "userName": user.Username})
}

func (uc *UserController) Create(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "Dados do usuário inválidos"})
		return
	}
	id, err := uc.UserUseCase.Create(user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(201, gin.H{"userId": id})
}
