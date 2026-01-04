package ports

import (
	"github.com/scorazag/gopayhub/internal/core/domain"
)

// PaymentRepository define qué puede hacer la base de datos
type PaymentRepository interface {
	GetClientByApiKey(apiKey string) (*domain.Client, error)
	GetMerchantByID(id uint) (*domain.Merchant, error)
	CreateTransaction(tx *domain.Transaction) error
	CreateDeposit(tx *domain.Deposit) error
	GetIdempotencyKey(key string) (*domain.IdempotencyKey, error)
	SaveIdempotencyKey(key *domain.IdempotencyKey) error
	GetClientBalance(clientID uint) (float64, error) // <--- El nuevo superpoder
	CreateCashOut(cashout *domain.CashOut) error
}

// PaymentService define qué lógica de negocio exponemos
type PaymentService interface {
	ProcessPayment(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.Transaction, error)
}

// DepositService - Contrato exclusivo para depósitos
type DepositService interface {
	ProcessDeposit(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.Deposit, error)
}

// CashOutService - Contrato exclusivo para retiros
type CashOutService interface {
	ProcessCashOut(amount float64, merchantID uint, clientID uint, reference string, idemKey string) (*domain.CashOut, error)
}
