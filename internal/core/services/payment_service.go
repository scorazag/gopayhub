package services

import (
	// Importante añadir esto
	"encoding/json"
	"errors"

	"github.com/scorazag/gopayhub/internal/core/domain"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

type paymentService struct {
	repo ports.PaymentRepository // Aquí guardamos la interfaz
}

// Constructor del servicio
func NewPaymentService(repo ports.PaymentRepository) ports.PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) ProcessPayment(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.Transaction, error) {

	// 0. BUSCAR IDEMPOTENCIA
	if idemKey != "" {
		existingKey, err := s.repo.GetIdempotencyKey(idemKey)
		// Si no hay error y encontramos la llave...
		if err == nil && existingKey != nil && existingKey.Key != "" {
			var oldTx domain.Transaction
			// Convertimos el JSON guardado en la DB de nuevo a un struct de Transacción
			if err := json.Unmarshal([]byte(existingKey.ResponseJSON), &oldTx); err == nil {
				return &oldTx, nil // Retornamos la transacción original sin hacer nada más
			}
		}
	}

	// 1. REGLAS DE NEGOCIO
	if amount <= 0 {
		return nil, errors.New("el monto debe ser mayor a cero")
	}

	// 2. VERIFICAR MERCHANT
	merchant, err := s.repo.GetMerchantByID(merchantID)
	if err != nil {
		return nil, errors.New("proveedor de servicio no encontrado")
	}

	// 3. CREAR OBJETO TRANSACCIÓN
	tx := &domain.Transaction{
		Amount:         amount,
		MerchantID:     merchant.ID,
		ClientID:       clientID,
		Reference:      reference,
		Status:         "COMPLETED",
		IdempotencyKey: idemKey,
	}

	// 4. GUARDAR TRANSACCIÓN PRIMERO
	if err := s.repo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	// 5. GUARDAR LLAVE DE IDEMPOTENCIA CON EL RESULTADO
	if idemKey != "" {
		// Convertimos la transacción real (ya con su ID y fecha) a JSON
		txJSON, _ := json.Marshal(tx)

		newIdem := &domain.IdempotencyKey{
			Key:          idemKey,
			ResponseJSON: string(txJSON),
			StatusCode:   201,
		}
		// Guardamos la llave para futuros reintentos
		_ = s.repo.SaveIdempotencyKey(newIdem)
	}

	return tx, nil
}
