# GoBackend

A structured, high-performance Go-based backend powered by **Gin** and a custom **Fastapify** wrapper for automated OpenAPI (Swagger) documentation.

## Features

- **Gin HTTP Framework:** Fast and lightweight routing.
- **Fastapify Wrapper:** Automatically generates Swagger UI (`/docs`) and OpenAPI JSON (`/openapi.json`) straight from strongly typed handlers and structs.
- **MongoDB Integration:** Built-in connection management.
- **Structured Modules:** Code is organized by domain (Movies, Users, etc.).
- **Validation:** Uses `go-playground/validator` (via Gin binding) for robust schema checking.

## Getting Started

### Prerequisites
* Go 1.18+
* MongoDB URI configured in your `.env` file

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sharathcx/GoBackend.git
   cd GoBackend
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Environment Setup:**
   Create a `.env` file in the root directory and add your connection strings:
   ```env
   MONGO_URI=mongodb://localhost:27017
   PORT=8000
   ```

4. **Run the API:**
   ```bash
   go run main.go
   ```

### Viewing API Documentation
Once the server is running, navigate to:
- **Swagger UI:** `http://localhost:8000/docs`
- **OpenAPI JSON:** `http://localhost:8000/openapi.json`

## Project Architecture

```
GoBackend/
├── fastapify/           # Custom OpenAPI tracking & UI logic
├── modules/             # Business domains (movies, users)
│   ├── movie/           # Movie routes, schemas, DB logic
│   └── user/            # User routes, schemas, DB logic
├── database/            # Global MongoDB connector
├── globals/             # Environment variables and app state
└── main.go              # Entry point & API bootstrapping
```
