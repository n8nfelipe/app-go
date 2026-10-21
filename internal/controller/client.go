package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/usecase"
)

type ClientController struct {
	ClientUseCase *usecase.ClientUseCase
}

type ClientForm struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func NewClientController(clientUseCase *usecase.ClientUseCase) *ClientController {
	return &ClientController{
		ClientUseCase: clientUseCase,
	}
}

func parseClientID(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		respondError(ctx, apperror.New(apperror.CodeInvalidInput, "ID de cliente inválido", err), "validando ID do cliente")
		return 0, false
	}
	return id, true
}

func authenticatedUserID(ctx *gin.Context) (uint, bool) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		respondError(ctx, apperror.New(apperror.CodeUnauthorized, "Não autenticado", nil), "validando autenticação")
		return 0, false
	}
	return userID, true
}

func (c *ClientController) GetClientByID(ctx *gin.Context) {
	id, ok := parseClientID(ctx)
	if !ok {
		return
	}
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	client, err := c.ClientUseCase.GetClientByID(id, userID)
	if err != nil {
		respondError(ctx, err, "buscar client")
		return
	}

	ctx.JSON(http.StatusOK, client)
}

func (c *ClientController) GetClients(ctx *gin.Context) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	clients, err := c.ClientUseCase.GetClients(userID)
	if err != nil {
		respondError(ctx, err, "listar clients")
		return
	}
	if clients == nil {
		clients = []model.Client{}
	}

	ctx.JSON(http.StatusOK, clients)
}

func (c *ClientController) UpdateClient(ctx *gin.Context) {
	id, ok := parseClientID(ctx)
	if !ok {
		return
	}
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	client, err := c.ClientUseCase.GetClientByID(id, userID)
	if err != nil {
		respondError(ctx, err, "buscar client para atualizar")
		return
	}

	var clientForm ClientForm
	if err := ctx.ShouldBindJSON(&clientForm); err != nil {
		respondError(ctx, apperror.New(apperror.CodeInvalidInput, "Corpo da requisição inválido", err), "validando corpo da requisição")
		return
	}

	client.Name = clientForm.Name
	client.Email = clientForm.Email

	if err := c.ClientUseCase.UpdateClient(client); err != nil {
		respondError(ctx, err, "atualizar client")
		return
	}

	ctx.JSON(http.StatusOK, client)
}

func (c *ClientController) CreateClient(ctx *gin.Context) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	var clientForm ClientForm
	if err := ctx.ShouldBindJSON(&clientForm); err != nil {
		respondError(ctx, apperror.New(apperror.CodeInvalidInput, "Corpo da requisição inválido", err), "validando corpo da requisição")
		return
	}

	client := &model.Client{
		Name:   clientForm.Name,
		Email:  clientForm.Email,
		UserID: userID,
	}

	id, err := c.ClientUseCase.CreateClient(client)
	if err != nil {
		respondError(ctx, err, "criar client")
		return
	}
	client.ID = id
	ctx.JSON(http.StatusCreated, client)
}

func (c *ClientController) DeleteClient(ctx *gin.Context) {
	id, ok := parseClientID(ctx)
	if !ok {
		return
	}
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	if err := c.ClientUseCase.DeleteClient(id, userID); err != nil {
		respondError(ctx, err, "excluir client")
		return
	}

	ctx.Status(http.StatusNoContent)
}
