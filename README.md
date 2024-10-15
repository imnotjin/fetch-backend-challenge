# Fetch Backend Internship Challenge

## Project Overview

This project implements a REST API for managing a points-based reward system. The API allows tracking of user points across multiple payers, handling point transactions, and managing point balances.

---

## Prerequisites

Before running this project, ensure you have the following installed on your system:

- [Go](https://golang.org/dl/) (version 1.23.2 or higher)
- [Docker](https://docs.docker.com/get-docker/) (version 27.2.0 or higher)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.29.2 or higher)
- [Make](https://www.gnu.org/software/make/) (optional, for using Makefile commands)

---

## Setup Instructions

You can set up and run this project either using the provided Makefile or manually.

### Option 1: Using Makefile

1. Clone the Repository:

   ```bash
   git clone https://github.com/imnotjin/fetch-backend-challenge.git
   cd fetch-backend-challenge
   ```

2. Set up environment files:

   ```bash
   make setup
   ```

3. Start the application:

   ```bash
   make up
   ```

4. To stop the application:

   ```bash
   make down
   ```

### Option 2: Manual Setup

1. Clone the Repository:

   ```bash
   git clone https://github.com/imnotjin/fetch-backend-challenge.git
   cd fetch-backend-challenge
   ```

2. Environment Variables:
   Create a `.env` file with the following content:

   ```env
   DB_USER=postgres
   DB_PASSWORD=your_secure_password
   DB_NAME=points_db
   DB_HOST=db
   DB_PORT=5432
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=points_db
   ```

   Create a `.env.test` file with the following content:

   ```env
   DB_USER=postgres
   DB_PASSWORD=your_secure_password
   DB_NAME=test_db
   DB_HOST=localhost
   DB_PORT=5433
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=test_db
   ```

3. Docker Setup:

   - Ensure Docker and Docker Compose are installed.
   - Start the PostgreSQL databases (production & test) and API server:
     ```bash
     docker-compose up --build
     ```
   - To stop the containers:
     ```bash
     docker-compose down
     ```

## Running the Application

Once the containers are up, the API will be running at `http://localhost:8000`.

### API Documentation

You can access the API documentation via Swagger at:

- **URL**: http://localhost:8000/swagger/index.html

## Testing

This project includes an integration test to ensure the functionality of the API endpoints and their interaction with the database.

### Running Tests

You can run the test using either of these methods:

1. Using Makefile:

   ```bash
   make test
   ```

2. Manually:
   ```bash
   go test ./handlers -v
   ```

### Integration Test

Location: `./handlers`

This test focuses on verifying the correct behavior of the API endpoints as described in the project requirements. It covers:

1. Adding points through the `/add` endpoint
2. Spending points through the `/spend` endpoint
3. Checking point balance through the `/balance` endpoint

The test follows this sequence:

1. Adds points using five different transactions
2. Spends a specified amount of points
3. Verifies the final balance

**Expected outcome**:

After adding the following transactions:

- { "payer": "DANNON", "points": 300, "timestamp": "2022-10-31T10:00:00Z" }
- { "payer": "UNILEVER", "points": 200, "timestamp": "2022-10-31T11:00:00Z" }
- { "payer": "DANNON", "points": -200, "timestamp": "2022-10-31T15:00:00Z" }
- { "payer": "MILLER COORS", "points": 10000, "timestamp": "2022-11-01T14:00:00Z" }
- { "payer": "DANNON", "points": 1000, "timestamp": "2022-11-02T14:00:00Z" }

And spending 5000 points, the final balance should be:

```json
{
  "DANNON": 1000,
  "UNILEVER": 0,
  "MILLER COORS": 5300
}
```

This test ensures that the API correctly handles point transactions and maintains accurate balances per payer.
