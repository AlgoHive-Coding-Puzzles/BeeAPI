<p align="center">
  <img width="150px" src="https://raw.githubusercontent.com/AlgoHive-Coding-Puzzles/Ressources/refs/heads/main/images/beeapi-logo.png" title="Algohive">
</p>

<h1 align="center">BeeAPI Go</h1>

## Self-Hostable API for AlgoHive

BeeAPI Go is a Go implementation of the BeeAPI service, designed to load puzzles from `.alghive` files, organized into themes, and serve them to the AlgoHive platform independently. It comes with a Swagger UI to test the API endpoints and a simple web interface to manage the puzzles.

![Swagger](img/swagger.png)

## AlgoHive

AlgoHive is a web, self-hostable platform that allows developers to create puzzles for other developers to solve. Each puzzle contains two parts to solve, allowing developers to test their skills in a variety of ways. The puzzles are created using a proprietary file format that is compiled into a single file for distribution.

## Installation

### Local

To use BeeAPI Go, you need to have Go 1.16 or higher installed on your system.

```bash
# Clone the repository
git clone https://github.com/AlgoHive-Coding-Puzzles/BeeAPI-Go.git
cd BeeAPI-Go

# Build the application
go build -o beeapi

# Run the application
./beeapi
```

Feed the `puzzles` directory with themes folders containing `.alghive` files or use the Web interface to manage the puzzles for the BeeAPI instance.

```
beeapi/
├── puzzles/
│   ├── theme1/
│   │   ├── puzzle1.alghive
│   │   ├── puzzle2.alghive
│   ├── theme2/
│   │   ├── puzzle3.alghive
│   │   ├── puzzle4.alghive
```

### Docker

You can also run the API using Docker. To build the Docker image:

```bash
docker build -t beeapi-go .
```

Then, run the Docker container:

```bash
docker run -d -p 8080:8080 --name beeapi-go -v $(pwd)/puzzles:/app/puzzles -v $(pwd)/data:/app/data beeapi-go
```

## API Authentication

BeeAPI Go uses API key authentication for protected endpoints (POST, PUT, DELETE). When the server starts for the first time, it generates a unique API key and saves it in a `.api-key` file in the root directory.

To authenticate your requests to protected endpoints:

1. Look for the API key in the `.api-key` file or in the server logs when it starts
2. Include the API key in your requests' Authorization header:

```bash
curl -X POST "http://localhost:5000/theme" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d "name=new-theme"
```

If you're using Docker, you can access the API key by:

```bash
docker exec beeapi-go cat /app/.api-key
```

Make sure to keep this key secure as it provides administrative access to your BeeAPI instance.

## API Documentation

The API documentation is available at `/swagger/index.html` when the server is running.

## Environment Variables

The following environment variables can be used to configure the server:

- `SERVER_NAME`: The name of the server (default: "Local")
- `SERVER_DESCRIPTION`: A description of the server (default: "Local Dev Server")
- `USER_*`: Create users automatically, for example `USER_ADMIN=adminpassword`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
