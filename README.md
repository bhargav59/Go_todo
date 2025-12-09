# 📝 Go Todo REST API

A production-ready, scalable RESTful API for task/todo management built with Go, Gin framework, and JWT authentication.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Tests](https://img.shields.io/badge/Tests-15%2F15%20Passing-success)

## ✨ Features

- **🔐 JWT Authentication** - Secure user registration and login
- **📋 Full CRUD Operations** - Create, read, update, delete todos
- **👤 User Ownership** - Users can only access their own todos
- **📄 Pagination** - Efficient listing with page/per_page support
- **🔍 Filtering** - Filter todos by completion status
- **📊 Statistics** - Get todo stats (total, completed, pending)
- **⚡ Rate Limiting** - Prevent API abuse
- **📝 Structured Logging** - Request tracking with unique IDs
- **🐳 Docker Ready** - Dockerfile and docker-compose included
- **🗄️ Dual Database Support** - PostgreSQL (production) / SQLite (development)

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, SQLite used by default)
- Docker & Docker Compose (optional)

### Local Development

```bash
# Clone the repository
git clone https://github.com/bhargav59/Go_todo.git
cd Go_todo

# Install dependencies
make deps

# Run the server (uses SQLite)
make run
```

The API will be available at `http://localhost:8080`

### Using Docker

```bash
# Start with PostgreSQL
make docker-up

# View logs
make docker-logs

# Stop containers
make docker-down
```

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Register new user | ❌ |
| POST | `/api/auth/login` | Login and get JWT | ❌ |
| GET | `/api/auth/profile` | Get current user profile | ✅ |

### Todos

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/todos` | Create a new todo | ✅ |
| GET | `/api/todos` | List all todos (paginated) | ✅ |
| GET | `/api/todos/:id` | Get a specific todo | ✅ |
| PUT | `/api/todos/:id` | Update a todo | ✅ |
| DELETE | `/api/todos/:id` | Delete a todo | ✅ |
| GET | `/api/todos/stats` | Get todo statistics | ✅ |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | API health status |

## 🔧 Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project",
    "description": "Finish the Go Todo API",
    "priority": "high"
  }'
```

### List Todos with Pagination

```bash
curl "http://localhost:8080/api/todos?page=1&per_page=10&completed=false" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📁 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── handlers/             # HTTP request handlers
│   ├── middleware/           # Auth, logging, rate limiting
│   ├── models/               # Database models
│   ├── repository/           # Data access layer
│   └── services/             # Business logic
├── pkg/
│   ├── database/             # Database connection
│   └── utils/                # JWT, response helpers
├── tests/                    # Integration tests
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # Docker Compose config
├── Makefile                  # Build commands
└── README.md
```

## ⚙️ Configuration

Environment variables (see `.env.example`):

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | Server port |
| `ENVIRONMENT` | development | Environment (development/production) |
| `DB_HOST` | sqlite | Database host (use `sqlite` for SQLite) |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | todo_api | Database name |
| `JWT_SECRET` | (required) | JWT signing secret |
| `JWT_EXPIRY` | 86400 | Token expiry in seconds (24h) |

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race
```

## 🛠️ Available Make Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the application |
| `make build` | Build binary to `./bin/` |
| `make test` | Run all tests |
| `make test-coverage` | Generate coverage report |
| `make clean` | Remove build artifacts |
| `make deps` | Download dependencies |
| `make docker-build` | Build Docker image |
| `make docker-up` | Start Docker containers |
| `make docker-down` | Stop Docker containers |
| `make swagger` | Generate Swagger docs |

## 📦 Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **ORM**: [GORM](https://gorm.io/) - Go ORM library
- **Database**: PostgreSQL / SQLite
- **Authentication**: JWT (golang-jwt/jwt)
- **Testing**: testify
- **Container**: Docker & Docker Compose

## 🔒 Security Features

- Password hashing with bcrypt
- JWT token authentication
- Rate limiting (100 requests/minute per IP)
- Input validation
- SQL injection prevention via GORM
- CORS support

## 📈 Performance

- Connection pooling for database
- Efficient pagination
- Goroutine-safe rate limiter
- Graceful shutdown support

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 👨‍💻 Author

**Bhaskar**

- GitHub: [@bhargav59](https://github.com/bhargav59)

---

⭐ Star this repo if you find it helpful!
