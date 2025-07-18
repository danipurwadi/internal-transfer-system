# Internal Transfer System

This project is a simple internal transfer system that allows users to create accounts, view their balance, and transfer funds to other accounts. The project is built with Go and uses a PostgreSQL database to store the data.

Note: The design of the project is designed using philosophies from [Ardan Lab's Service Starter Kit](https://github.com/ardanlabs/service). Some of the packages used are also based on the starter kit, although the actual implementation is adapted based on what I find useful for this project.

## Assumptions

1. Maximum value of account balance is below `99,999,999,999,999`
2. Balance and transaction amounts are accurate up to 5 decimal points. Values that are more accurate than that can be rounded up/down to 5 d.p.

## Getting Started

To get started with this project, you will need to have the following installed on your local machine:

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [make](https://www.gnu.org/software/make/)
  - `make` is optional but is handy to run the various commands easily
  - You can view the commands in `makefile` directly and execute them at the root of the project

Once you have the prerequisites installed, you can follow these steps to get the project up and running:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/danipurwadi/internal-transfer-system.git
    ```

2.  **Install the dependencies:**

    ```bash
    make dev-gotooling
    ```

3.  **Start the application:**

    ```bash
    make start
    ```

This will start the application and the PostgreSQL database in Docker containers. The application will be available at `http://localhost:8080`.

**Other helpful commands**

a. Run test for project

```bash
make test
```

b. Run a golang script to load test the application

```bash
make load-test
```

## Project Structure

The project is divided into the following main components:

- **`app`:** This directory contains the application layer, which is responsible for handling HTTP requests, decoding and encoding JSON, and calling the business layer.
- **`business`:** This directory contains the business layer, which is responsible for implementing the business logic of the application.
- **`foundation`:** This directory contains the foundation layer, which provides common functionality that is used by the other layers, such as logging, error handling, and database access.
- **`zarf`:** This directory contains the configuration for the project. The name zarf is means a sleeve that protects your hand from hot containers.

## API Endpoints

The application exposes the following REST API endpoints:

### 1. Health Check

- **GET `/health`**
  - Description: Checks the health status of the application.
  - Response:
    - `200 OK`
    ```json
    {
      "status": true
    }
    ```

### 2. Account Management

- **POST `/accounts`**

  - Description: Creates a new account with a specified initial balance.
  - Request Body:
    ```json
    {
      "account_id": 123,
      "initial_balance": "100.00"
    }
    ```
  - Response:
    - `201 Created` (on success, no response body)
    - `400 Bad Request` (e.g., invalid JSON, missing fields)
    - `409 Conflict` (if `account_id` already exists)

- **GET `/accounts/{account_id}`**
  - Description: Retrieves the current balance for a given account.
  - Path Parameters:
    - `account_id` (integer): The ID of the account to query.
  - Response:
    - `200 OK`
    ```json
    {
      "account_id": "123",
      "balance": "100.00"
    }
    ```
    - `400 Bad Request` (e.g., invalid `account_id` format)
    - `404 Not Found` (if `account_id` does not exist)

### 3. Transaction Management

- **POST `/transactions`**
  - Description: Creates a new transaction to transfer funds between two accounts.
  - Request Body:
    ```json
    {
      "source_account_id": 123,
      "destination_account_id": 456,
      "amount": "50.00"
    }
    ```
  - Response:
    - `201 Created` (on success, no response body)
    - `400 Bad Request` (e.g., invalid JSON, missing fields, `source_account_id` equals `destination_account_id`, negative `amount`)
    - `404 Not Found` (if `source_account_id` or `destination_account_id` does not exist)
    - `422 Unprocessable Entity` (if `source_account_id` has insufficient funds)

## Available Commands

The following `make` commands are available:

- `make dev-gotooling`: Install the Go tooling dependencies.
- `make tidy`: Tidy the `go.mod` file.
- `make start`: Start the application and the database.
- `make stop`: Stop the application and the database.
- `make exit`: Stop the application and the database and remove the volumes.
- `make test`: Run the tests.
- `make test-race`: Run the tests with the race detector.
- `make stats`: Open the statistics page in the browser.
- `make load-test`: Run a golang script to load test the application.
