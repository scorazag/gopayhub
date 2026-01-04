package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scorazag/gopayhub/internal/adapters/repository/postgres"
)

func AuthMiddleware(repo *postgres.PaymentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer la key del header X-API-KEY
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere API Key"})
			c.Abort()
			return
		}

		// 2. Validar contra la base de datos
		client, err := repo.GetClientByApiKey(apiKey)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "API Key inválida o cliente inactivo"})
			c.Abort()
			return
		}

		// 3. Guardar el cliente en el contexto por si lo necesitamos en el controlador
		c.Set("client_id", client.ID)
		c.Set("client_name", client.Name)

		c.Next()
	}
}
