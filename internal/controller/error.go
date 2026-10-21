package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/apperror"
)

func respondError(ctx *gin.Context, err error, operation string) {
	status := http.StatusInternalServerError
	message := "Erro interno no servidor"

	var applicationError apperror.Error
	if errors.As(err, &applicationError) {
		switch applicationError.Code() {
		case apperror.CodeNotFound:
			status = http.StatusNotFound
			message = applicationError.PublicMessage()
		case apperror.CodeUnauthorized:
			status = http.StatusUnauthorized
			message = applicationError.PublicMessage()
		case apperror.CodeInvalidInput:
			status = http.StatusBadRequest
			message = applicationError.PublicMessage()
		case apperror.CodeInternal:
			log.Printf("erro ao %s: %v", operation, err)
		default:
			log.Printf("erro ao %s: %v", operation, err)
		}
	} else {
		log.Printf("erro ao %s: %v", operation, err)
	}

	ctx.JSON(status, gin.H{"error": message})
}
