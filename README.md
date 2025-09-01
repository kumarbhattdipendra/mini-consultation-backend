# Backend API

## Overview

This backend system provides various endpoints to manage users, guides, and bookings. It is built with the [Gin Gonic web framework](https://github.com/gin-gonic/gin) in Go and connected to a PostgreSQL database using GORM.



## API Endpoints

### Public Endpoints

- **Health Check**
  - **GET** `/health`: Checks the health status of the API.

- **Authentication**
  - **POST** `/auth/register`: Register a new user.
  - **POST** `/auth/login`: Login an existing user.

- **Guides**
  - **GET** `/guides`: Retrieve a list of guides.

### Protected Endpoints

These endpoints require authentication.

- **Bookings**
  - **GET** `/bookings`: Retrieve user's bookings.
  - **POST** `/bookings`: Create a new booking.

## Middleware

- **AuthMiddleware**: This middleware is used to protect routes that require authentication. It requires a valid JWT token.

## Database

This backend uses PostgreSQL as its database system. The database connection is facilitated through the GORM ORM.

## Environment Variables

- `JWT_SECRET`: This is the secret key used for signing JWT tokens.

## Database Setup

Before running the project, create the necessary enum type in PostgreSQL:

```sql
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');
```

## Getting Started

1. Clone the repository.
2. Set up the environment variables in a `.env` file.
3. Execute the SQL command to create the enum type in the database.
4. Run the application with `go run main.go`.

