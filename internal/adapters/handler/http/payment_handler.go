package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

// Definimos el struct para leer el JSON que viene de afuera
type PaymentRequest struct {
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	MerchantID uint    `json:"merchant_id" binding:"required"`
	Reference  string  `json:"reference" binding:"required"`
}

type PaymentHandler struct {
	service ports.PaymentService
}

func NewPaymentHandler(service ports.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) ProcessTransaction(c *gin.Context) {
	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al identificar al cliente"})
		return
	}

	// 3. Obtener Idempotency-Key del Header
	// Es estándar usar "X-Idempotency-Key" o "Idempotency-Key"
	idemKey := c.GetHeader("X-Idempotency-Key")

	// 4. Llamar al servicio
	tx, err := h.service.ProcessPayment(
		req.Amount,
		req.MerchantID,
		clientID.(uint),
		req.Reference,
		idemKey,
	)

	if err != nil {
		// Si el error es de negocio, devolvemos 422
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// 5. Éxito: Devolvemos el objeto 'tx' completo.
	// Esto es vital para que en reintentos de idempotencia el cliente reciba
	// la misma data que la primera vez.
	c.JSON(http.StatusCreated, tx)
}
