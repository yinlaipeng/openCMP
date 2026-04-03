# openCMP 开发环境配置与部署指南

## 1. 项目概述

openCMP 是一个基于 Go + Gin + Gorm 的单体架构多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

### 1.1 技术栈
- **后端**: Go, Gin, Gorm, PostgreSQL/MySQL
- **前端**: Vue 3, TypeScript, Element Plus
- **数据库**: PostgreSQL 12+, MySQL 8+
- **缓存**: Redis (可选)
- **消息队列**: RabbitMQ/Kafka (可选)
- **容器化**: Docker, Docker Compose
- **部署**: Kubernetes (可选)

## 2. 开发环境搭建

### 2.1 系统要求
- **操作系统**: Linux/macOS/Windows 10+
- **Go**: 1.19+
- **Node.js**: 18+
- **数据库**: PostgreSQL 12+ 或 MySQL 8+
- **Docker**: 20+ (推荐)

### 2.2 后端环境配置

#### 2.2.1 安装 Go
```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang
```

#### 2.2.2 克装依赖
```bash
cd backend
go mod tidy
```

#### 2.2.3 数据库配置
```bash
# 安装 PostgreSQL (macOS)
brew install postgresql

# 安装 PostgreSQL (Ubuntu)
sudo apt install postgresql postgresql-contrib

# 启动 PostgreSQL
brew services start postgresql  # macOS
sudo systemctl start postgresql   # Linux
```

#### 2.2.4 创建数据库
```sql
CREATE DATABASE opencmp;
CREATE USER opencmp WITH PASSWORD 'opencmp123';
GRANT ALL PRIVILEGES ON DATABASE opencmp TO opencmp;
```

### 2.3 前端环境配置

#### 2.3.1 安装 Node.js
```bash
# 使用 nvm (推荐)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18

# 或直接安装
# macOS
brew install node

# Ubuntu
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### 2.3.2 安装前端依赖
```bash
cd frontend
npm install
```

## 3. 项目配置

### 3.1 后端配置

#### 3.1.1 创建配置文件
```bash
# backend/configs/config.yaml
server:
  port: 8080
  mode: debug

database:
  driver: postgres
  dsn: "host=localhost user=opencmp password=opencmp123 dbname=opencmp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
  max_idle_conns: 10
  max_open_conns: 100

auth:
  jwt_secret: "your-super-secret-jwt-key-change-this-in-production"
  token_expire_hours: 24

log:
  level: info
  format: json
```

#### 3.1.2 环境变量配置
```bash
# backend/.env
DATABASE_URL=postgres://opencmp:opencmp123@localhost:5432/opencmp?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
LOG_LEVEL=info
```

### 3.2 前端配置

#### 3.2.1 创建环境配置
```bash
# frontend/.env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=openCMP
VITE_PORT=3000
VITE_PROXY_TARGET=http://localhost:8080
```

#### 3.2.2 配置代理 (vite.config.ts)
```typescript
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: process.env.VITE_PROXY_TARGET || 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  }
})
```

## 4. 开发模式启动

### 4.1 后端开发模式
```bash
# 启动后端服务
cd backend
go run cmd/server/main.go

# 或使用 air 热重载
go install github.com/cosmtrek/air@latest
air -c .air.conf
```

### 4.2 前端开发模式
```bash
# 启动前端开发服务器
cd frontend
npm run dev
```

### 4.3 使用 Docker Compose 开发
```bash
# docker-compose.dev.yml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    container_name: opencmp_postgres_dev
    environment:
      POSTGRES_DB: opencmp
      POSTGRES_USER: opencmp
      POSTGRES_PASSWORD: opencmp123
    ports:
      - "5432:5432"
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U opencmp"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    container_name: opencmp_backend_dev
    depends_on:
      postgres:
        condition: service_healthy
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://opencmp:opencmp123@postgres:5432/opencmp?sslmode=disable
      - JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
    volumes:
      - ./backend:/app
    command: ["air", "-c", ".air.conf"]

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    container_name: opencmp_frontend_dev
    depends_on:
      - backend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
    environment:
      - VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 5. 生产环境部署

### 5.1 Docker 部署

#### 5.1.1 后端 Dockerfile
```dockerfile
# backend/Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./main"]
```

#### 5.1.2 前端 Dockerfile
```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### 5.1.3 生产环境 Docker Compose
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    container_name: opencmp_postgres
    environment:
      POSTGRES_DB: opencmp
      POSTGRES_USER: opencmp
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U opencmp"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: opencmp_redis
    restart: unless-stopped
    volumes:
      - redis_data:/data

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: opencmp_backend
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/opencmp?sslmode=disable
      - REDIS_ADDR=redis:6379
      - JWT_SECRET=${JWT_SECRET}
      - PORT=8080
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: opencmp_frontend
    depends_on:
      - backend
    ports:
      - "80:80"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    container_name: opencmp_nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

### 5.2 Kubernetes 部署

#### 5.2.1 部署文件
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: opencmp-backend
  labels:
    app: opencmp-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: opencmp-backend
  template:
    metadata:
      labels:
        app: opencmp-backend
    spec:
      containers:
      - name: backend
        image: opencmp/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: opencmp-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: opencmp-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: opencmp-backend-service
spec:
  selector:
    app: opencmp-backend
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

## 6. CI/CD 配置

### 6.1 GitHub Actions
```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.19'
        
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        
    - name: Install Go dependencies
      run: |
        cd backend
        go mod download
        
    - name: Run Go tests
      run: |
        cd backend
        go test -v ./...
        
    - name: Install Node dependencies
      run: |
        cd frontend
        npm ci
        
    - name: Run Node tests
      run: |
        cd frontend
        npm run test
        
    - name: Build frontend
      run: |
        cd frontend
        npm run build

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - name: Checkout
      uses: actions/checkout@v3
      
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
      
    - name: Login to DockerHub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}
        
    - name: Build and push backend
      uses: docker/build-push-action@v4
      with:
        context: ./backend
        file: ./backend/Dockerfile
        push: true
        tags: ${{ secrets.DOCKER_USERNAME }}/opencmp-backend:latest
        
    - name: Build and push frontend
      uses: docker/build-push-action@v4
      with:
        context: ./frontend
        file: ./frontend/Dockerfile
        push: true
        tags: ${{ secrets.DOCKER_USERNAME }}/opencmp-frontend:latest
