# Wallet Microservice

A Go-based microservice for managing digital wallets and transactions.

## Features

- Create and manage digital wallets
- Support for multiple currencies (defaults to USD)
- Credit and debit wallet operations
- Transaction history tracking
- PostgreSQL database backend with GORM auto-migration
- RESTful API endpoints

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15+

## Quick Start

### Using Docker Compose

```bash
# Start the application with database
docker-compose up --build

# Run tests in containers
make test-docker

# Clean up
make clean
```

### Local Development

```bash
# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=wallet_db

# Run the application (auto-migrates database)
go run ./cmd/main.go

# Run tests
make test-unit          # Unit tests only
make test-integration   # Integration tests (requires database)
make test              # All tests
```

## Testing

### Test Types

- **Unit Tests**: Fast, isolated tests using mocks
- **Integration Tests**: Tests against real database
- **Containerized Tests**: Full environment testing in Docker

### Running Tests

```bash
# Unit tests (fast, no database needed)
make test-unit

# Integration tests (requires database)
make test-integration

# All tests in Docker containers
make test-docker

# Tests with coverage
make test-coverage

# Tests with race detection
make test-race
```

### Test Tags

- `unit`: Unit tests using mocks
- `integration`: Tests requiring database connection

## API Endpoints

- `POST /wallets` - Create a new wallet
- `GET /wallets/:id` - Get wallet by ID
- `GET /wallets/user/:userID` - Get wallet by user ID
- `PUT /wallets/:id` - Update wallet
- `DELETE /wallets/:id` - Delete wallet
- `POST /wallets/:id/credit` - Credit wallet
- `POST /wallets/:id/debit` - Debit wallet
- `GET /wallets/:id/transactions` - Get transaction history

## Project Structure

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── database/               # Database connection and auto-migration
│   ├── handlers/               # HTTP request handlers
│   ├── models/                 # Data models with GORM tags
│   ├── repositories/           # Data access layer
│   └── services/               # Business logic layer
├── docker-compose.yml          # Main application services
├── docker-compose.test.yml     # Test environment services
├── Dockerfile                  # Main application image
├── Dockerfile.test             # Test environment image
└── Makefile                    # Build and test commands
```

## Database Schema

The database schema is automatically created using GORM auto-migration:

- **wallets**: User wallet information with balance and currency
- **transactions**: Transaction history with credit/debit operations
- **Indexes**: Optimized for user lookups and transaction queries
- **Constraints**: Foreign key relationships and data validation
- **Triggers**: Automatic `updated_at` timestamp updates

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `password` | Database password |
| `DB_NAME` | `wallet_db` | Database name |
| `GIN_MODE` | `debug` | Gin framework mode |

## Development

### Adding New Models

1. **Define the model** in `internal/models/` with GORM tags
2. **Add to AutoMigrate** in `internal/database/connection.go`
3. **Create repository methods** in `internal/repositories/`
4. **Add service logic** in `internal/services/`
5. **Write tests** for new functionality

### GORM Auto-Migration

The application automatically:
- Creates tables based on model structs
- Adds indexes and constraints from GORM tags
- Manages foreign key relationships
- Creates triggers for timestamp updates

### Database Migrations

No manual migrations needed! GORM handles:
- Table creation and updates
- Schema changes
- Index management
- Constraint enforcement

## CI/CD

The project includes GitHub Actions workflows for:
- Running tests on push/PR
- Code coverage reporting
- Integration testing with PostgreSQL service

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

[Add your license here]
