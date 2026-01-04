package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

type DepositRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reference string  `json:"reference" binding:"required"`
	StoreName string  `json:"store_name"` // Opcional: El nombre del punto de venta
}

type DepositHandler struct {
	service ports.DepositService
}

func NewDepositHandler(service ports.DepositService) *DepositHandler {
	return &DepositHandler{service: service}
}

func (h *DepositHandler) ProcessDeposit(c *gin.Context) {
	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Obtenemos el client_id del Middleware de Auth
	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo identificar al cliente"})
		return
	}

	idemKey := c.GetHeader("X-Idempotency-Key")

	// Llamamos al servicio (aquí pasamos 0 o un valor por defecto para merchantID si no aplica)
	res, err := h.service.ProcessDeposit(req.Amount, 0, clientID.(uint), req.Reference, idemKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
