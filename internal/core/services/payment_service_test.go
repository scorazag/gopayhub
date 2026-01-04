package services

import (
	"encoding/json"
	"testing"

	"github.com/scorazag/gopayhub/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo ahora implementará TODOS los métodos de ports.PaymentRepository
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetMerchantByID(id uint) (*domain.Merchant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Merchant), args.Error(1)
}

func (m *MockRepo) CreateTransaction(tx *domain.Transaction) error {
	return m.Called(tx).Error(0)
}

func (m *MockRepo) GetClientByApiKey(apiKey string) (*domain.Client, error) {
	args := m.Called(apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockRepo) GetIdempotencyKey(key string) (*domain.IdempotencyKey, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdempotencyKey), args.Error(1)
}

func (m *MockRepo) SaveIdempotencyKey(key *domain.IdempotencyKey) error {
	return m.Called(key).Error(0)
}

func (m *MockRepo) CreateDeposit(deposit *domain.Deposit) error {
	args := m.Called(deposit)
	return args.Error(0)
}

func (m *MockRepo) CreateCashOut(cashout *domain.CashOut) error {
	args := m.Called(cashout)
	return args.Error(0)
}

func (m *MockRepo) GetClientBalance(clientID uint) (float64, error) {
	args := m.Called(clientID)
	return args.Get(0).(float64), args.Error(1)
}

// --- TEST 1: MONTO CERO ---
func TestProcessPayment_AmountZero(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewPaymentService(mockRepo)

	// No necesitamos configurar mocks aquí porque el código falla ANTES de tocar el repo
	tx, err := service.ProcessPayment(0, 1, 1, "REF-123", "")

	assert.Nil(t, tx)
	assert.Equal(t, "el monto debe ser mayor a cero", err.Error())
}

// --- TEST 2: IDEMPOTENCIA (Llave existente) ---
func TestProcessPayment_IdempotencyHit(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewPaymentService(mockRepo)

	// Preparamos una transacción vieja "guardada" en JSON
	oldTx := domain.Transaction{Amount: 100, Reference: "PAGO-ANTERIOR"}
	oldTxJSON, _ := json.Marshal(oldTx)

	existingKey := &domain.IdempotencyKey{
		Key:          "key-repetida",
		ResponseJSON: string(oldTxJSON),
	}

	// Configuramos el mock: "Cuando pregunten por esta llave, devuélvela"
	mockRepo.On("GetIdempotencyKey", "key-repetida").Return(existingKey, nil)

	// Ejecución
	tx, err := service.ProcessPayment(100, 1, 1, "REF-123", "key-repetida")

	// Aserciones
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "PAGO-ANTERIOR", tx.Reference)

	// Verificamos que NO se intentó crear una nueva transacción
	mockRepo.AssertNotCalled(t, "CreateTransaction", mock.Anything)
}

func TestProcessPayment_SuccessNewKey(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewPaymentService(mockRepo)

	merchant := &domain.Merchant{ID: 1, Name: "Test Merchant"}
	idemKey := "nueva-llave-123"

	// 1. Mock: No existe la llave todavía
	mockRepo.On("GetIdempotencyKey", idemKey).Return(nil, nil)

	// 2. Mock: El merchant existe
	mockRepo.On("GetMerchantByID", uint(1)).Return(merchant, nil)

	// 3. Mock: Se crea la transacción (usamos Anything porque el UUID se genera adentro)
	mockRepo.On("CreateTransaction", mock.Anything).Return(nil)

	// 4. Mock: Se guarda la llave de idempotencia
	mockRepo.On("SaveIdempotencyKey", mock.Anything).Return(nil)

	// Ejecución
	tx, err := service.ProcessPayment(150.0, 1, 1, "REF-ABC", idemKey)

	// Aserciones
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 150.0, tx.Amount)

	// Verificamos que se llamaron a los métodos de guardado
	mockRepo.AssertExpectations(t)
}
