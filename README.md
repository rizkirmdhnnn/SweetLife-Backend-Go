# SweetLife Backend

SweetLife is a smart app designed to help people maintain a healthy diet, especially for diabetics or those who care about their health. The app uses an advanced technology called Machine Learning to scan food, provide nutritional information, and recommend food choices that are suitable for each user's needs.

## Features
- nanti deh lagi males ngisi

## Prerequisites

Before running this program, make sure you have installed:
- Go version 1.20 or higher
- PostgreSQL Database
- Google Cloud Platform account for storage bucket
- Mailgun account for email service

## Installation

```bash
# Clone the repository
git clone https://github.com/rizkirmdhnnn/SweetLife-Backend-Go

# Navigate to project directory
cd SweetLife-Backend-Go

# Install dependencies
go mod download

# Copy environment file
cp .example.env .env
```

## Usage

```bash
# Run the program
go run main.go

# Build the program
go build -o sweetlife-backend

# Run the built application
./sweetlife-backend
```

## Project Structure

```
.
├── config/             # Application configuration files
├── dto/                # Data Transfer Objects
├── email/              # Email templates and handlers
├── errors/             # Custom error definitions
├── handlers/           # HTTP request handlers (Controllers in Laravel)
├── helpers/            # Helper functions and utilities
├── middleware/         # HTTP middleware
├── models/             # Data models and structs
├── repositories/       # Data access layer
├── routers/            # Route definitions
├── services/          # Business logic layer
├── templates/         # HTML/template files
├── .example.env       # Environment variables template
├── .gitignore        # Git ignore file
├── Dockerfile        # Docker configuration
├── go.mod            # Go modules file
├── go.sum            # Go modules checksums
└── main.go           # Application entry point
```

## Configuration

Configure your application by copying `.example.env` to `.env` and setting the appropriate values:

```env
# Application Settings
APP_HOST="http://localhost"
APP_PORT="3000"
APP_ENV="development"
APP_KEY="verysecretkey"
APP_DEBUG="True"

# Database Configuration
DB_HOST="localhost"
DB_USER=""
DB_PASSWORD=""
DB_NAME=""
DB_PORT=""

# JWT Configuration
JWTSIGNKEY="secretkey"

# Mail Configuration (Mailgun)
MAILGUNKEY=""
MAILGUNDOMAIN=""
MAILFROM=""

# Storage Configuration (Google Cloud)
STORAGE_BUCKET=""
STORAGE_FOLDER=""
PROJECT_ID=""
```

### Environment Variables Description

#### Application Settings
- `APP_HOST`: Base URL for the application
- `APP_PORT`: Port where the application will run
- `APP_ENV`: Current environment (development/production)
- `APP_KEY`: Application encryption key
- `APP_DEBUG`: Enable/disable debug mode

#### Database Settings
- `DB_HOST`: PostgreSQL server host
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_PORT`: Database port

#### JWT Configuration
- `JWTSIGNKEY`: Secret key for JWT token signing

#### Mail Settings (Mailgun)
- `MAILGUNKEY`: Mailgun API key
- `MAILGUNDOMAIN`: Mailgun domain
- `MAILFROM`: Default sender email address

#### Storage Settings (Google Cloud)
- `STORAGE_BUCKET`: Google Cloud Storage bucket name
- `STORAGE_FOLDER`: Storage folder path
- `PROJECT_ID`: Google Cloud project identifier

## API Documentation
https://sweetlife.apidog.io

