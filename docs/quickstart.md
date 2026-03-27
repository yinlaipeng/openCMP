# 快速入门指南

## 1. 环境准备

### 安装 Go

```bash
# macOS
brew install go@1.21

# Linux
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```

### 安装 MySQL

```bash
# macOS
brew install mysql@5.7

# Linux (Ubuntu)
sudo apt-get install mysql-server-5.7
```

## 2. 初始化数据库

```bash
mysql -u root -p < scripts/init.sql
```

## 3. 配置应用

复制配置文件并修改:

```bash
cp configs/config.yaml configs/config.local.yaml
```

修改数据库连接字符串:

```yaml
database:
  dsn: "root:your-password@tcp(127.0.0.1:3306)/opencmp?charset=utf8mb4&parseTime=True&loc=Local"
```

## 4. 运行服务

```bash
# 安装依赖
go mod download

# 运行
go run cmd/server/main.go -config configs/config.local.yaml
```

服务启动后，访问 `http://localhost:8080/health` 检查健康状态。

## 5. API 测试

### 5.1 添加云账户

```bash
curl -X POST http://localhost:8080/api/v1/cloud-accounts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的阿里云账号",
    "provider_type": "alibaba",
    "credentials": {
      "access_key_id": "LTAI5t...",
      "access_key_secret": "xxxxx"
    }
  }'
```

### 5.2 列出云账户

```bash
curl -X GET http://localhost:8080/api/v1/cloud-accounts
```

### 5.3 验证云账户

```bash
curl -X POST http://localhost:8080/api/v1/cloud-accounts/1/verify
```

### 5.4 创建 VPC

```bash
curl -X POST http://localhost:8080/api/v1/network/vpcs \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 1,
    "name": "my-vpc",
    "cidr": "10.0.0.0/16"
  }'
```

### 5.5 列出 VPC

```bash
curl -X GET "http://localhost:8080/api/v1/network/vpcs?account_id=1"
```

### 5.6 创建虚拟机

```bash
curl -X POST http://localhost:8080/api/v1/compute/vms \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 1,
    "name": "test-vm",
    "instance_type": "ecs.t5-lc1m1.small",
    "image_id": "m-bp1xxxxxxxx",
    "vpc_id": "vpc-xxx",
    "subnet_id": "vsw-xxx",
    "disk_size": 40
  }'
```

## 6. 常见问题

### Q: 如何获取阿里云 AccessKey?

A: 登录阿里云控制台 -> 用户头像 -> AccessKey 管理

### Q: 数据库连接失败？

A: 检查 MySQL 服务是否运行，用户名密码是否正确

### Q: 如何查看日志？

A: 日志输出到控制台，配置文件设置 `log.level` 调整日志级别

## 7. 下一步

- 阅读完整 [API 文档](README.md#api-端点)
- 阅读 [架构设计](README.md#架构设计)
- 尝试接入其他云厂商
