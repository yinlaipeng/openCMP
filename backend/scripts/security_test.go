//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func main() {
	fmt.Println("🔐 开始安全性测试...")

	// 1. 测试未经身份验证的访问
	fmt.Println("\n🔍 测试未经身份验证的访问...")
	if err := testUnauthenticatedAccess(); err != nil {
		fmt.Printf("❌ 未经身份验证的访问测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 未经身份验证的访问被正确拒绝")
	}

	// 2. 测试无效令牌访问
	fmt.Println("\n🔍 测试无效令牌访问...")
	if err := testInvalidTokenAccess(); err != nil {
		fmt.Printf("❌ 无效令牌访问测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 无效令牌访问被正确拒绝")
	}

	// 3. 测试暴力破解防护
	fmt.Println("\n🔍 测试暴力破解防护...")
	if err := testBruteForceProtection(); err != nil {
		fmt.Printf("❌ 暴力破解防护测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 暴力破解防护测试完成")
	}

	// 4. 测试权限绕过
	fmt.Println("\n🔍 测试权限绕过...")
	if err := testPrivilegeEscalation(); err != nil {
		fmt.Printf("❌ 权限绕过测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 权限绕过测试完成")
	}

	// 5. 测试敏感信息泄露
	fmt.Println("\n🔍 测试敏感信息泄露...")
	if err := testSensitiveInformationDisclosure(); err != nil {
		fmt.Printf("❌ 敏感信息泄露测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 敏感信息泄露测试完成")
	}

	// 6. 测试SQL注入
	fmt.Println("\n🔍 测试SQL注入...")
	if err := testSQLInjection(); err != nil {
		fmt.Printf("❌ SQL注入测试失败: %v\n", err)
	} else {
		fmt.Println("✅ SQL注入测试完成")
	}

	// 7. 测试认证功能
	fmt.Println("\n🔍 测试认证功能...")
	if token, err := testAuthenticationFlow(); err != nil {
		fmt.Printf("❌ 认证功能测试失败: %v\n", err)
	} else {
		fmt.Printf("✅ 认证功能测试成功，获取到令牌: %s\n", token)
		
		// 使用有效令牌测试授权
		fmt.Println("\n🔍 使用有效令牌测试授权...")
		if err := testAuthorizationWithValidToken(token); err != nil {
			fmt.Printf("❌ 授权测试失败: %v\n", err)
		} else {
			fmt.Println("✅ 授权测试成功")
		}
	}

	fmt.Println("\n🛡️  安全性测试完成！")
}

// testUnauthenticatedAccess 测试未经身份验证的访问
func testUnauthenticatedAccess() error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("预期状态码 %d，实际得到 %d", http.StatusUnauthorized, resp.StatusCode)
	}

	return nil
}

// testInvalidTokenAccess 测试无效令牌访问
func testInvalidTokenAccess() error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_12345")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("预期状态码 %d，实际得到 %d", http.StatusUnauthorized, resp.StatusCode)
	}

	return nil
}

// testBruteForceProtection 测试暴力破解防护
func testBruteForceProtection() error {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// 尝试多次无效登录
	for i := 0; i < 5; i++ {
		reqBody, _ := json.Marshal(LoginRequest{
			Username: "nonexistent",
			Password: "wrongpassword",
		})

		resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	// 再次尝试有效登录，看是否被限制
	reqBody, _ := json.Marshal(LoginRequest{
		Username: "admin",
		Password: "admin123",
	})

	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 应该仍然允许有效用户登录（如果系统实现了智能限流）
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️  检测到可能的暴力破解防护机制（状态码: %d, 响应: %s）\n", resp.StatusCode, string(body))
	}

	return nil
}

// testPrivilegeEscalation 测试权限提升
func testPrivilegeEscalation() error {
	// 首先使用普通用户登录（假设存在）
	token, err := login("demo", "demopass")
	if err != nil {
		// 如果普通用户不存在，跳过此测试
		fmt.Println("⚠️  普通用户不存在，跳过权限提升测试")
		return nil
	}

	// 尝试执行管理员操作（例如创建用户）
	client := &http.Client{Timeout: 10 * time.Second}
	userData := map[string]interface{}{
		"name":      "privilege_test",
		"display_name": "Privilege Escalation Test",
		"email":     "privilege_test@example.com",
		"password":  "Password123!",
		"domain_id": 1,
	}
	reqBody, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", baseURL+"/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 普通用户不应该能够创建其他用户
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return fmt.Errorf("安全漏洞：普通用户能够执行管理员操作")
	}

	fmt.Printf("✅ 正确阻止了权限提升尝试（状态码: %d）\n", resp.StatusCode)
	return nil
}

// testSensitiveInformationDisclosure 测试敏感信息泄露
func testSensitiveInformationDisclosure() error {
	// 尝试获取用户列表，检查是否包含敏感信息
	token, err := login("admin", "admin123")
	if err != nil {
		return fmt.Errorf("无法登录以进行测试: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("获取用户列表失败: %s", string(body))
	}

	// 检查响应中是否包含明文密码
	if containsSensitiveInfo(string(body)) {
		return fmt.Errorf("发现敏感信息泄露：响应中包含潜在的敏感数据")
	}

	return nil
}

// containsSensitiveInfo 检查响应中是否包含敏感信息
func containsSensitiveInfo(response string) bool {
	// 检查常见的敏感字段
	sensitiveKeywords := []string{
		"password", "passwd", "pwd", "secret", "key", "token", "authorization",
		"api_key", "private_key", "secret_key", "credentials",
	}

	responseLower := fmt.Sprintf("%s", response)
	for _, keyword := range sensitiveKeywords {
		if containsIgnoreCase(responseLower, keyword) {
			// 进一步检查是否在密码字段之外的地方出现
			return true
		}
	}

	return false
}

// containsIgnoreCase 检查字符串是否包含另一个字符串（忽略大小写）
func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) &&
		containsSubstring(toLower(str), toLower(substr))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if str[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result = append(result, c+('a'-'A'))
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

// testSQLInjection 测试SQL注入
func testSQLInjection() error {
	token, err := login("admin", "admin123")
	if err != nil {
		return fmt.Errorf("无法登录以进行测试: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	
	// 尝试在用户名参数中注入SQL
	maliciousUsername := "admin' OR '1'='1"
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users?name=%s", baseURL, maliciousUsername), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("✅ 正确阻止了SQL注入尝试（状态码: %d）\n", resp.StatusCode)
		return nil
	}

	// 检查响应是否异常（例如返回过多数据）
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	
	if items, ok := result["items"].([]interface{}); ok {
		if len(items) > 100 { // 假设正常情况下不会有这么多用户
			return fmt.Errorf("可能存在SQL注入漏洞：返回了异常多的数据项")
		}
	}

	fmt.Println("✅ 未检测到明显的SQL注入漏洞")
	return nil
}

// testAuthenticationFlow 测试认证流程
func testAuthenticationFlow() (string, error) {
	// 尝试使用正确的凭据登录
	token, err := login("admin", "admin123")
	if err != nil {
		return "", err
	}

	// 使用获取的令牌访问受保护的资源
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/auth/current-user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("使用有效令牌访问失败，状态码: %d", resp.StatusCode)
	}

	return token, nil
}

// testAuthorizationWithValidToken 使用有效令牌测试授权
func testAuthorizationWithValidToken(token string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// 尝试访问用户列表（需要适当权限）
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("授权失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// login 登录
func login(username, password string) (string, error) {
	reqBody, _ := json.Marshal(LoginRequest{
		Username: username,
		Password: password,
	})

	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", err
	}

	return loginResp.Token, nil
}