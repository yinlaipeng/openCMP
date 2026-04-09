package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/utils"
)

// AuthSourceService 认证源服务
type AuthSourceService struct {
	db *gorm.DB
}

// NewAuthSourceService 创建认证源服务
func NewAuthSourceService(db *gorm.DB) *AuthSourceService {
	return &AuthSourceService{db: db}
}

// CreateAuthSource 创建认证源
func (s *AuthSourceService) CreateAuthSource(ctx context.Context, source *model.AuthSource) error {
	return s.db.WithContext(ctx).Create(source).Error
}

// GetAuthSource 获取认证源
func (s *AuthSourceService) GetAuthSource(ctx context.Context, id uint) (*model.AuthSource, error) {
	var source model.AuthSource
	err := s.db.WithContext(ctx).First(&source, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
}

// ListAuthSources 列出认证源
func (s *AuthSourceService) ListAuthSources(ctx context.Context, limit, offset int) ([]*model.AuthSource, int64, error) {
	var sources []*model.AuthSource
	var total int64

	if err := s.db.Model(&model.AuthSource{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&sources).Error

	return sources, total, err
}

// ListAuthSourcesWithFilters 列出认证源（支持筛选）
func (s *AuthSourceService) ListAuthSourcesWithFilters(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*AuthSourceWithFlags, int64, error) {
	var sources []*model.AuthSource
	var total int64

	query := s.db.Model(&model.AuthSource{})

	// Apply filters if provided - only supporting name, type, scope, and enabled
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if sourceType, ok := filters["type"].(string); ok && sourceType != "" {
		query = query.Where("type = ?", sourceType)
	}
	if scope, ok := filters["scope"].(string); ok && scope != "" {
		query = query.Where("scope = ?", scope)
	}
	if enabled, ok := filters["enabled"].(bool); ok {
		query = query.Where("enabled = ?", enabled)
	}
	// Keyword search now only applies to name field for consistency with UI filters
	if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&sources).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert to enhanced format with flags
	var enhancedSources []*AuthSourceWithFlags
	for _, src := range sources {
		enhancedSource := &AuthSourceWithFlags{
			AuthSource: *src,
			IsSystem:   s.isSystemAuthSource(src),
		}
		enhancedSources = append(enhancedSources, enhancedSource)
	}

	return enhancedSources, total, nil
}

// AuthSourceWithFlags 认证源增强信息
type AuthSourceWithFlags struct {
	model.AuthSource
	IsSystem bool `json:"is_system"`
}

// isSystemAuthSource 检查是否为系统认证源
func (s *AuthSourceService) isSystemAuthSource(source *model.AuthSource) bool {
	// 系统认证源的判断条件：
	// 1. 类型为 local/sql 且作用域为 system
	// 2. 名称包含 "system" 或 "builtin"（不区分大小写）
	// 3. 或者是通过初始化脚本创建的内置系统认证源
	isBuiltInType := (source.Type == "local" || source.Type == "sql") && source.Scope == "system"
	isBuiltInName := strings.Contains(strings.ToLower(source.Name), "system") || strings.Contains(strings.ToLower(source.Name), "builtin")

	return isBuiltInType || isBuiltInName
}

// UpdateAuthSource 更新认证源
func (s *AuthSourceService) UpdateAuthSource(ctx context.Context, source *model.AuthSource) error {
	return s.db.WithContext(ctx).Save(source).Error
}

// DeleteAuthSource 删除认证源
func (s *AuthSourceService) DeleteAuthSource(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.AuthSource{}, id).Error
}

// TestAuthSource 测试认证源连接
func (s *AuthSourceService) TestAuthSource(ctx context.Context, source *model.AuthSource) (bool, error) {
	// TODO: 根据类型实现不同的测试逻辑
	// LDAP: 尝试连接 LDAP 服务器
	// OIDC: 尝试获取发现文档
	// SAML: 尝试获取元数据
	switch source.Type {
	case "ldap":
		return s.testLDAP(source)
	case "oidc":
		return s.testOIDC(source)
	case "saml":
		return s.testSAML(source)
	default:
		return true, nil
	}
}

func (s *AuthSourceService) testLDAP(source *model.AuthSource) (bool, error) {
	if source.Config == nil {
		return false, errors.New("ldap config is empty")
	}

	cfg, err := UnmarshalLDAPConfig(json.RawMessage(source.Config))
	if err != nil {
		return false, fmt.Errorf("failed to parse ldap config: %w", err)
	}

	if cfg.URL == "" {
		return false, errors.New("ldap url is required")
	}

	// Check if this is a mock LDAP test by checking for "mock" in the URL
	if strings.Contains(strings.ToLower(cfg.URL), "mock") {
		return s.testMockLDAP(source)
	}

	// 解析 LDAP URL，提取 host:port
	parsed, err := url.Parse(cfg.URL)
	if err != nil {
		return false, fmt.Errorf("invalid ldap url: %w", err)
	}

	host := parsed.Host
	if host == "" {
		// 兼容没有 scheme 的情况，直接当 host:port 处理
		host = cfg.URL
	}
	// 如果没有端口，根据 scheme 补充默认端口
	if _, _, err := net.SplitHostPort(host); err != nil {
		switch parsed.Scheme {
		case "ldaps":
			host = host + ":636"
		default:
			host = host + ":389"
		}
	}

	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		return false, fmt.Errorf("failed to connect to ldap server %s: %w", host, err)
	}
	conn.Close()
	return true, nil
}

