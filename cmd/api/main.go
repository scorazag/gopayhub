package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	// 1. Le damos el alias 'gormPostgres' al driver oficial
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/scorazag/gopayhub/internal/adapters/handler/http"
	"github.com/scorazag/gopayhub/internal/adapters/handler/http/middleware"

	// 2. Le damos el alias 'repoPostgres' a TU carpeta
	repoPostgres "github.com/scorazag/gopayhub/internal/adapters/repository/postgres"

	"github.com/scorazag/gopayhub/internal/core/domain"
	"github.com/scorazag/gopayhub/internal/core/services"
)

func main() {
	// 1. Conexión a la Base de Datos (DB)
	dsn := "host=localhost user=user password=password dbname=gopayhub port=5432 sslmode=disable TimeZone=America/Mexico_City"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// GORM hace un buen manejo de logs, útil para debugging
		Logger: nil,
	})
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Configuración opcional de pool de conexiones
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 2. Migraciones
	log.Println("Ejecutando migraciones...")
	// Habilitar extensión uuid-ossp si no existe
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatalf("Error al habilitar la extensión uuid-ossp: %v", err)
	}

	err = db.AutoMigrate(
		&domain.Client{},
		&domain.Merchant{},
		&domain.Transaction{},
		&domain.IdempotencyKey{},
		&domain.Deposit{},
	)
	if err != nil {
		log.Fatalf("Error durante la migración de la DB: %v", err)
	}
	log.Println("Migraciones completadas exitosamente.")

	// 3. Inicialización de la Arquitectura Hexagonal (Inyección de Dependencias)

	// Repositorio (Capa de Infraestructura)
	// El repositorio solo sabe interactuar con la DB
	repo := repoPostgres.NewPaymentRepository(db)

	// Servicio (Capa de Core/Negocio)
	// El servicio recibe el repositorio, NO la DB.
	paymentService := services.NewPaymentService(repo)
	depositService := services.NewDepositService(repo)

	// Handler (Capa de Adaptadores/Gin)
	// El handler recibe el servicio.
	paymentHandler := http.NewPaymentHandler(paymentService)
	depositHandler := http.NewDepositHandler(depositService)

	// 4. Configuración de Rutas y Servidor Gin

	// Grupo de rutas API
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// El health check lo dejamos fuera del auth para que AWS/Docker puedan revisarlo
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK"})
		})

		// Aplicamos el middleware a partir de aquí
		api.Use(middleware.AuthMiddleware(repo))

		// Esta ruta ahora está protegida
		api.POST("/transactions", paymentHandler.ProcessTransaction)
		api.POST("/deposits", depositHandler.ProcessDeposit)
	}

	log.Println("Servidor GoPayHub iniciado en :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
