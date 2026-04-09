# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

openCMP is an open-source multi-cloud management platform built with Go + Gin + Gorm for the backend and Vue 3 + TypeScript + Element Plus for the frontend. The platform enables unified management and quick access to resources across multiple cloud providers through a unified interface and adapter pattern.

## Architecture

The application follows a layered architecture:

```
┌──────────────────────────────────────────────────────────────────────┐
│                         API Layer (Gin)                               │
│                    RESTful API / HTTP Handlers                        │
├──────────────────────────────────────────────────────────────────────┤
│                        Service Layer                                  │
│           Business logic layer (resource mgmt/account mgmt/task orchestration/permission control)              │
├──────────────────────────────────────────────────────────────────────┤
│                     Cloud Provider Layer                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ Alibaba  │  │ Tencent  │  │   AWS    │  │  Azure   │             │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │             │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘             │
├──────────────────────────────────────────────────────────────────────┤
│                  Cloud Interface Layer (standardized interfaces)      │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐     │
│  │   Compute  │  │  Network   │  │  Storage   │  │  Database  │     │
│  │ ICompute   │  │ INetwork   │  │ IStorage   │  │ IDatabase  │     │
│  └────────────┘  └────────────┘  └────────────┘  └────────────┘     │
├──────────────────────────────────────────────────────────────────────┤
│                      Data Layer (Gorm)                                │
│              MySQL/PostgreSQL / Cloud account config / Resource metadata storage            │
└──────────────────────────────────────────────────────────────────────┘
```

### Backend Structure

- **cmd/server/main.go**: Main application entry point
- **internal/handler/**: HTTP handlers for API endpoints
- **internal/service/**: Business logic services
- **internal/model/**: Database models
- **internal/middleware/**: HTTP middleware (logging, auth, recovery)
- **pkg/**: Shared utilities and cloud provider interfaces
- **configs/config.yaml**: Application configuration

### Frontend Structure

- **src/views/**: Vue components for different modules (IAM, Cloud Accounts, Compute, Network, etc.)
- **src/api/**: API client functions
- **src/router.ts**: Vue Router configuration
- **src/layout/**: Main layout component
- **src/utils/request.ts**: HTTP request utilities

## Development Commands

### Backend Commands

```bash
# Navigate to backend directory
cd backend

# Install dependencies
make deps

# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Run the application
make run

# Clean build artifacts
make clean

# Docker build
make docker-build

# Start MySQL with docker-compose
make docker-up

# Stop MySQL
make docker-down
```

### Frontend Commands

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Run tests
npm run test
```

## Key Modules

### IAM Module
- **Domain Management**: Tenant isolation and multi-domain support
- **Project Management**: Resource grouping within domains
- **User Management**: Full user lifecycle management
- **Group Management**: Bulk user management
- **Role Management**: RBAC-based roles with system and custom roles
- **Permission Management**: Fine-grained permission control
- **Policy Management**: Flexible policy engine
- **Authentication Sources**: Support for LDAP, local, SQL authentication

### Cloud Resource Management
- **Compute Resources**: VMs, images, key pairs, instance templates
- **Network Resources**: VPC, subnets, security groups, EIPs, load balancers
- **Storage Resources**: Disks, snapshots, object storage
- **Database Services**: RDS, Redis, MongoDB

### Message Center
- **Internal Messages**: System notifications and message push
- **Message Subscriptions**: User-selected event types
- **Notification Channels**: Email, enterprise WeChat, DingTalk, Webhook
- **Recipient Management**: Notification targets (users/user groups/roles)
- **Robot Management**: Enterprise WeChat robots, DingTalk robots

### Multi-Cloud Management
- **Cloud Account Management**: Unified management of multiple cloud accounts
- **Sync Policies**: Configure sync scope by resource type
- **Resource Sync Rules**: Automatic mapping to projects via tags
- **Scheduled Sync Tasks**: Cron-based scheduled sync

## Database Models

Key models include:
- CloudAccount: Cloud provider account configuration
- Domain: Tenant isolation units
- Project: Resource grouping within domains
- User: User accounts
- Group: User groups
- Role: Permission roles
- Permission: Fine-grained permissions
- Policy: Policy engine configurations
- Message: Internal messaging system
- NotificationChannel: Notification channels
- Robot: Chatbot integrations

## Testing

To run all backend tests:
```bash
cd backend
go test -v ./...
```

To run specific tests:
```bash
# Run specific test file
go test -v internal/service/user_test.go

# Run specific test function
go test -v -run TestUserCreate
```

## Configuration

The application uses `configs/config.yaml` for configuration with sections for:
- server: Port and mode settings
- database: Connection details
- auth: JWT secret and token expiration
- log: Logging level and format
- providers: Cloud provider enablement
- message_center: Messaging configuration
- multicloud_sync: Multi-cloud synchronization settings

## Security Features

- JWT Token Authentication
- Role-Based Access Control (RBAC)
- Fine-grained permission control
- Password encryption with bcrypt
- SQL injection protection
- XSS and CSRF protection