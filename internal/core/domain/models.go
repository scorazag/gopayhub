package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`    // Ej: "Oxxo Sucursal Centro"
	ApiKey    string `gorm:"uniqueIndex;not null"` // Ej: "sk_live_12345"
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
}

type Merchant struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"size:100"` // Ej: "CFE", "Netflix"
	ServiceType    string `gorm:"index"`    // Ej: "ELECTRICITY", "STREAMING"
	IntegrationURL string
	CreatedAt      time.Time
}

type Transaction struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Amount         float64   `gorm:"not null"`
	Currency       string    `gorm:"size:3;default:'MXN'"`
	Status         string    `gorm:"size:20;index"` // PENDING, COMPLETED, FAILED
	Reference      string    `gorm:"not null"`      // Referencia del recibo de luz
	ClientID       uint
	Client         Client
	MerchantID     uint
	Merchant       Merchant
	CreatedAt      time.Time
	IdempotencyKey string `gorm:"size:100;index"` // Relación lógica
}

// Tabla para evitar doble cobro
type IdempotencyKey struct {
	Key          string `gorm:"primaryKey"` // El UUID que manda Oxxo
	ResponseJSON string // Guardamos qué le respondimos la primera vez
	StatusCode   int    // 200, 400, etc.
	CreatedAt    time.Time
}

type Deposit struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Amount         float64   `gorm:"not null"`
	Currency       string    `gorm:"size:3;default:'MXN'"`
	Status         string    `gorm:"size:20;index"` // PENDING, COMPLETED, FAILED
	Reference      string    `gorm:"not null"`
	StoreName      string    `gorm:"size:100"`       // Ej: "OXXO Tacubaya"
	ExternalID     string    `gorm:"size:100;index"` // ID que te da el corresponsal
	ClientID       uint      `gorm:"not null"`
	Client         Client    `gorm:"foreignKey:ClientID"`
	CreatedAt      time.Time
	IdempotencyKey string `gorm:"size:100;index"`
}

type CashOut struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Amount         float64   `gorm:"not null"`
	Currency       string    `gorm:"size:3;default:'MXN'"`
	Status         string    `gorm:"size:20;index"` // PENDING, COMPLETED, FAILED
	Reference      string    `gorm:"not null"`
	StoreName      string    `gorm:"size:100"`       // Ej: "OXXO Tacubaya"
	ExternalID     string    `gorm:"size:100;index"` // ID que te da el corresponsal
	ClientID       uint      `gorm:"not null"`
	Client         Client    `gorm:"foreignKey:ClientID"`
	CreatedAt      time.Time
	IdempotencyKey string `gorm:"size:100;index"`
}

func (d *Deposit) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return nil
}

func (c *CashOut) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	// Aquí podrías forzar el Status inicial si quisieras
	if c.Status == "" {
		c.Status = "PENDING"
	}
	return nil
}

// Hook BeforeCreate: Se ejecuta automáticamente antes de insertar en la DB
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	// Generamos el UUID desde Go
	t.ID = uuid.New()
	// También podrías poner validaciones básicas aquí
	// por ejemplo, asegurar que el monto sea positivo
	return nil
}
