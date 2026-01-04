package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessCashOut_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewCashOutService(mockRepo)

	// Mockeamos: El cliente tiene $1000 y el guardado es exitoso
	mockRepo.On("GetClientBalance", uint(1)).Return(1000.0, nil)
	mockRepo.On("CreateCashOut", mock.Anything).Return(nil)

	res, err := service.ProcessCashOut(200.0, 0, 1, "REF-CASH-01", "idem-999")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 200.0, res.Amount)
	assert.Equal(t, "COMPLETED", res.Status)
	mockRepo.AssertExpectations(t)
}

func TestProcessCashOut_InsufficientFunds(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewCashOutService(mockRepo)

	// Mockeamos: El cliente solo tiene $50
	mockRepo.On("GetClientBalance", uint(1)).Return(50.0, nil)

	// Intenta sacar $100
	res, err := service.ProcessCashOut(100.0, 0, 1, "REF-CASH-02", "")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "insufficient funds", err.Error())
	// Verificamos que no se intentó guardar nada
	mockRepo.AssertNotCalled(t, "CreateCashOut", mock.Anything)
}

func TestProcessCashOut_AmountZero(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewCashOutService(mockRepo)

	_, err := service.ProcessCashOut(-10.0, 0, 1, "REF-CASH-03", "")

	assert.Error(t, err)
	assert.Equal(t, "el monto debe ser mayor a cero", err.Error())
}
