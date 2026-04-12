# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

openCMP is an open-source multi-cloud management platform built with Go + Gin + Gorm for the backend and Vue 3 + TypeScript + Element Plus for the frontend. The platform enables unified management and quick access to resources across multiple cloud providers through a unified interface and adapter pattern.

The project implements a modular architecture with core IAM (Identity and Access Management) module, cloud resource management, message center, and multi-cloud management capabilities. The platform supports tenant isolation through domains, resource grouping through projects, and fine-grained access control through roles and permissions.

## Architecture

The application follows a layered monolith architecture:

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
- **internal/migration/**: Database migration scripts

### Frontend Structure

- **src/views/**: Vue components for different modules (IAM, Cloud Accounts, Compute, Network, etc.)
- **src/api/**: API client functions
- **src/router.ts**: Vue Router configuration
- **src/layout/**: Main layout component with environment switching
- **src/utils/request.ts**: HTTP request utilities
- **src/utils/projectContext.ts**: Project context management for switching between management console and project-specific views
- **src/types/**: TypeScript type definitions

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

## Multi-Agent Development Framework

The project utilizes a multi-agent development framework with the following roles:

### Agent Roles & Responsibilities

- **Product Design Agent**: Responsible for module design, outputs executable specifications to `docs/specs/`
- **Task Distribution Agent**: Reads design docs, breaks down into backend/frontend executable tasks
- **Backend Development Agent**: Implements Go backend (handler/service/model/provider)
- **Frontend Development Agent**: Implements Vue 3 frontend pages and API calls
- **Test Agent**: Validates backend interfaces and frontend functionality

### Document-Driven Pipeline

```
Product Design Agent -> Specification Documents (docs/specs/)
    ↓
Task Distribution Agent -> Task Documents (docs/tasks/)
    ↓ (Parallel execution)
Backend Agent         Frontend Agent
  └ model/migration     └ views/api
  └ service             └ store/router
  └ handler/router      └ components
    ↓ (Backend completes first for API integration)
Test Agent Validation
```

### Module Development Priority

1. **IAM Core Modules**: Domain → Project → Group → User → Role → Permission → Authentication Sources
2. **Platform Functions**: Message Center (Inbox → Subscriptions → Channels → Receivers → Robots)
3. **Multi-Cloud Enhancement**: Cloud Account Enhancement → Sync Policies → Resource Sync Rules → Scheduled Sync Tasks

## Key Modules

### IAM Module
- **Domain Management**: Tenant isolation top-level unit, domain administrator, domain configuration, inter-domain switching
- **Project Management**: Domain resource grouping, project members, project quotas, project resource views
- **User Management**: Complete user lifecycle management, password policy, MFA support
- **Group Management**: User bulk management, group-level permission assignment
- **Role Management**: RBAC model-based roles, system roles and custom roles
- **Permission Management**: Fine-grained permission control, resource-level permissions
- **Policy Management**: Flexible policy engine supporting conditional policies and policy statements
- **Authentication Sources**: Support for multiple authentication sources (LDAP, local, SQL)

### Cloud Resource Management
- **Compute Resources**: VMs, images, key pairs, instance templates
- **Network Resources**: VPC, subnets, security groups, EIPs, load balancers
- **Storage Resources**: Disks, snapshots, object storage
- **Database Services**: RDS, Redis, MongoDB
- **Middleware Services**: Message queues, data analysis middleware

### Message Center
- **Internal Messages**: System notifications and message push
- **Message Subscriptions**: User-selected event types
- **Notification Channels**: Email, enterprise WeChat, DingTalk, Webhook
- **Recipient Management**: Notification targets (users/user groups/roles)
- **Robot Management**: Enterprise WeChat robots, DingTalk robots configuration

### Multi-Cloud Management
- **Cloud Account Management**: Unified management of multiple cloud accounts
- **Sync Policies**: Configure sync scope by resource type
- **Resource Sync Rules**: Automatic mapping to projects via tags
- **Scheduled Sync Tasks**: Cron-based scheduled sync

### Cloud Vendor Adapters
- **Unified Resource Classification**: Host, Network, Storage, Database, Middleware
- **Vendor Support**: Alibaba Cloud, Tencent Cloud, AWS, Azure
- **Adapter Pattern**: Standardized interfaces through ICompute, INetwork, IStorage, IDatabase

## Database Models

Key models include:
- CloudAccount: Cloud provider account configuration
- Domain: Tenant isolation units with domain configuration
- Project: Resource grouping within domains, with resource allocation tracking
- User: User accounts with authentication details
- Group: User groups for bulk management
- Role: Permission roles with role assignments
- Permission: Fine-grained permissions in format `<module>:<resource>:<action>`
- Policy: Policy engine configurations with policy statements
- Message: Internal messaging system with notification tracking
- NotificationChannel: Notification channels with configuration
- Robot: Chatbot integrations (WeChat, DingTalk, Feishu, etc.)
- Subscription: Event subscription rules for notifications
- Receiver: Notification recipient targets (users, groups, roles)

## Environment Switching Feature

A key feature of the platform is the ability to switch between different environments:
- **Management Console**: Shows all available resources and functions for global management
- **Project Mode**: Shows only project-related resources based on selected project context
- **Implementation**: Uses localStorage to store project context and dynamic sidebar menu rendering
- **Project-Specific Views**: Dedicated routes and components for project-related resources (alerts, messages, robots)

## Permission System

The permission system uses a hierarchical structure:
- Format: `<module>:<resource>:<action>` (e.g., `compute:vm:list`, `iam:user:create`)
- Domain-level permission inheritance
- Project-specific resource access controls
- Role-based access control (RBAC)

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
- Multi-factor authentication (MFA) support
- API request signature verification