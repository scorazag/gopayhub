package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

type CashOutRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reference string  `json:"reference" binding:"required"`
	StoreName string  `json:"store_name"`
}

type CashOutHandler struct {
	service ports.CashOutService
}

func NewCashOutHandler(service ports.CashOutService) *CashOutHandler {
	return &CashOutHandler{service: service}
}

func (h *CashOutHandler) ProcessCashOut(c *gin.Context) {
	var req CashOutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Recuperamos el client_id del middleware de autenticación
	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo identificar al cliente"})
		return
	}

	idemKey := c.GetHeader("X-Idempotency-Key")

	// Ejecutamos el retiro
	res, err := h.service.ProcessCashOut(req.Amount, 0, clientID.(uint), req.Reference, idemKey)
	if err != nil {
		// Si el error es "insufficient funds", regresamos un 422 (Unprocessable Entity)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
