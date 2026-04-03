# openCMP IAM 模块后端开发指南

## 1. 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── handler/             # HTTP Handler (Gin)
│   │   ├── auth.go          # 认证相关
│   │   ├── domain.go        # 域管理
│   │   ├── project.go       # 项目管理
│   │   ├── user.go          # 用户管理
│   │   ├── group.go         # 用户组管理
│   │   ├── role.go          # 角色管理
│   │   ├── permission.go    # 权限管理
│   │   └── policy.go        # 策略管理
│   ├── service/             # 业务逻辑层
│   │   ├── user_service.go  # 用户服务
│   │   ├── role_service.go  # 角色服务
│   │   ├── permission_service.go # 权限服务
│   │   ├── policy_service.go # 策略服务
│   │   └── ...
│   ├── model/               # 数据模型 (Gorm)
│   │   ├── user.go          # 用户模型
│   │   ├── role.go          # 角色模型
│   │   ├── permission.go    # 权限模型
│   │   ├── policy.go        # 策略模型
│   │   └── ...
│   ├── middleware/          # Gin 中间件
│   │   ├── auth.go          # 认证中间件
│   │   └── permission.go    # 权限中间件
│   └── utils/               # 工具函数
│       ├── jwt.go           # JWT 工具
│       └── password.go      # 密码工具
├── pkg/
│   └── cloudprovider/       # 云适配器层
├── configs/                 # 配置文件
└── scripts/                 # 脚本
```

## 2. 数据模型设计

### 2.1 用户模型 (User)
```go
type User struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;index" json:"name"`
    DisplayName string    `gorm:"type:varchar(100);not null;index" json:"display_name"`
    Email       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
    Phone       string    `gorm:"type:varchar(20)" json:"phone"`
    Password    string    `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
    Enabled     bool      `gorm:"default:true" json:"enabled"`
    MFEnabled   bool      `gorm:"default:false" json:"mfa_enabled"`
    DomainID    uint      `gorm:"not null;index" json:"domain_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // 关联关系
    Domain Domain `gorm:"foreignKey:DomainID" json:"domain"`
    Group  Group  `gorm:"foreignKey:GroupID" json:"group"`
}
```

### 2.2 角色模型 (Role)
```go
type Role struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;index" json:"name"`
    DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name"`
    Description string    `gorm:"type:text" json:"description"`
    Type        string    `gorm:"type:varchar(20);default:'custom'" json:"type"` // system, custom
    Enabled     bool      `gorm:"default:true" json:"enabled"`
    IsPublic    bool      `gorm:"default:false" json:"is_public"` // 是否公开
    DomainID    uint      `gorm:"not null;index" json:"domain_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // 关联关系
    Domain Domain `gorm:"foreignKey:DomainID" json:"domain"`
}
```

### 2.3 权限模型 (Permission)
```go
type Permission struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;index" json:"name"` // 权限标识
    DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name"` // 显示名称
    Resource    string    `gorm:"type:varchar(50);not null;index" json:"resource"` // 资源类型
    Action      string    `gorm:"type:varchar(50);not null;index" json:"action"` // 操作类型
    Type        string    `gorm:"type:varchar(20);default:'custom'" json:"type"` // system, custom
    Description string    `gorm:"type:text" json:"description"`
    DomainID    uint      `gorm:"not null;index" json:"domain_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // 关联关系
    Domain Domain `gorm:"foreignKey:DomainID" json:"domain"`
}
```

### 2.4 策略模型 (Policy)
```go
type Policy struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;index" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    Scope       string    `gorm:"type:varchar(20);default:'system'" json:"scope"` // system, domain, project
    Type        string    `gorm:"type:varchar(20);default:'custom'" json:"type"` // system, custom
    Enabled     bool      `gorm:"default:true" json:"enabled"`
    IsSystem    bool      `gorm:"default:false" json:"is_system"` // 是否系统策略
    IsPublic    bool      `gorm:"default:false" json:"is_public"` // 是否公开
    DomainID    uint      `gorm:"index" json:"domain_id"` // 可选，仅当 scope 为 domain 时使用
    Policy      string    `gorm:"type:jsonb" json:"policy"` // 策略内容，JSON格式
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // 关联关系
    Domain Domain `gorm:"foreignKey:DomainID" json:"domain"`
}
```

