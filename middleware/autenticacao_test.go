package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/n8nfelipe/app-go.git/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Setenv("JWT_SECRET", "segredo-de-teste")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protegido", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userID":   c.GetUint("userID"),
			"username": c.GetString("username"),
		})
	})
	return router
}

func doGet(router *gin.Engine, path, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func doRequest(router *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	return doGet(router, "/protegido", authHeader)
}

func TestGenerateToken_RetornaTokenValido(t *testing.T) {
	t.Setenv("JWT_SECRET", "segredo-de-teste")

	tokenString, err := middleware.GenerateToken(1, "felipe")

	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &middleware.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims := token.Claims.(*middleware.CustomClaims)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "felipe", claims.Username)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestGenerateToken_SemSecretFalha(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := middleware.GenerateToken(1, "felipe")

	require.Error(t, err)
}

func TestAuthMiddleware_SemToken(t *testing.T) {
	w := doRequest(setupRouter(t), "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_FormatoInvalido(t *testing.T) {
	w := doRequest(setupRouter(t), "Token abc123")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenInvalido(t *testing.T) {
	w := doRequest(setupRouter(t), "Bearer token-de-mentira")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenExpirado(t *testing.T) {
	router := setupRouter(t)

	claims := middleware.CustomClaims{
		UserID:   1,
		Username: "felipe",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)

	w := doRequest(router, "Bearer "+tokenString)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenAssinaturaErrada(t *testing.T) {
	router := setupRouter(t)

	claims := middleware.CustomClaims{
		UserID:   1,
		Username: "felipe",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("outro-segredo"))
	require.NoError(t, err)

	w := doRequest(router, "Bearer "+tokenString)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenValido(t *testing.T) {
	router := setupRouter(t)

	tokenString, err := middleware.GenerateToken(7, "ana")
	require.NoError(t, err)

	w := doRequest(router, "Bearer "+tokenString)

	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		UserID   uint   `json:"userID"`
		Username string `json:"username"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, uint(7), body.UserID)
	assert.Equal(t, "ana", body.Username)
}

func TestAuthMiddleware_SemSecretRetorna500(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protegido", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := doRequest(router, "Bearer qualquer-coisa")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRateLimit_BloqueiaAposBurst(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/limitado", middleware.RateLimit(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	var bloqueadas int
	for i := 0; i < 10; i++ {
		if doGet(router, "/limitado", "").Code == http.StatusTooManyRequests {
			bloqueadas++
		}
	}

	assert.GreaterOrEqual(t, bloqueadas, 1, "esperava ao menos uma resposta 429 após o burst")
}
