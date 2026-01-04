package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessDeposit_ExceedsLimit(t *testing.T) {
	// Setup
	mockRepo := new(MockRepo) // Usamos el mismo MockRepo que ya tiene CreateDeposit
	service := NewDepositService(mockRepo)

	// Ejecución: Intentamos depositar $11,000 (El límite es 10k)
	res, err := service.ProcessDeposit(11000.0, 0, 1, "DEP-001", "")

	// Aserciones
	assert.Nil(t, res)
	assert.Equal(t, "el monto excede el límite permitido para depósitos en efectivo", err.Error())

	// Verificamos que NUNCA se llamó al repo para guardar
	mockRepo.AssertNotCalled(t, "CreateDeposit", mock.Anything)
}

func TestProcessDeposit_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockRepo)
	service := NewDepositService(mockRepo)

	// Configuramos el mock para que acepte el guardado
	mockRepo.On("CreateDeposit", mock.Anything).Return(nil)

	// Ejecución
	res, err := service.ProcessDeposit(500.0, 0, 1, "DEP-OK", "idem-123")

	// Aserciones
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 500.0, res.Amount)
	assert.Equal(t, "COMPLETED", res.Status)

	// Verificamos que se llamó al guardado exactamente una vez
	mockRepo.AssertExpectations(t)
}
