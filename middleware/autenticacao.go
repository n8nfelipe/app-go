package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta usada para assinar os tokens (lida da env a cada uso)
func jwtSecretKey() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET não definido")
	}
	return []byte(secret), nil
}

// CustomClaims define os dados que serão guardados dentro do Token (Payload)
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ----------------------------------------------------
// 1. GERADOR DE TOKEN (Simulação após o Login bem-sucedido)
// ----------------------------------------------------
func GenerateToken(userID uint, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // Expira em 2h
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := jwtSecretKey()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

// ----------------------------------------------------
// 2. MIDDLEWARE DE AUTENTICAÇÃO (Protetor de Rotas)
// ----------------------------------------------------
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, err := jwtSecretKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuração de autenticação inválida no servidor"})
			c.Abort()
			return
		}

		// 1. Extrai o cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho Authorization não enviado"})
			c.Abort() // Interrompe a execução dos próximos handlers
			return
		}

		// 2. Verifica o formato "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido (Use 'Bearer <token>')"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Faz o Parse e valida o Token JWT
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Valida se o algoritmo de assinatura é o esperado (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		// 4. Se for válido, extrai as Claims e injeta no Contexto do Gin
		if claims, ok := token.Claims.(*CustomClaims); ok {
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next() // Permite continuar para a rota protegida
	}
}
