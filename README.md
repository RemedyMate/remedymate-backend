# 🏥 RemedyMate Backend

A comprehensive health guidance and symptom triage API built with Go, featuring AI-powered medical advice, conversation management, and admin tools for healthcare professionals.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

## 🎯 Overview

RemedyMate is an intelligent health guidance platform that provides:

- **AI-powered symptom triage** using Google Gemini LLM
- **Multi-language support** (English & Amharic)
- **Conversation-based health assessments**
- **Admin tools** for managing red flags and feedback
- **Secure authentication** with JWT tokens
- **MongoDB-based data persistence**

## ✨ Features

### 🔍 Core Health Features

- **Symptom Triage**: AI-powered classification (GREEN/YELLOW/RED)
- **Topic Mapping**: Intelligent symptom-to-topic mapping
- **Guidance Cards**: Personalized health recommendations
- **Conversation Flow**: Interactive health assessments
- **Multi-language Support**: English and Amharic

### 👥 User Management

- **User Registration & Authentication**
- **JWT-based security**
- **Role-based access control** (Admin/SuperAdmin)
- **User status tracking**
- **OAuth integration** (Google, Facebook)

### 🛠️ Admin Features

- **Red Flag Management**: Create, update, delete medical red flags
- **Feedback Management**: Review and manage user feedback
- **Analytics Dashboard**: Health insights and statistics
- **Content Management**: Manage health guidance content

### 🔒 Security & Privacy

- **Stateless architecture**
- **Client-side session management**
- **Data encryption**
- **Input validation**
- **Rate limiting**

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │    │    Use Cases    │    │  Repositories   │
│                 │    │                 │    │                 │
│ • Auth          │◄──►│ • User          │◄──►│ • User          │
│ • RemedyMate    │    │ • RemedyMate    │    │ • Conversation │
│ • Conversation  │    │ • Conversation  │    │ • RedFlag       │
│ • Admin         │    │ • Admin         │    │ • Feedback      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │   Services      │    │   Database      │
│                 │    │                 │    │                 │
│ • Auth          │    │ • LLM (Gemini)  │    │ • MongoDB       │
│ • Validation    │    │ • Content       │    │ • Collections  │
│ • CORS          │    │ • Triage        │    │ • Indexes       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Tech Stack

### Backend

- **Go 1.24.5** - Core language
- **Gin Framework** - HTTP router and middleware
- **MongoDB** - Primary database
- **JWT** - Authentication tokens
- **Bcrypt** - Password hashing

### AI & LLM

- **Google Gemini** - Large Language Model
- **Custom prompts** - Medical triage and guidance
- **Multi-language processing**

### Infrastructure

- **Docker** - Containerization
- **Environment variables** - Configuration
- **Structured logging**

## 🚀 Installation

### Prerequisites

- Go 1.24.5 or higher
- MongoDB 4.4 or higher
- Google Gemini API key

### 1. Clone Repository

```bash
git clone https://github.com/your-org/remedymate-backend.git
cd remedymate-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Setup

```bash
cp env.example .env
```

### 4. Configure Environment Variables

Edit `.env` file with your configuration:

```env
# Server Configuration
PORT=8080

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
DB_NAME=remedymate

# JWT Configuration
JWT_SECRET_KEY=your_super_secret_jwt_key_minimum_32_characters_here
JWT_EXPIRY_HOURS=24

# Gemini LLM Configuration
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-1.5-flash

# OAuth2 Configuration (Optional)
GOOGLE_CLIENT_ID=your_google_client_id_here
GOOGLE_CLIENT_SECRET=your_google_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/google/callback

FACEBOOK_CLIENT_ID=your_facebook_client_id_here
FACEBOOK_CLIENT_SECRET=your_facebook_client_secret_here
FACEBOOK_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/facebook/callback
```

### 5. Start MongoDB

```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or using local installation
mongod
```

### 6. Run Application

```bash
go run delivery/main.go
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### 🔐 Authentication Endpoints

#### Register User

```http
POST /auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password123"
}
```

#### Login User

```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure_password123"
}
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "admin"
  }
}
```

### 🏥 Health Guidance Endpoints

#### Get Remedy (Main Endpoint)

```http
POST /remedy
Content-Type: application/json

{
  "text": "I have a severe headache and feel dizzy",
  "language": "en"
}
```

**Response:**

```json
{
  "session_id": "sess_12345",
  "triage": {
    "level": "YELLOW",
    "red_flags": [],
    "message": "Your symptoms require attention but are not immediately dangerous."
  },
  "guidance_card": {
    "title": "Headache Management",
    "content": "Here's what you can do...",
    "self_care": ["Rest in a dark room", "Apply cold compress"],
    "seek_care_if": ["Headache worsens", "Vision changes"]
  }
}
```

#### Start/Continue Conversation

```http
POST /conversation
Content-Type: application/json

{
  "session_id": "sess_12345",
  "user_input": "The headache started this morning",
  "language": "en"
}
```

