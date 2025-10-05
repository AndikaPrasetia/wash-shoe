# WashShoe - Laundry Management System

WashShoe is a web-based laundry management system built with Go, PostgreSQL, and Redis. The system provides features for managing laundry orders, users, and other laundry business processes.

## Key Features

- **Authentication & Authorization**: Login, signup, logout, and refresh token system with JWT
- **User Management**: CRUD user operations with role-based authorization
- **Order Management**: Creation, updates, and tracking of laundry order status
- **Redis Cache**: For faster performance and session management
- **Secured API**: API protected with authentication middleware
- **Refresh Token**: Secure refresh token mechanism

## Technologies Used

- **Backend**: Go (Golang)
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT (JSON Web Token)
- **ORM**: SQLC (SQL Compiler)
- **Database Migration**: Custom migration system

## Prerequisites

Make sure you have installed:

- [Go 1.24+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/)

## Installation and Setup

### 1. Clone Repository

```bash
git clone https://github.com/AndikaPrasetia/wash-shoe.git
cd wash-shoe
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Copy the environment file from the example:

```bash
cp example.env .env
```

Then adjust the configurations in `.env` file:

```env
# DB Config
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_DRIVER=postgres

# API Config
API_HOST=localhost
API_PORT=8080

# Token Config
APP_NAME=wash-shoe
JWT_SECRET=your_super_secret_jwt_key
ACCESS_TOKEN_EXP=15
REFRESH_TOKEN_EXP=10080

# Redis Config
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Running the Application

### Method 1: Local Development (Without Docker)

1. **Ensure PostgreSQL and Redis are running**

2. **Setup database**
   ```bash
   # Run database migrations
   # (adjust according to the migration system used in the project)
   ```

3. **Run the application**
   ```bash
   go run cmd/app/main.go
   ```

4. The application will run on `http://localhost:8080`

### Method 2: Docker (Recommended)

1. **Ensure Docker and Docker Compose are installed**

2. **Build and run with Docker Compose**
   ```bash
   docker compose up --build
   ```

3. The application will run on `http://localhost:8080`

4. To run in background:
   ```bash
   docker compose up -d
   ```

## Project Structure

```
wash-shoe/
├── .dockerignore            # Files ignored by Docker
├── .env                     # Environment configuration (don't commit!)
├── .gitignore               # Files ignored by Git
├── Dockerfile               # Dockerfile for application container
├── docker-compose.yml       # Multi-container services configuration
├── go.mod                   # Dependencies management
├── go.sum                   # Dependencies checksum
├── schema.sql               # Database schema
├── Makefile                 # Build commands
├── cmd/                     # Main application
│   └── app/                 # Server entry point
├── internal/                # Internal application code
│   ├── config/              # Application configuration
│   ├── db/                  # Database related code
│   ├── delivery/            # Handlers and HTTP layer
│   ├── dto/                 # Data Transfer Objects
│   ├── middleware/          # HTTP middleware
│   ├── model/               # Data models
│   ├── redis/               # Redis client
│   ├── repository/          # Database query layer
│   ├── sqlc/                # SQLC generated code
│   ├── usecase/             # Business logic
│   └── utils/               # Utility functions
├── migrations/              # Database migrations
└── tmp/                     # Temporary files
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh token

### User Management
- `DELETE /api/v1/users/:id` - Delete user (requires login)

### Home
- `POST /api/v1/home` - Home page (requires login)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL database host | localhost |
| `DB_PORT` | PostgreSQL database port | 5432 |
| `DB_USER` | Database username | - |
| `DB_PASS` | Database password | - |
| `DB_NAME` | Database name | - |
| `API_HOST` | API host | localhost |
| `API_PORT` | API port | 8080 |
| `JWT_SECRET` | JWT secret key | - |
| `ACCESS_TOKEN_EXP` | Access token expiration (minutes) | 15 |
| `REFRESH_TOKEN_EXP` | Refresh token expiration (minutes) | 10080 |
| `REDIS_ADDR` | Redis address | localhost:6379 |
| `REDIS_PASSWORD` | Redis password | empty |
| `REDIS_DB` | Redis database | 0 |

## Development

### Build Project

```bash
go build -o bin/main ./cmd/app
```

### Database Migration

(Note: Adjust according to the migration system used in the project)

## Testing

(Note: Add instructions for running tests if available)

## Deployment

### Deploy to Railway

1. Create an account at [Railway](https://railway.app)
2. Install Railway CLI:
   ```bash
   npm install -g @railway/cli
   ```
3. Login to Railway:
   ```bash
   railway login
   ```
4. Connect to project:
   ```bash
   railway init
   ```
5. Deploy:
   ```bash
   railway up
   ```

### Deploy with Docker

1. Build image:
   ```bash
   docker build -t wash-shoe-app .
   ```
2. Push to registry (if needed):
   ```bash
   docker tag wash-shoe-app username/wash-shoe-app
   docker push username/wash-shoe-app
   ```

## Contributing

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Name: Andika Prasetia
Email: andikaprasetia@proton.me
Project Link: [https://github.com/AndikaPrasetia/wash-shoe](https://github.com/AndikaPrasetia/wash-shoe)

---