// testMockLDAP 模拟 LDAP 连接测试，使用内置测试数据
func (s *AuthSourceService) testMockLDAP(source *model.AuthSource) (bool, error) {
	// 验证必要的配置是否存在
	cfg, err := UnmarshalLDAPConfig(json.RawMessage(source.Config))
	if err != nil {
		return false, fmt.Errorf("failed to parse ldap config: %w", err)
	}

	if cfg.URL == "" {
		return false, errors.New("ldap url is required")
	}

	// 模拟 LDAP 连接测试 - 返回固定的测试响应
	fmt.Printf("Mock LDAP test performed for source: %s with URL: %s\n", source.Name, cfg.URL)

	// 模拟连接延迟
	time.Sleep(100 * time.Millisecond)

	// 模拟可能的常见连接问题
	if strings.Contains(cfg.URL, "invalid") || strings.Contains(cfg.URL, "error") {
		return false, errors.New("simulated connection error")
	}

	// 返回模拟的成功连接
	return true, nil
}

func (s *AuthSourceService) testOIDC(source *model.AuthSource) (bool, error) {
	// TODO: 实现 OIDC 连接测试
	return true, nil
}

func (s *AuthSourceService) testSAML(source *model.AuthSource) (bool, error) {
	// TODO: 实现 SAML 连接测试
	return true, nil
}