### 💬 Feedback Endpoints

#### Submit Public Feedback

```http
POST /feedbacks
Content-Type: application/json

{
  "sessionId": "sess_12345",
  "topicKey": "headache",
  "language": "en",
  "rating": 4,
  "message": "The guidance was helpful"
}
```

### 🛠️ Admin Endpoints (Authentication Required)

#### Red Flag Management

**List Red Flags:**

```http
GET /admin/redflags
Authorization: Bearer <jwt_token>
```

**Create Red Flag:**

```http
POST /admin/redflags
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "keywords": ["chest pain", "heart attack"],
  "language": "en",
  "level": "RED",
  "description": "Severe chest pain indicating cardiac emergency"
}
```

**Update Red Flag:**

```http
PUT /admin/redflags/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "keywords": ["chest pain", "heart attack", "myocardial infarction"],
  "level": "RED",
  "description": "Updated description"
}
```

**Delete Red Flag:**

```http
DELETE /admin/redflags/{id}
Authorization: Bearer <jwt_token>
```

#### Feedback Management

**List Feedbacks:**

```http
GET /admin/feedbacks?limit=20&offset=0&language=en
Authorization: Bearer <jwt_token>
```

**Get Specific Feedback:**

```http
GET /admin/feedbacks/{id}
Authorization: Bearer <jwt_token>
```

**Delete Feedback:**

```http
DELETE /admin/feedbacks/{id}
Authorization: Bearer <jwt_token>
```

## 🗄️ Database Schema

### Collections

#### Users Collection

```json
{
  "_id": "ObjectId",
  "username": "string",
  "email": "string",
  "passwordHash": "string",
  "personalInfo": {
    "firstName": "string",
    "lastName": "string",
    "age": "number",
    "gender": "string",
    "profilePictureUrl": "string"
  },
  "role": "admin|superadmin",
  "createdBy": "string",
  "updatedBy": "string",
  "createdAt": "Date",
  "updatedAt": "Date",
  "lastLogin": "Date"
}
```

#### Red Flags Collection

```json
{
  "_id": "ObjectId",
  "keywords": ["string"],
  "language": "string",
  "level": "RED|YELLOW",
  "description": "string",
  "isDeleted": "boolean",
  "createdAt": "Date",
  "updatedAt": "Date",
  "deletedAt": "Date",
  "createdBy": "string",
  "updatedBy": "string",
  "deletedBy": "string"
}
```

#### Feedback Collection

```json
{
  "_id": "ObjectId",
  "sessionId": "string",
  "topicKey": "string",
  "language": "string",
  "rating": "number",
  "message": "string",
  "isDeleted": "boolean",
  "createdAt": "Date",
  "deletedAt": "Date"
}
```

### Indexes

- **Users**: `username` (unique), `email` (unique)
- **User Status**: `userId` (unique)
- **Conversations**: `user_id`, `status`, `created_at`
- **Red Flags**: `language`, `level`, `isDeleted`
- **Feedback**: `sessionId`, `language`, `rating`

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run specific test package
go test ./usecase/...

# Run with coverage
go test -cover ./...
```

### Test Data

The application includes test data in the `data/` directory:

- `approved_block.json` - Approved content blocks
- `red_flag_rules.json` - Red flag detection rules
- `yellow_flag_rules.json` - Yellow flag detection rules

### Postman Collection

Import the provided Postman collection for comprehensive API testing:

1. **Authentication Tests**

   - Register user
   - Login user
   - Refresh token
   - Logout

2. **Health Guidance Tests**

   - Submit symptoms
   - Get triage results
   - Start conversations
   - Continue conversations

3. **Admin Tests**
   - Manage red flags
   - Manage feedback
   - View analytics

## 🚀 Deployment

### Docker Deployment

#### Dockerfile

```dockerfile
FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o remedymate-backend delivery/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/remedymate-backend .
COPY --from=builder /app/data ./data

EXPOSE 8080
CMD ["./remedymate-backend"]
```

#### Docker Compose

```yaml
version: "3.8"
services:
  remedymate-backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MONGO_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    depends_on:
      - mongodb

  mongodb:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

volumes:
  mongodb_data:
```

## 🤝 Contributing

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit changes: `git commit -m 'Add amazing feature'`
6. Push to branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Write tests for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

### Getting Help

- **Documentation**: Check this README and API docs
- **Issues**: Create GitHub issues for bugs
- **Discussions**: Use GitHub discussions for questions

### Contact

- **Email**: support@remedymate.com
- **GitHub**: [@remedymate-backend](https://github.com/your-org/remedymate-backend)

## 🔄 Changelog

### Version 1.0.0

- Initial release
- Core health guidance features
- Admin management tools
- Multi-language support
- JWT authentication
- MongoDB integration

---

**Built with ❤️ for better healthcare accessibility**
