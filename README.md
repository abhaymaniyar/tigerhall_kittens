# Tigerhall Kittens Service

Tigerhall Kittens is a service designed to track sightings of tigers in the wild, allowing users to report and list
sightings via a RESTful API.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them:

- Go (version 1.18 or later)
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

## Deployment

1. Create a dev environment file from the sample file

 ```sh
 $cp env.sample development.env
 ```

2. Modify the values in `development.env` to match your needs.

```sh
./run-server.sh [env-file]
```

This command will load environment variables of the specified env file and run the server on the port specified in the
env variables file.

### Detailed API Spec with sample responses

https://documenter.getpostman.com/view/3144528/2sA2xnxVD2

### Postman Collection

https://api.postman.com/collections/3144528-69009c19-e07f-4d08-b5bc-3981e8ceeb36?access_key=PMAT-01HS4YCA9CSY9J19XF7PF8G6B2

### Built With

- Gin - The web framework used
- GORM - ORM library for Go
- Goose - DB migration library
- Viper - Configuration management
- jwt-go - JWT implementation for Go

