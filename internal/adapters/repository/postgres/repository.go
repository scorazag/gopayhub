package postgres

import (
	"github.com/scorazag/gopayhub/internal/core/domain"
	"gorm.io/gorm"
)

// PaymentRepository implementa la interfaz de puertos
type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetMerchantByID(id uint) (*domain.Merchant, error) {
	var merchant domain.Merchant
	err := r.db.First(&merchant, id).Error
	return &merchant, err
}

func (r *PaymentRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *PaymentRepository) CreateDeposit(tx *domain.Deposit) error {
	return r.db.Create(tx).Error
}

func (r *PaymentRepository) GetIdempotencyKey(key string) (*domain.IdempotencyKey, error) {
	var idempotencyKey domain.IdempotencyKey
	err := r.db.Where("key = ?", key).First(&idempotencyKey).Error
	return &idempotencyKey, err
}

func (r *PaymentRepository) SaveIdempotencyKey(key *domain.IdempotencyKey) error {
	return r.db.Create(key).Error
}

func (r *PaymentRepository) GetClientByApiKey(apiKey string) (*domain.Client, error) {
	var client domain.Client
	// Buscamos un cliente activo que coincida con la key
	err := r.db.Where("api_key = ? AND is_active = ?", apiKey, true).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *PaymentRepository) GetClientBalance(clientID uint) (float64, error) {
	var totalDeposits float64
	var totalPayments float64
	var totalCashOuts float64

	// Sumar Depósitos
	r.db.Model(&domain.Deposit{}).Where("client_id = ? AND status = ?", clientID, "COMPLETED").Select("COALESCE(sum(amount), 0)").Scan(&totalDeposits)

	// Sumar Pagos (Transactions)
	r.db.Model(&domain.Transaction{}).Where("client_id = ? AND status = ?", clientID, "COMPLETED").Select("COALESCE(sum(amount), 0)").Scan(&totalPayments)

	// Sumar Retiros (CashOuts)
	r.db.Model(&domain.CashOut{}).Where("client_id = ? AND status = ?", clientID, "COMPLETED").Select("COALESCE(sum(amount), 0)").Scan(&totalCashOuts)

	// Saldo final
	return totalDeposits - totalPayments - totalCashOuts, nil
}

func (r *PaymentRepository) CreateCashOut(cashout *domain.CashOut) error {
	return r.db.Create(cashout).Error
}
