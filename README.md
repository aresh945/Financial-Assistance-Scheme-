# Financial Assistance Scheme Management System

A backend system for managing financial assistance schemes and applications

## Overview
This system allows administrators to:
- Manage financial assistance schemes
- Process applicant registrations
- Determine scheme eligibility
- Track applications and their statuses

## Installation

### Database Setup
1. Create a PostgreSQL database:
```sql
CREATE DATABASE tutorial1;
```

2. Run the schema file:
```bash
psql -d tutorial1 -f database/GTDB.sql
```

### Project Setup
1. Clone the repository:
```bash
git clone https://github.com/yourusername/Financial-Assistance-Scheme
cd Financial-Assistance-Scheme
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure database connection in cmd/api/main.go:
```go
dbConfig := &database.Config{
    Host:     "localhost",
    Port:            // Your PostgreSQL port
    User:         // Your PostgreSQL username
    Password: 
    DBName:   "tutorial1",
    SSLMode:  "disable",
}
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Documentation

### Applicants

#### Get All Applicants
```http
GET /api/applicants
```

#### Create Applicant
```http
POST /api/applicants
```
Example request:
```json
{
    "id": "01913b7a-4493-74b2-93f8-e684c4ca935c",
    "name": "James Smith",
    "employment_status": "unemployed",
    "marital_status": "single",
    "sex": "male",
    "date_of_birth": "1990-07-01T00:00:00Z",
    "household": []
}
```

### Schemes
#### Get All Schemes
```http
GET /api/schemes
```

#### Get Eligible Schemes
```http
GET /api/schemes/eligible?applicant={id}
```

### Applications
#### Get All Applications
```http
GET /api/applications
```

#### Create Application
```http
POST /api/applications
```
Example request:
```json
{
    "application_id": "01913b90-5d23-7abc-9def-123456789abc",
    "applicant_id": "01913b7a-4493-74b2-93f8-e684c4ca935c",
    "scheme_id": "01913b89-9a43-7163-8757-01cc254783f3",
    "status": "pending"
}
```

## Project Structure
```
financial_assistance/
├── docs/
│   └── financial-assistance-api.postman_collection.json
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── models/              # Data structures
│   ├── repository/          # Database interactions
│   ├── service/            # Business logic
│   └── handler/            # HTTP handlers
├── pkg/
│   └── database/           # Database utilities
├── GTDB.sql 
│ 
└── go.mod
```
