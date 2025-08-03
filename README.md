# ShoppinggCart

# E-commerce API

A simple e-commerce REST API built with Go, Gin, and GORM.

## Features

- User registration and authentication
- Shopping cart management
- Order processing
- Item listing
- Single device login (one token per user)
- SQLite database (easily switchable to PostgreSQL)

## Tech Stack

- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Object Relational Mapping)
- **Database**: SQLite (development), PostgreSQL (production ready)
- **Testing**: Ginkgo & Gomega
- **Authentication**: Token-based (Bearer tokens)

## API Endpoints

### User Management
- `POST /users` - Create new user account
- `POST /users/login` - User login (returns token)
- `GET /users` - List all users

### Cart Management (Protected)
- `POST /carts` - Add item to cart
- `GET /carts` - Get user's cart
- `GET /carts/all` - List all carts (admin)

### Order Management (Protected)
- `POST /orders` - Create order from cart
- `GET /orders` - List user's orders

### Items
- `GET /items` - List all items (with pagination)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
\`\`\`bash
git clone <repository-url>
cd ecommerce-api
\`\`\`

2. Install dependencies:
\`\`\`bash
go mod download
\`\`\`

3. Run the application:
\`\`\`bash
go run .
\`\`\`

The server will start on `http://localhost:8080`

### Using Make Commands

\`\`\`bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with Ginkgo
make test-ginkgo

# Format code
make fmt

# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run
\`\`\`

## Usage Examples

### 1. Create a User
\`\`\`bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
\`\`\`

### 2. Login
\`\`\`bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
\`\`\`

### 3. Add Item to Cart
\`\`\`bash
curl -X POST http://localhost:8080/carts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "item_id": 1,
    "quantity": 2
  }'
\`\`\`

### 4. Create Order
\`\`\`bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
\`\`\`

### 5. List Items
\`\`\`bash
curl -X GET "http://localhost:8080/items?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
\`\`\`

## Testing

Run tests using Ginkgo:

\`\`\`bash
# Install Ginkgo CLI
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Run tests
ginkgo -v
\`\`\`

Or use standard Go testing:

\`\`\`bash
go test -v ./...
\`\`\`

## Database Schema

The application uses the following models:

- **User**: User accounts with authentication
- **Item**: Products available for purchase
- **Cart**: User shopping carts
- **CartItem**: Items in a cart with quantities
- **Order**: Completed orders
- **OrderItem**: Items in an order with price snapshots

## Configuration

Environment variables:
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)

## Production Deployment

### Docker

1. Build the Docker image:
\`\`\`bash
docker build -t ecommerce-api .
\`\`\`

2. Run with Docker Compose:
\`\`\`bash
docker-compose up -d
\`\`\`

### Database Migration

For production, consider switching to PostgreSQL:

1. Update the database connection in `main.go`
2. Update the Docker Compose file to include PostgreSQL
3. Run migrations

## Security Considerations

- Passwords are hashed using bcrypt
- Token-based authentication
- Single device login enforcement
- Input validation on all endpoints
- SQL injection protection via GORM

## Future Enhancements

- JWT tokens instead of simple tokens
- Role-based access control
- Inventory management
- Payment integration
- Email notifications
- Rate limiting
- API documentation with Swagger
- Caching with Redis
- Logging and monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and ensure they pass
6. Submit a pull request