## 3. 服务层开发

### 3.1 用户服务 (UserService)
```go
type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User, password string) error {
    // 验证域是否存在
    var domain model.Domain
    if err := s.DB.First(&domain, user.DomainID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("domain with ID %d not found", user.DomainID)
        }
        return err
    }

    // 检查邮箱是否已存在
    var existingUserByEmail model.User
    if err := s.DB.Where("email = ?", user.Email).First(&existingUserByEmail).Error; err == nil {
        return fmt.Errorf("user with email '%s' already exists", user.Email)
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    // 检查用户名是否已存在
    var existingUserByUsername model.User
    if err := s.DB.Where("name = ?", user.Name).First(&existingUserByUsername).Error; err == nil {
        return fmt.Errorf("user with name '%s' already exists", user.Name)
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    // 哈希密码
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %v", err)
    }
    user.Password = hashedPassword

    // 设置默认值
    if user.Enabled {
        user.Enabled = true
    }

    return s.DB.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
    var user model.User
    err := s.DB.Preload("Domain").Preload("Group").First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with ID %d not found", id)
        }
        return nil, err
    }
    // 不返回密码
    user.Password = ""
    return &user, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *UserService) GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error) {
    var permissions []model.Permission

    err := s.DB.WithContext(ctx).
        Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
        Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
        Where("user_roles.user_id = ?", userID).
        Find(&permissions).Error

    if err != nil {
        return nil, err
    }

    return permissions, nil
}
```

### 3.2 权限检查服务
```go
type PermissionService struct {
    DB *gorm.DB
}

// CheckUserPermission 检查用户是否有指定权限
func (s *PermissionService) CheckUserPermission(userID uint, resource, action string) (bool, error) {
    // 检查用户是否为超级管理员
    isAdmin, err := s.isUserAdmin(userID)
    if err != nil {
        return false, err
    }
    if isAdmin {
        return true, nil
    }

    // 检查用户直接权限
    hasDirectPermission, err := s.hasDirectPermission(userID, resource, action)
    if err != nil {
        return false, err
    }
    if hasDirectPermission {
        return true, nil
    }

    // 检查用户组权限
    hasGroupPermission, err := s.hasGroupPermission(userID, resource, action)
    if err != nil {
        return false, err
    }
    if hasGroupPermission {
        return true, nil
    }

    // 检查用户角色权限
    hasRolePermission, err := s.hasRolePermission(userID, resource, action)
    if err != nil {
        return false, err
    }
    if hasRolePermission {
        return true, nil
    }

    // 检查策略权限
    hasPolicyPermission, err := s.hasPolicyPermission(userID, resource, action)
    if err != nil {
        return false, err
    }
    if hasPolicyPermission {
        return true, nil
    }

    return false, nil
}

// isUserAdmin 检查用户是否为管理员
func (s *PermissionService) isUserAdmin(userID uint) (bool, error) {
    // 获取管理员角色
    var adminRole model.Role
    if err := s.DB.Where("name = ? AND type = ?", "admin", "system").First(&adminRole).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }

    // 检查用户是否拥有管理员角色
    var userRole model.UserRole
    if err := s.DB.Where("user_id = ? AND role_id = ?", userID, adminRole.ID).First(&userRole).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }

    return true, nil
}
```

## 4. Handler 层开发

