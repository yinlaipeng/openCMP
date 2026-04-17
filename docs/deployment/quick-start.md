# openCMP 快速部署指南

## 系统要求

- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ RAM
- 10GB+ 磁盘空间

## 快速部署

### 1. 准备配置

```bash
# 克隆项目
git clone https://github.com/opencmp/opencmp.git
cd opencmp

# 复制环境配置模板
cp .env.example .env
```

### 2. 配置环境变量

编辑 `.env` 文件，修改以下关键配置：

```bash
# 数据库配置（生产环境必须修改）
MYSQL_ROOT_PASSWORD=your-secure-password
MYSQL_PASSWORD=your-secure-password

# JWT密钥（生产环境必须修改）
JWT_SECRET=your-secure-jwt-key-at-least-32-characters

# 服务模式
SERVER_MODE=release
```

### 3. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 验证部署

访问以下地址验证服务正常：

- 前端: http://localhost:80
- API健康检查: http://localhost:8080/health
- 登录测试: 使用 admin/admin123 登录

## 服务说明

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| MySQL | opencmp-mysql | 3306 | 数据库服务 |
| Redis | opencmp-redis | 6379 | 缓存服务 |
| Backend | opencmp-backend | 8080 | API服务 |
| Frontend | opencmp-frontend | 80 | 前端界面 |
| Nginx | opencmp-nginx | 8081/443 | 反向代理 |

## 常用操作

### 停止服务
```bash
docker-compose down
```

### 停止并清理数据
```bash
docker-compose down -v
```

### 重启单个服务
```bash
docker-compose restart backend
```

### 查看日志
```bash
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 进入容器调试
```bash
docker exec -it opencmp-backend sh
docker exec -it opencmp-mysql mysql -u root -p
```

## 数据库管理

### 备份数据库
```bash
docker exec opencmp-mysql mysqldump -u root -p${MYSQL_ROOT_PASSWORD} opencmp > backup.sql
```

### 恢复数据库
```bash
docker exec -i opencmp-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} opencmp < backup.sql
```

## SSL配置

1. 准备SSL证书文件：
```bash
mkdir ssl
cp your-cert.pem ssl/cert.pem
cp your-key.pem ssl/key.pem
```

2. 编辑 `nginx.conf`，启用HTTPS配置块

3. 重启nginx服务：
```bash
docker-compose restart nginx
```

## 性能优化

### 数据库连接池
在 `.env` 中配置：
```bash
DB_MAX_IDLE_CONNS=20
DB_MAX_OPEN_CONNS=100
```

### Redis缓存
确保Redis服务正常运行：
```bash
docker exec opencmp-redis redis-cli ping
# 应返回 PONG
```

## 故障排除

### 服务无法启动
1. 检查端口占用：`lsof -i :8080`
2. 检查容器日志：`docker-compose logs`
3. 检查磁盘空间：`df -h`

### 数据库连接失败
1. 确认MySQL服务正常：`docker-compose ps mysql`
2. 检查密码配置：`docker exec opencmp-mysql mysql -u root -p`
3. 检查网络：`docker network ls`

### 前端无法访问后端
1. 检查backend服务：`curl http://localhost:8080/health`
2. 检查nginx配置：`docker exec opencmp-nginx nginx -t`
3. 检查容器网络：`docker network inspect opencmp_default`

## 生产环境建议

1. **安全配置**:
   - 修改所有默认密码
   - 启用HTTPS
   - 配置防火墙规则

2. **性能配置**:
   - 增加数据库连接池大小
   - 配置Redis持久化
   - 启用Nginx gzip压缩

3. **监控配置**:
   - 配置日志收集
   - 设置健康检查告警
   - 监控资源使用率

4. **备份策略**:
   - 定期数据库备份
   - 配置文件备份
   - 测试恢复流程