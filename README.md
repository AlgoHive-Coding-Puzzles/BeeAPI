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
docker run -d -p 8080:8080 --name beeapi-go -v $(pwd)/puzzles:/app/puzzles beeapi-go
```

## API Documentation

The API documentation is available at `/swagger/index.html` when the server is running.

## Environment Variables

The following environment variables can be used to configure the server:

- `SERVER_NAME`: The name of the server (default: "Local")
- `SERVER_DESCRIPTION`: A description of the server (default: "Local Dev Server")
- `USER_*`: Create users automatically, for example `USER_ADMIN=adminpassword`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