### 4.1 用户 Handler
```go
type UserHandler struct {
    Service *service.UserService
    Logger  *zap.Logger
}

func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
    return &UserHandler{
        Service: service,
        Logger:  logger,
    }
}

// CreateUser 创建用户
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "User details (without password)"
// @Param password body string true "User password"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req createUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.Logger.Error("Failed to bind user request", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := &model.User{
        Name:        req.Name,
        DisplayName: req.DisplayName,
        Email:       req.Email,
        Phone:       req.Phone,
        Enabled:     req.Enabled,
        DomainID:    req.DomainID,
    }

    if err := h.Service.CreateUser(user, req.Password); err != nil {
        h.Logger.Error("Failed to create user", zap.Error(err))
        if err.Error()[:7] == "domain" || err.Error()[:4] == "user" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    // 不返回密码
    user.Password = ""
    c.JSON(http.StatusCreated, user)
}

// GetUserPermissions 获取用户权限
// @Summary Get user permissions
// @Description Get all permissions assigned to a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []model.Permission
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/permissions [get]
func (h *UserHandler) GetUserPermissions(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    ctx := c.Request.Context()
    permissions, err := h.Service.GetUserPermissions(ctx, uint(id))
    if err != nil {
        h.Logger.Error("Failed to get user permissions", zap.Uint("user_id", uint(id)), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
        return
    }

    c.JSON(http.StatusOK, permissions)
}

// Request body structs
type createUserRequest struct {
    Name        string `json:"name" binding:"required"`
    DisplayName string `json:"display_name" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    Phone       string `json:"phone"`
    Password    string `json:"password" binding:"required,min=8"`
    Enabled     bool   `json:"enabled"`
    DomainID    uint   `json:"domain_id" binding:"required"`
}
```

## 5. 中间件开发

### 5.1 认证中间件
```go
// AuthMiddleware 认证中间件
func AuthMiddleware(logger *zap.Logger, jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // 解析Bearer token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        // 解析JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // 验证签名算法
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            logger.Info("Invalid token", zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 获取声明
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            logger.Error("Invalid token claims")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 检查用户ID是否存在
        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            logger.Error("User ID not found in token")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        userID := uint(userIDFloat)

        // 将用户信息存储到上下文中
        c.Set("user_id", userID)
        c.Set("user_name", claims["user_name"])
        c.Set("user_email", claims["user_email"])

        c.Next()
    }
}
```

### 5.2 权限中间件
```go
// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(logger *zap.Logger, db *gorm.DB, requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }

        // 检查用户权限
        hasPermission, err := checkUserPermission(db, userID.(uint), requiredPermission)
        if err != nil {
            logger.Error("Error checking user permission", zap.Error(err))
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            c.Abort()
            return
        }

        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission '%s' required", requiredPermission)})
            c.Abort()
            return
        }

        c.Next()
    }
}

