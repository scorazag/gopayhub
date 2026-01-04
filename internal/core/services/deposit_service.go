package services

import (
	"errors"

	"github.com/scorazag/gopayhub/internal/core/domain"
	"github.com/scorazag/gopayhub/internal/core/ports"
)

type depositService struct {
	repo ports.PaymentRepository
}

func NewDepositService(repo ports.PaymentRepository) ports.DepositService {
	return &depositService{repo: repo}
}

func (s *depositService) ProcessDeposit(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.Deposit, error) {
	// 1. Validaciones
	if amount > 10000 {
		return nil, errors.New("el monto excede el límite permitido para depósitos en efectivo")
	}
	if amount <= 0 {
		return nil, errors.New("monto inválido")
	}

	// 2. Crear objeto
	deposit := &domain.Deposit{
		Amount:    amount,
		ClientID:  clientID,
		Reference: reference,
		Status:    "COMPLETED",
	}

	// 3. Guarda el deposito
	err := s.repo.CreateDeposit(deposit)
	if err != nil {
		return nil, err
	}

	return deposit, nil
}
