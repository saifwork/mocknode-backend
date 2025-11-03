# Mock Service - MockNode Backend

## Overview
A modern backend built in Golang that allows developers to easily create, manage, and test APIs through projects and collections. This is the engine behind MockNode, perfect for developers, learners, and testers who want to simulate API endpoints without building a custom backend.

## Project Status
- **Type**: Pure Backend API Service (no frontend)
- **Language**: Go 1.24.5
- **Framework**: Gin Web Framework
- **Database**: MongoDB (cloud-hosted)
- **Port**: 8080 (localhost)
- **Last Updated**: November 3, 2025

## Architecture
This is a RESTful API service with the following components:
- **AuthService**: Manages users, login, JWTs, email verification
- **ProjectService**: User projects CRUD
- **CollectionService**: Schema & validation management
- **RecordService**: Dynamic data CRUD
- **ConfigHandler**: Field types & prebuilt mock data endpoints

## Key Features
- JWT-based authentication with email verification
- User projects (free users can create up to 2 projects)
- Dynamic collections with field validation (up to 3 per project for free users)
- CRUD operations for records
- Prebuilt mock data endpoints (users, products, posts, todos, comments, carts)
- Health check endpoint

## Project Structure
```
cmd/server/          - Main application entry point
internal/
  api/               - HTTP handlers and routing
  core/              - Core functionality (config, MongoDB, Redis clients)
  dtos/              - Data transfer objects
  middlewares/       - HTTP middlewares (auth, CORS, logging, recovery)
  models/            - Data models
  services/          - Business logic
  utils/             - Utility functions
data/                - Prebuilt mock data (JSON files)
```

## Running the Project
The backend runs automatically via the `backend` workflow on port 8080.

### Available Endpoints
- **Health**: `GET /health`
- **Auth**: `/auth/*` (signup, login, verify-email, forgot-password, reset-password)
- **Projects**: `/api/projects/*` (CRUD operations)
- **Collections**: `/api/projects/:pid/collections/*` and `/api/collections/*`
- **Records**: `/api/collections/:cid/records/*`
- **Config**: `/config/types` and `/config/preset/*` (users, products, posts, etc.)

## Configuration
The application uses environment variables with sensible defaults:
- APP_PORT: 8080
- MONGO_URI: Cloud MongoDB instance (hardcoded default)
- JWT_SECRET: supersecret (change in production)
- Redis configuration (optional, currently disabled in code)

All configuration is managed through `internal/core/config/config.go`.

## Dependencies
- **gin-gonic/gin**: Web framework
- **mongo-driver**: MongoDB client
- **go-redis**: Redis client (optional)
- **golang-jwt**: JWT token handling
- **godotenv**: Environment variable loading
- **golang.org/x/crypto**: Password hashing

## Recent Changes
- **Nov 3, 2025**: Initial Replit setup
  - Installed Go 1.23
  - Configured backend workflow on port 8080
  - Updated .gitignore for Go projects
  - Verified MongoDB connection and all endpoints working
  - Health check confirmed operational

## Notes
- This is a backend-only service with no frontend
- MongoDB is hosted externally (cloud instance)
- Redis client code is commented out in main.go
- The service uses hardcoded MongoDB credentials in the config defaults
- Email functionality uses Gmail SMTP (credentials in config defaults)