// checkUserPermission 检查用户是否具有指定权限
func checkUserPermission(db *gorm.DB, userID uint, permissionName string) (bool, error) {
    // 检查用户是否为超级管理员（拥有所有权限）
    var adminRole model.Role
    if err := db.Where("name = ? AND type = ?", "admin", "system").First(&adminRole).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 如果没有找到管理员角色，则继续检查普通权限
        } else {
            return false, err
        }
    } else {
        // 检查用户是否拥有管理员角色
        var userRole model.UserRole
        if err := db.Where("user_id = ? AND role_id = ?", userID, adminRole.ID).First(&userRole).Error; err != nil {
            if !errors.Is(err, gorm.ErrRecordNotFound) {
                return false, err
            }
        } else {
            // 用户是管理员，拥有所有权限
            return true, nil
        }
    }

    // 检查用户直接拥有的权限
    if hasDirectPermission(db, userID, permissionName) {
        return true, nil
    }

    // 检查用户所属组的权限
    if hasGroupPermission(db, userID, permissionName) {
        return true, nil
    }

    // 检查用户角色的权限
    if hasRolePermission(db, userID, permissionName) {
        return true, nil
    }

    return false, nil
}
```

## 6. 数据库迁移

### 6.1 自动迁移
```go
// 在 main.go 中
func main() {
    // ... 数据库连接代码 ...

    // 自动迁移表结构
    err = db.AutoMigrate(
        &model.User{},
        &model.Role{},
        &model.Permission{},
        &model.UserRole{},
        &model.RolePermission{},
        &model.Group{},
        &model.GroupRole{},
        &model.Policy{},
        &model.PolicyStatement{},
        &model.RolePolicy{},
        &model.Domain{},
        &model.Project{},
        &model.AuthSource{},
    )
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // ... 其他代码 ...
}
```

### 6.2 初始化默认数据
```go
// 初始化默认数据
func initDefaultData(db *gorm.DB) error {
    ctx := context.Background()

    // 创建默认域
    defaultDomain := &model.Domain{
        Name:        "Default",
        Description: "默认域",
        Enabled:     true,
    }

    var existingDomain model.Domain
    if err := db.Where("name = ?", "Default").First(&existingDomain).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            if err := db.Create(defaultDomain).Error; err != nil {
                return fmt.Errorf("failed to create default domain: %v", err)
            }
        } else {
            return fmt.Errorf("failed to check default domain: %v", err)
        }
    } else {
        defaultDomain = &existingDomain
    }

    // 创建管理员角色
    adminRole := &model.Role{
        Name:        "admin",
        DisplayName: "系统管理员",
        Description: "系统管理员角色，拥有所有权限",
        Type:        "system",
        Enabled:     true,
        IsPublic:    true,
    }

    var existingRole model.Role
    if err := db.Where("name = ?", "admin").First(&existingRole).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            if err := db.Create(adminRole).Error; err != nil {
                return fmt.Errorf("failed to create admin role: %v", err)
            }
        } else {
            return fmt.Errorf("failed to check admin role: %v", err)
        }
    } else {
        adminRole = &existingRole
    }

    // 创建默认管理员用户
    adminUser := &model.User{
        Name:        "admin",
        DisplayName: "管理员",
        Email:       "admin@example.com",
        Password:    "admin123", // 密码将在创建时被哈希
        DomainID:    defaultDomain.ID,
        Enabled:     true,
    }

    var existingUser model.User
    if err := db.Where("name = ?", "admin").First(&existingUser).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // 对密码进行哈希处理
            hashedPassword, err := utils.HashPassword(adminUser.Password)
            if err != nil {
                return fmt.Errorf("failed to hash password: %v", err)
            }
            adminUser.Password = hashedPassword

            if err := db.Create(adminUser).Error; err != nil {
                return fmt.Errorf("failed to create admin user: %v", err)
            }

            // 为管理员用户分配管理员角色
            userService := service.NewUserService(db)
            if err := userService.AssignUserRole(ctx, adminUser.ID, adminRole.ID, defaultDomain.ID); err != nil {
                return fmt.Errorf("failed to assign admin role to admin user: %v", err)
            }
        } else {
            return fmt.Errorf("failed to check admin user: %v", err)
        }
    }

    // 创建基本权限
    permissions := []model.Permission{
        {Name: "user:list", DisplayName: "用户列表", Description: "查看用户列表", Resource: "user", Action: "list", Type: "system"},
        {Name: "user:create", DisplayName: "创建用户", Description: "创建新用户", Resource: "user", Action: "create", Type: "system"},
        // ... 其他权限
    }

    for _, perm := range permissions {
        var existingPerm model.Permission
        if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                if err := db.Create(&perm).Error; err != nil {
                    return fmt.Errorf("failed to create permission %s: %v", perm.Name, err)
                }
            } else {
                return fmt.Errorf("failed to check permission %s: %v", perm.Name, err)
            }
        }
    }

    // 为管理员角色分配所有权限
    roleService := service.NewRoleService(db)
    for _, perm := range permissions {
        var existingPerm model.Permission
        if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
            return fmt.Errorf("failed to get permission %s: %v", perm.Name, err)
        }

        if err := roleService.AssignPermissionToRole(ctx, adminRole.ID, existingPerm.ID); err != nil {
            return fmt.Errorf("failed to assign permission %s to admin role: %v", perm.Name, err)
        }
    }

    return nil
}
```

## 7. 安全最佳实践

### 7.1 密码安全
```go
package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt哈希密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash 比较密码和哈希值
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 7.2 JWT 安全
```go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
    UserID      uint   `json:"user_id"`
    UserName    string `json:"user_name"`
    UserEmail   string `json:"user_email"`
    jwt.RegisteredClaims
}

// GenerateJWTToken 生成JWT token
func GenerateJWTToken(userID uint, userName, userEmail, secret string, expireHours int) (string, error) {
    expirationTime := time.Now().Add(time.Hour * time.Duration(expireHours))
    claims := &JWTClaims{
        UserID:    userID,
        UserName:  userName,
        UserEmail: userEmail,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

## 8. 测试

### 8.1 单元测试
```go
// service/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    db, tearDown := setupTestDB()
    defer tearDown()

    service := NewUserService(db)

    t.Run("Successful user creation", func(t *testing.T) {
        user := &model.User{
            Name:        "test-user",
            DisplayName: "Test User",
            Email:       "test@example.com",
            DomainID:    1,
        }

        err := service.CreateUser(user, "password123")
        assert.NoError(t, err)
        assert.NotZero(t, user.ID)
        assert.Equal(t, "test-user", user.Name)
    })

    t.Run("Duplicate email", func(t *testing.T) {
        user1 := &model.User{
            Name:        "user1",
            DisplayName: "User 1",
            Email:       "duplicate@example.com",
            DomainID:    1,
        }
        user2 := &model.User{
            Name:        "user2",
            DisplayName: "User 2",
            Email:       "duplicate@example.com", // Same email
            DomainID:    1,
        }

        err1 := service.CreateUser(user1, "password123")
        assert.NoError(t, err1)

        err2 := service.CreateUser(user2, "password123")
        assert.Error(t, err2)
        assert.Contains(t, err2.Error(), "already exists")
    })
}
```

### 8.2 集成测试
```go
// handler/user_handler_test.go
func TestUserHandler_CreateUser(t *testing.T) {
    db, tearDown := setupTestDB()
    defer tearDown()

    service := NewUserService(db)
    handler := NewUserHandler(service, testLogger)

    t.Run("Successful user creation", func(t *testing.T) {
        w := httptest.NewRecorder()
        ginCtx, _ := gin.CreateTestContext(w)

        userData := model.User{
            Name:        "test-create-user",
            DisplayName: "Test Create User",
            Email:       "create@example.com",
            DomainID:    1,
        }
        jsonData, _ := json.Marshal(createUserRequest{
            Name:        userData.Name,
            DisplayName: userData.DisplayName,
            Email:       userData.Email,
            Password:    "password123",
            DomainID:    userData.DomainID,
        })

        req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
        req.Header.Set("Content-Type", "application/json")
        ginCtx.Request = req

        handler.CreateUser(ginCtx)

        assert.Equal(t, http.StatusCreated, w.Code)

        var responseUser model.User
        err := json.Unmarshal(w.Body.Bytes(), &responseUser)
        assert.NoError(t, err)
        assert.Equal(t, "test-create-user", responseUser.Name)
        assert.Equal(t, "Test Create User", responseUser.DisplayName)
    })
}
```

## 9. 配置管理

### 9.1 配置文件
```yaml
# configs/config.yaml
server:
  port: 8080
  mode: debug  # debug, release, test

