package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/stretchr/testify/assert"
)

func TestRespondError_TraduzErroTipadoParaHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	respondError(ctx, apperror.New(apperror.CodeNotFound, "Cliente não encontrado", nil), "buscar cliente")

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.JSONEq(t, `{"error":"Cliente não encontrado"}`, recorder.Body.String())
}

func TestRespondError_NaoExpõeErroInterno(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	respondError(ctx, apperror.Wrap(assert.AnError, "consultar banco"), "buscar cliente")

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"error":"Erro interno no servidor"}`, recorder.Body.String())
}
