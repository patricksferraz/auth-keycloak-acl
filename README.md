# Auth Service with Keycloak Integration

[![Go Report Card](https://goreportcard.com/badge/github.com/patricksferraz/auth-service)](https://goreportcard.com/report/github.com/patricksferraz/auth-service)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/patricksferraz/auth-service?status.svg)](https://godoc.org/github.com/patricksferraz/auth-service)

A robust authentication and authorization service built with Go, featuring Keycloak integration, gRPC support, and comprehensive API documentation.

## 🌟 Features

- 🔐 Keycloak Integration for Identity and Access Management
- 🚀 gRPC Support for High-Performance Communication
- 📚 Swagger/OpenAPI Documentation
- 🔄 Kafka Integration for Event-Driven Architecture
- 📊 Elastic APM Integration for Performance Monitoring
- 🐳 Docker and Kubernetes Support
- 🧪 Comprehensive Test Suite
- 🔒 CORS Support
- 📝 Structured Logging with Logrus

## 🏗️ Architecture

This service follows a clean architecture pattern with the following layers:

- **Domain**: Core business logic and entities
- **Application**: Use cases and business rules
- **Infrastructure**: External services integration
- **API**: HTTP/gRPC endpoints and controllers

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose
- Keycloak instance
- Kafka (optional, for event-driven features)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/patricksferraz/auth-service.git
cd auth-service
```

2. Copy the environment file and configure it:
```bash
cp .env.example .env
```

3. Build and run with Docker Compose:
```bash
docker-compose up -d
```

### Development

1. Install dependencies:
```bash
go mod download
```

2. Run the service:
```bash
go run main.go
```

## 📚 API Documentation

Once the service is running, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## 🧪 Testing

Run the test suite:
```bash
make test
```

For integration tests:
```bash
docker-compose -f docker-compose.test.yml up
```

## 🛠️ Configuration

The service can be configured through environment variables or configuration files. See `.env.example` for available options.

## 📦 Deployment

### Docker

Build the Docker image:
```bash
docker build -t auth-service .
```

### Kubernetes

Kubernetes manifests are available in the `k8s/` directory.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Keycloak](https://www.keycloak.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [gRPC](https://grpc.io/)
- [Kafka](https://kafka.apache.org/)
- [Elastic APM](https://www.elastic.co/apm)

## 📫 Contact

Project Link: [https://github.com/patricksferraz/auth-service](https://github.com/patricksferraz/auth-service)