database:
  driver: mysql  # mysql, postgres, sqlite
  dsn: "root:password@tcp(localhost:3306)/opencmp?charset=utf8mb4&parseTime=True&loc=Local"
  max_idle_conns: 10
  max_open_conns: 100

auth:
  jwt_secret: "your-secret-key-change-this-in-production"  # JWT密钥
  token_expire_hours: 24  # Token过期时间（小时）

log:
  level: info    # debug, info, warn, error
  format: json   # json, console
```

### 9.2 配置加载
```go
// configs/config.go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Auth     AuthConfig     `yaml:"auth"`
    Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
    Driver       string `yaml:"driver"`
    DSN          string `yaml:"dsn"`
    MaxIdleConns int    `yaml:"max_idle_conns"`
    MaxOpenConns int    `yaml:"max_open_conns"`
}

type AuthConfig struct {
    JWTSecret        string `yaml:"jwt_secret"`
    TokenExpireHours int    `yaml:"token_expire_hours"`
}

type LogConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

## 10. 常见问题

### 10.1 如何添加新的权限？
1. 在数据库中创建权限记录
2. 在角色中分配该权限
3. 在API中使用权限中间件保护端点

### 10.2 如何实现基于角色的访问控制？
1. 用户-角色关联 (UserRole)
2. 角色-权限关联 (RolePermission)
3. 在中间件中检查用户角色和权限

### 10.3 如何处理多租户权限？
1. 使用域 (Domain) 隔离不同租户
2. 在权限检查时验证资源归属
3. 实现跨域权限管理功能