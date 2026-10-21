package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/apperror"
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

func (uc *UserController) Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		respondError(ctx, apperror.New(apperror.CodeInvalidInput, "Dados de login inválidos", err), "validando login")
		return
	}

	user, err := uc.UserUseCase.GetByUsernameAndPassword(loginData.Username, loginData.Password)
	if err != nil {
		respondError(ctx, err, "autenticando usuário")
		return
	}
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		respondError(ctx, apperror.Wrap(err, "gerando token"), "gerar token")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "userName": user.Username})
}

func (uc *UserController) Create(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		respondError(ctx, apperror.New(apperror.CodeInvalidInput, "Dados do usuário inválidos", err), "validando usuário")
		return
	}
	id, err := uc.UserUseCase.Create(user)
	if err != nil {
		respondError(ctx, err, "criar usuário")
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"userId": id})
}