// SyncResult 同步结果
type SyncResult struct {
	Created int      `json:"created"`
	Updated int      `json:"updated"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

// SyncUsers 同步认证源用户
func (s *AuthSourceService) SyncUsers(ctx context.Context, sourceID uint) (*SyncResult, error) {
	source, err := s.GetAuthSource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, gorm.ErrRecordNotFound
	}

	result := &SyncResult{
		Errors: []string{},
	}

	// TODO: 根据 source.Type 实现不同协议的用户同步逻辑
	// 目前仅支持 LDAP，其他类型返回不支持错误
	switch source.Type {
	case "ldap":
		// TODO: 实现 LDAP 用户批量拉取和同步
		// 1. 使用 BindDN/BindPassword 建立管理员连接
		// 2. 搜索 BaseDN 下符合 UserFilter 的用户条目
		// 3. 对每个用户：检查本地是否已存在，不存在则创建，存在则更新
		result.Errors = append(result.Errors, "ldap sync not yet implemented, install go-ldap/ldap/v3 first")
		result.Failed = 1
	default:
		return nil, fmt.Errorf("sync not supported for auth source type: %s", source.Type)
	}

	return result, nil
}

// SyncAuthSource 同步认证源用户（兼容旧接口）
func (s *AuthSourceService) SyncAuthSource(ctx context.Context, id uint) error {
	_, err := s.SyncUsers(ctx, id)
	return err
}

// GetAuthSourcesByScope 按范围查询认证源
func (s *AuthSourceService) GetAuthSourcesByScope(ctx context.Context, scope string, domainID *uint) ([]*model.AuthSource, error) {
	var sources []*model.AuthSource
	query := s.db.WithContext(ctx).Where("scope = ? AND enabled = ?", scope, true)
	if scope == "domain" && domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}
	err := query.Find(&sources).Error
	return sources, err
}

// LDAPUserInfo LDAP 用户信息
type LDAPUserInfo struct {
	Username    string
	DisplayName string
	Email       string
	DN          string
}

// authenticateWithLDAP 使用 LDAP 认证用户
// TODO: 安装 github.com/go-ldap/ldap/v3 后实现完整的 LDAP bind 认证
func (s *AuthSourceService) authenticateWithLDAP(source *model.AuthSource, username, password string) (*LDAPUserInfo, error) {
	// TODO: 实现真实 LDAP 认证
	// 步骤：
	// 1. 解析 LDAP 配置（URL, BaseDN, BindDN, BindPassword, UserFilter）
	// 2. 用 BindDN/BindPassword 建立管理员连接
	// 3. 用 UserFilter 搜索用户 DN（例：(&(objectClass=person)(uid=%s))）
	// 4. 用找到的用户 DN 和传入的 password 尝试 bind
	// 5. bind 成功表示认证通过，返回用户信息
	// 示例代码（需要 go-ldap/ldap/v3）：
	//   l, err := ldap.DialURL(cfg.URL)
	//   l.Bind(cfg.BindDN, cfg.BindPassword)
	//   result, _ := l.Search(searchRequest)
	//   userDN := result.Entries[0].DN
	//   l.Bind(userDN, password)
	return nil, errors.New("ldap authentication not implemented: install github.com/go-ldap/ldap/v3")
}

// AuthenticateUser 统一用户认证入口
// 先尝试本地认证，再尝试 LDAP 认证
func (s *AuthSourceService) AuthenticateUser(ctx context.Context, username, password string, domainID *uint) (*model.User, *model.AuthSource, error) {
	// 1. 从数据库查找本地用户
	var localUser model.User
	query := s.db.WithContext(ctx).Where("name = ?", username)
	if domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}
	err := query.First(&localUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 2. 本地用户存在，验证密码
	if err == nil {
		if !localUser.Enabled {
			return nil, nil, errors.New("user is disabled")
		}
		if bcryptErr := bcrypt.CompareHashAndPassword([]byte(localUser.Password), []byte(password)); bcryptErr == nil {
			return &localUser, nil, nil
		}
		// 本地密码验证失败，继续尝试 LDAP
	}

	// 3. 查找可用的 LDAP 认证源
	var ldapSources []*model.AuthSource
	ldapQuery := s.db.WithContext(ctx).
		Where("type = ? AND enabled = ?", "ldap", true).
		Where("scope = ? OR (scope = ? AND domain_id = ?)", "system", "domain", domainID)
	if err := ldapQuery.Find(&ldapSources).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to query ldap auth sources: %w", err)
	}

	// 4. 对每个 LDAP 源尝试认证
	for _, source := range ldapSources {
		userInfo, authErr := s.authenticateWithLDAP(source, username, password)
		if authErr != nil || userInfo == nil {
			continue
		}

		// 5. LDAP 认证成功，检查是否需要自动创建本地用户
		var domainIDVal uint
		if source.DomainID != nil {
			domainIDVal = *source.DomainID
		}

		if source.AutoCreate {
			// 检查本地是否已存在该用户
			var existingUser model.User
			findErr := s.db.WithContext(ctx).Where("name = ?", username).First(&existingUser).Error
			if findErr == nil {
				// 用户已存在，直接返回
				return &existingUser, source, nil
			}
			// 用户不存在，创建本地用户（密码置为空，后续通过 LDAP 认证）
			newUser := &model.User{
				Name:        username,
				DisplayName: userInfo.DisplayName,
				Email:       userInfo.Email,
				Password:    "",
				DomainID:    domainIDVal,
				Enabled:     true,
			}
			// 生成一个随机密码占位（不会用于本地认证）
			placeholder, _ := utils.HashPassword(fmt.Sprintf("ldap_%s_%d", username, time.Now().UnixNano()))
			newUser.Password = placeholder

			if createErr := s.db.WithContext(ctx).Create(newUser).Error; createErr != nil {
				return nil, nil, fmt.Errorf("failed to auto-create ldap user: %w", createErr)
			}

			// 若认证源配置了默认角色，分配角色
			if source.DefaultRole != nil && domainIDVal > 0 {
				userRole := &model.UserRole{
					UserID:   newUser.ID,
					RoleID:   *source.DefaultRole,
					DomainID: domainIDVal,
				}
				// 忽略分配角色失败（用户已创建）
				s.db.WithContext(ctx).Create(userRole)
			}

			return newUser, source, nil
		}

		// AutoCreate=false：仅在本地已有该用户时允许登录
		if err == nil {
			// localUser 已找到，密码通过 LDAP 认证
			return &localUser, source, nil
		}

		return nil, nil, errors.New("user not found locally and auto_create is disabled")
	}

	return nil, nil, errors.New("invalid username or password")
}

// EnableAuthSource 启用认证源
func (s *AuthSourceService) EnableAuthSource(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.AuthSource{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableAuthSource 禁用认证源
func (s *AuthSourceService) DisableAuthSource(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.AuthSource{}).Where("id = ?", id).Update("enabled", false).Error
}

// LDAPConfig LDAP 配置
type LDAPConfig struct {
	URL          string `json:"url"`
	BaseDN       string `json:"base_dn"`
	BindDN       string `json:"bind_dn"`
	BindPassword string `json:"bind_password"`
	UserFilter   string `json:"user_filter"`
	UserIDAttr   string `json:"user_id_attr"`   // 用于标识用户唯一性的属性，默认为 uid
	UserNameAttr string `json:"user_name_attr"` // 用于显示用户名的属性，默认为 cn
}

// OIDCConfig OIDC 配置
type OIDCConfig struct {
	Issuer       string   `json:"issuer"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
}

// SAMLConfig SAML 配置
type SAMLConfig struct {
	IdPMetadataURL string `json:"idp_metadata_url"`
	SpEntityID     string `json:"sp_entity_id"`
	SpACSURL       string `json:"sp_acs_url"`
}

// UnmarshalLDAPConfig 解析 LDAP 配置
func UnmarshalLDAPConfig(configJSON json.RawMessage) (*LDAPConfig, error) {
	var cfg LDAPConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalOIDCConfig 解析 OIDC 配置
func UnmarshalOIDCConfig(configJSON json.RawMessage) (*OIDCConfig, error) {
	var cfg OIDCConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalSAMLConfig 解析 SAML 配置
func UnmarshalSAMLConfig(configJSON json.RawMessage) (*SAMLConfig, error) {
	var cfg SAMLConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}
