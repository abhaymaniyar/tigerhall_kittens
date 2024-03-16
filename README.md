# Tigerhall Kittens Service

Tigerhall Kittens is a service designed to track sightings of tigers in the wild, allowing users to report and list
sightings via a RESTful API.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them:

- Go (version 1.18 or later)
- Docker (optional for containerization)
- PostgreSQL (or any preferred database, adjust accordingly)

### Installing

A step by step series of examples that tell you how to get a development environment running:

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/tigerhall-kittens.git

### Running Tests

To run the unit tests for the project, use the following command:

```sh
go test ./...
```

For integration tests, ensure your test database is correctly set up and use:

```sh
go test -tags=integration ./...
```

## Deployment

To deploy this project on a live system, consider containerizing it with Docker:

1. Build the Docker image:

```sh
docker build -t tigerhall-kittens .
```

2. Run the Docker container:

```sh
docker run -p 8080:8080 tigerhall-kittens
```

Adjust the port mappings and any other configurations as necessary.

## API Contracts

### User Endpoints

- POST **/users**: Create a new user
    - Request body:
        ```json
        {
          "username": "johndoe",
          "password": "securepassword",
          "email": "johndoe@example.com"
        }

- POST **/login**: Authenticate a user

  Request body:
    ``` json
    {
      "username": "johndoe",
      "password": "securepassword"
    }

#### Tiger and Sighting Endpoints

- POST /tigers: Add a new tiger
- GET /tigers: List all tigers with pagination
- POST /sightings: Report a new sighting
- GET /tigers/{id}/sightings: List all sightings for a tiger

Please refer to the `api_docs.md` for detailed request and response structures.

### Built With

- Gin - The web framework used
- GORM - ORM library for Go
- Viper - Configuration management
- jwt-go - JWT implementation for Go