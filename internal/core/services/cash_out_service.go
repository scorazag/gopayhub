package services

import (
	"errors"

	"github.com/scorazag/gopayhub/internal/core/domain"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

type CashOutService struct {
	repo ports.PaymentRepository
}

func NewCashOutService(repo ports.PaymentRepository) ports.CashOutService {
	return &CashOutService{repo: repo}
}

func (s *CashOutService) ProcessCashOut(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.CashOut, error) {
	// 1. Validar que el monto sea positivo
	if amount <= 0 {
		return nil, errors.New("el monto debe ser mayor a cero")
	}
	// 2. OBTENER EL SALDO ACTUAL usando el repo
	balance, err := s.repo.GetClientBalance(clientID)
	if err != nil {
		return nil, err
	}
	// 3. VALIDACIÓN CLAVE: ¿Tiene dinero suficiente?
	if amount > balance {
		return nil, errors.New("insufficient funds")
	}
	// 4. Crear el objeto CashOut
	cashout := &domain.CashOut{
		Amount:    amount,
		ClientID:  clientID,
		Reference: reference,
		Status:    "COMPLETED",
	}
	// 5. Guardar en el repo
	err = s.repo.CreateCashOut(cashout)
	if err != nil {
		return nil, err
	}
	return cashout, nil
}
