# TrafficController

A Go application that serves JSON content based on UUID lookups in a MySQL database.

## Features

- Retrieve JSON content by UUID via RESTful API
- Health check endpoint for monitoring
- Clean architecture with proper separation of concerns
- Production-ready with Docker support

## Architecture

- **Domain Layer**: Core business logic and rules
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External concerns like databases
- **Interface Layer**: User interfaces like REST APIs

The application is divided into bounded contexts:

- **Traffic Domain**: Manages the storage and retrieval of traffic data
- **Health Domain**: Handles service health monitoring
- **Common Domain**: Shared domain concepts and utilities

## Project Structure

```
traffic-controller/
├── cmd/
│   └── server/            # Application entry point
├── internal/
│   ├── common/            # Shared components
│   ├── traffic/           # Traffic domain
│   └── health/            # Health check domain
├── pkg/                   # Public packages
├── api/                   # API definitions
└── docker-compose.yml     # Docker Compose configuration
```

## API Endpoints

- `GET /serve/{uuid}` - Returns JSON content for the specified UUID
- `GET /health` - Health check endpoint

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Docker (optional)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/githamo/stubhub-tc.git traffic-controller
   cd traffic-controller
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create the database schema:
   ```sql
   CREATE DATABASE IF NOT EXISTS trafficdb;
   USE trafficdb;

   CREATE TABLE endpoints (
       id INT AUTO_INCREMENT PRIMARY KEY,
       uuid VARCHAR(36) UNIQUE NOT NULL,
       content VARCHAR(40) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   );

   CREATE TABLE stub_contents (
       id INT AUTO_INCREMENT PRIMARY KEY,
       filename VARCHAR(64) NOT NULL,
       content JSON NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   );
   ```

4. Configure environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials & base64 encryption key
   ```

5. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Docker Deployment

1. Build and run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

## Testing

Run tests:
```bash
go test ./...
```

## Usage Examples

### Retrieve content by UUID

```bash
curl http://localhost:8080/serve/550e8400-e29b-41d4-a716-446655440000
```

### Check service health

```bash
curl http://localhost:8080/health
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.