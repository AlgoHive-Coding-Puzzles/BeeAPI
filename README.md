<p align="center">
  <img width="150px" src="https://raw.githubusercontent.com/AlgoHive-Coding-Puzzles/Ressources/refs/heads/main/images/beeapi-logo.png" title="Algohive">
</p>

<h1 align="center">BeeAPI Go</h1>

<p align="center">
  <b>A high-performance, self-hostable API service for AlgoHive puzzle distribution</b>
</p>

## Overview

BeeAPI Go is a robust Go implementation of the BeeAPI service, designed with a microservice architecture for high availability and horizontal scaling. It provides a comprehensive API for loading, managing, and serving puzzles from `.alghive` files to the AlgoHive platform. The service is built with performance and reliability in mind, allowing for independent deployment across multiple environments.

## Key Features

- **Dynamic Puzzle Management**: Automatically extracts, loads, and unloads puzzles from `.alghive` files
- **Stateless Architecture**: Designed for easy replication across multiple instances
- **High Availability Design**: Can be deployed in redundant configurations for zero downtime
- **Secure API Authentication**: API key-based authentication for protected endpoints
- **Comprehensive API**: Full suite of endpoints for puzzle discovery and interaction
- **Swagger Documentation**: Interactive API documentation for easy integration
- **Container-Ready**: Optimized for Docker and container orchestration platforms
- **Efficient Resource Usage**: Intelligent memory management for puzzle loading/unloading

## Technical Architecture

BeeAPI Go implements a clean architecture pattern with clear separation between:

- **Service Layer**: Handles core business logic including puzzle extraction and loading
- **Controller Layer**: Manages API endpoints and request/response handling
- **Model Layer**: Defines data structures for puzzles and themes
- **Middleware Layer**: Provides authentication and request processing

The service is designed to be stateless, making it ideal for horizontal scaling in high availability environments. Multiple BeeAPI instances can be deployed behind a load balancer to distribute traffic and provide redundancy.

## Puzzle Management System

BeeAPI employs a sophisticated puzzle management system:

1. **Initial Loading**: When the server starts, it scans the `puzzles` directory structure
2. **Extraction Process**: `.alghive` files are automatically extracted into their component files
3. **Memory Management**: Puzzles are loaded into memory with optimized resource usage
4. **Graceful Unloading**: On shutdown, puzzle resources are properly released
5. **Dynamic Reloading**: Themes and puzzles can be reloaded without service interruption

This design allows for efficient management of puzzle resources while maintaining high performance.

## AlgoHive Platform

AlgoHive is a web-based, self-hostable platform that enables developers to create and solve coding puzzles. Each puzzle contains two parts to solve, challenging developers to apply different skills and approaches. The puzzles are distributed using the proprietary `.alghive` file format, which BeeAPI is designed to process.

## Installation

### Local Development

```bash
# Clone the repository
git clone https://github.com/AlgoHive-Coding-Puzzles/BeeAPI.git
cd BeeAPI

# Build the application
go build -o beeapi

# Run the application
./beeapi
```

### Docker Deployment

```bash
# Build the Docker image
docker build -t beeapi-go .

# Run as a container
docker run -d -p 8080:8080 --name beeapi-go \
  -v $(pwd)/puzzles:/app/puzzles \
  -v $(pwd)/data:/app/data \
  beeapi-go
```

### High Availability Deployment

For production environments, BeeAPI can be deployed in a high availability configuration:

1. Deploy multiple instances behind a load balancer
2. Use a shared storage volume for puzzle files (NFS, S3, etc.)
3. Implement health checks for automatic instance recovery
4. Configure instance auto-scaling based on load

## Directory Structure

```
beeapi/
├── puzzles/                  # Root directory for puzzle content
│   ├── theme1/               # Theme directory
│   │   ├── puzzle1.alghive   # Compressed puzzle file
│   │   ├── puzzle1/          # Extracted puzzle directory (created at runtime)
│   │   ├── puzzle2.alghive
│   │   ├── puzzle2/
│   ├── theme2/
│   │   ├── puzzle3.alghive
│   │   ├── puzzle3/
```

## API Authentication

BeeAPI Go implements a secure API key authentication system for protected endpoints:

- An API key is automatically generated on first startup
- The key is stored in a `.api-key` file in the root directory
- Protected endpoints (POST, PUT, DELETE) require this key in the Authorization header
- The key can be retrieved via the authenticated `/apikey` endpoint

```bash
# Example authenticated request
curl -X POST "http://localhost:5000/theme" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d "name=new-theme"
```

## API Documentation

The API documentation is accessible through Swagger UI at `/swagger/index.html` when the server is running. This provides:

- Interactive API testing interface
- Complete endpoint documentation
- Request/response schema details
- Authentication requirements

## Configuration

BeeAPI can be configured using environment variables:

- `SERVER_NAME`: The name of the server (default: "Local")
- `SERVER_DESCRIPTION`: A description of the server (default: "Local Dev Server")
- `PORT`: The port to run the server on (default: 5000)
- `PYTHON_PATH`: Path to Python interpreter for puzzle execution (default: "python")

## License

This project is licensed under the MIT License - see the LICENSE file for details.
