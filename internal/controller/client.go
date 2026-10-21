package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
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

func (c *ClientController) GetClientByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid client ID"})
		return
	}

	client := c.ClientUseCase.GetClientByID(id)
	if client == nil {
		ctx.JSON(404, gin.H{"error": "Client not found"})
		return
	}

	ctx.JSON(200, client)
}

func (c *ClientController) GetClients(ctx *gin.Context) {
	clients := c.ClientUseCase.GetClients()
	if clients == nil {
		ctx.JSON(404, gin.H{"error": "Client not found"})
		return
	}

	ctx.JSON(200, clients)
}

func (c *ClientController) UpdateClient(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid client ID"})
		return
	}

	client := c.ClientUseCase.GetClientByID(id)
	if client == nil {
		ctx.JSON(404, gin.H{"error": "Client not found"})
		return
	}

	var clientForm ClientForm
	if err := ctx.ShouldBindJSON(&clientForm); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	client.Name = clientForm.Name
	client.Email = clientForm.Email

	err = c.ClientUseCase.UpdateClient(client)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, client)
}

func (c *ClientController) CreateClient(ctx *gin.Context) {
	var clientForm ClientForm
	if err := ctx.ShouldBindJSON(&clientForm); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	client := &model.Client{
		Name:  clientForm.Name,
		Email: clientForm.Email,
	}

	id, err := c.ClientUseCase.CreateClient(client)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create client"})
		return
	}
	client.ID = id
	ctx.JSON(201, client)
}
