# GoPayHub - Payment Gateway MVP

GoPayHub es una micro-pasarela de pagos de alto rendimiento desarrollada en **Go (Golang)** siguiendo los principios de **Arquitectura Hexagonal (Ports & Adapters)**. El sistema permite a comercios (como Oxxo o Walmart) procesar pagos de servicios (luz, agua, internet) de forma segura y eficiente.

## 🚀 Características Principales
- **Arquitectura Hexagonal**: Separación clara entre lógica de negocio, puertos e infraestructura.
- **Autenticación mediante API Key**: Middleware personalizado para validación de clientes en base de datos.
- **Persistencia con GORM**: Integración robusta con PostgreSQL.
- **Idempotencia**: Prevención de cobros duplicados mediante llaves únicas por transacción.
- **Dockerizado**: Entorno de desarrollo listo con Docker Compose.

---

## 🏗️ Estructura del Proyecto

```text
├── cmd
│   └── api/main.go           # Punto de entrada. Configura dependencias y arranca el server.
├── internal
│   ├── adapters              # Implementaciones externas (Infraestructura)
│   │   ├── handler/http      # Controladores Gin y Middleware de seguridad.
│   │   └── repository/postgres # Implementación de BD con GORM.
│   ├── core                  # El corazón de la aplicación
│   │   ├── domain            # Modelos y entidades de negocio (Transactions, Clients).
│   │   ├── ports             # Interfaces (Contratos) que definen el comportamiento.
│   │   └── services          # Lógica de negocio y reglas de validación.
│   └── pkg                   # Librerías compartidas y utilidades.
├── docker-compose.yml        # Configuración de PostgreSQL.
└── go.mod                    # Gestión de dependencias de Go.


CLIENTE (OXXO) 
    │ 
    ▼ [Request + API Key]
MIDDLEWARE (Seguridad) ────▶ [Busca Cliente en DB]
    │ 
    ▼ [Pasa petición limpia]
HANDLER (HTTP Adapter) ────▶ [Valida Formato JSON]
    │ 
    ▼ [Llamada a método]
SERVICE (Core/Negocio) ────▶ [Aplica Reglas: ¿Monto > 0?]
    │ 
    ▼ [Llamada a interfaz]
REPOSITORY (DB Adapter) ───▶ [SQL INSERT en Postgres]
    │ 
    └────────── REGRESA INFORMACIÓN ───────────┐
                                               │
CLIENTE ◀──── [JSON + 201 OK] ◀──── HANDLER ◀──┘


# GoPayHub 🚀

Procesador de pagos y movimientos financieros construido con **Arquitectura Hexagonal** en Go.

### Características:
* **Pagos:** Procesamiento de transacciones con Merchants.
* **Depósitos:** Carga de saldo en efectivo (Límite $10,000).
* **Cash-Out:** Retiros de efectivo con validación de saldo en tiempo real.
* **Idempotencia:** Seguridad en transacciones duplicadas mediante Headers.
* **Tecnologías:** Gin Gonic, GORM, Postgres y Unit Testing (Testify).

### Cómo correrlo:
1. `go mod tidy`
2. Configurar DSN en `main.go`
3. `go run cmd/api/main.go`