```

## 7. 监控与日志

### 7.1 日志配置
```go
// backend/internal/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level string, format string) (*zap.Logger, error) {
    var atom = zap.NewAtomicLevel()
    
    switch level {
    case "debug":
        atom.SetLevel(zap.DebugLevel)
    case "info":
        atom.SetLevel(zap.InfoLevel)
    case "warn":
        atom.SetLevel(zap.WarnLevel)
    case "error":
        atom.SetLevel(zap.ErrorLevel)
    default:
        atom.SetLevel(zap.InfoLevel)
    }

    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "timestamp"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.StacktraceKey = ""

    var encoder zapcore.Encoder
    if format == "console" {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    }

    core := zapcore.NewCore(encoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), atom)

    return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}
```

### 7.2 健康检查
```go
// backend/internal/handler/health.go
func HealthCheck(c *gin.Context) {
    // 检查数据库连接
    db, err := c.MustGet("db").(*gorm.DB)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "error":  "database not initialized",
        })
        return
    }

    sqlDB, err := db.DB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "error":  "database connection error",
        })
        return
    }

    if err := sqlDB.Ping(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "error":  "database ping failed",
        })
        return
    }

    // 检查其他依赖项...

    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "checks": map[string]string{
            "database": "ok",
        },
    })
}
```

## 8. 性能优化

### 8.1 数据库优化
```go
// 连接池配置
func initDB() *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.New(
            log.New(os.Stdout, "\r\n", log.LstdFlags),
            logger.Config{
                SlowThreshold: time.Second,
                LogLevel:      logger.Info,
                Colorful:      true,
            },
        ),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get sql.DB:", err)
    }

    // 连接池配置
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return db
}
```

### 8.2 缓存策略
```go
// Redis 缓存示例
type CacheService struct {
    client *redis.Client
}

func (c *CacheService) GetUser(id uint) (*model.User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 尝试从缓存获取
    val, err := c.client.Get(context.Background(), key).Result()
    if err == nil {
        var user model.User
        if err := json.Unmarshal([]byte(val), &user); err == nil {
            return &user, nil
        }
    }

    // 缓存未命中，从数据库获取
    user, err := c.getUserFromDB(id)
    if err != nil {
        return nil, err
    }

    // 存储到缓存
    data, _ := json.Marshal(user)
    c.client.Set(context.Background(), key, data, time.Hour)

    return user, nil
}
```

## 9. 安全配置

### 9.1 CORS 配置
```go
// backend/cmd/server/main.go
func setupCORS(r *gin.Engine) {
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

### 9.2 安全头配置
```go
// 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
        
        c.Next()
    }
}
```

## 10. 部署脚本

### 10.1 部署脚本
```bash
#!/bin/bash
# deploy.sh

set -e

echo "Starting deployment..."

# 构建后端
echo "Building backend..."
cd backend
go build -o opencmp-backend cmd/server/main.go

# 构建前端
echo "Building frontend..."
cd ../frontend
npm run build

# 复制构建产物到部署目录
echo "Copying build artifacts..."
mkdir -p ../deploy/dist
cp -r dist/* ../deploy/dist/

# 启动服务
echo "Starting services..."
cd ../deploy
docker-compose -f docker-compose.prod.yml up -d

echo "Deployment completed!"
```

### 10.2 回滚脚本
```bash
#!/bin/bash
# rollback.sh

echo "Rolling back to previous version..."

# 停止当前服务
docker-compose -f docker-compose.prod.yml down

# 启动备份版本
docker-compose -f docker-compose.backup.yml up -d

echo "Rollback completed!"
```

## 11. 故障排除

### 11.1 常见问题
1. **数据库连接失败**: 检查 DSN 配置和网络连接
2. **端口冲突**: 检查端口占用情况
3. **权限不足**: 检查数据库用户权限
4. **SSL 错误**: 配置正确的 SSL 模式

### 11.2 调试命令
```bash
# 查看容器日志
docker logs opencmp_backend

# 进入容器调试
docker exec -it opencmp_backend sh

# 检查数据库连接
docker exec -it opencmp_postgres psql -U opencmp -d opencmp

# 检查服务状态
docker-compose ps
```

## 12. 维护任务

### 12.1 数据库备份
```bash
# 备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/opencmp"
mkdir -p $BACKUP_DIR

pg_dump -h localhost -U opencmp -d opencmp > $BACKUP_DIR/opencmp_$DATE.sql
gzip $BACKUP_DIR/opencmp_$DATE.sql

# 保留最近7天的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

### 12.2 日志轮转
```bash
# logrotate 配置
/var/log/opencmp/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

这个部署指南涵盖了从开发环境搭建到生产环境部署的完整流程，包括 CI/CD、监控、安全配置和维护任务。