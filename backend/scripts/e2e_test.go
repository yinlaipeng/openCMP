//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	// 检查服务器是否运行
	if !isServerRunning() {
		fmt.Println("ERROR: 服务器未运行，请先启动后端服务")
		os.Exit(1)
	}

	fmt.Println("SUCCESS: 服务器正在运行，开始端到端测试...")

	// 1. 测试用户登录
	fmt.Println("\nTEST: 测试用户登录...")
	token, err := login("admin", "admin123")
	if err != nil {
		fmt.Printf("FAILED: 登录失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: 登录成功，获取到令牌: %s\n", token)

	// 2. 测试获取当前用户信息
	fmt.Println("\nTEST: 测试获取当前用户信息...")
	userInfo, err := getCurrentUser(token)
	if err != nil {
		fmt.Printf("FAILED: 获取当前用户信息失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: 成功获取当前用户信息: %+v\n", userInfo)

	// 3. 测试列出用户
	fmt.Println("\nTEST: 测试列出用户...")
	users, err := listUsers(token)
	if err != nil {
		fmt.Printf("FAILED: 列出用户失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: 成功列出用户，共 %d 个用户\n", len(users))

	// 4. 测试创建新用户
	fmt.Println("\nTEST: 测试创建新用户...")
	newUserID, err := createUser(token)
	if err != nil {
		fmt.Printf("FAILED: 创建用户失败: %v\n", err)
		// 不退出，继续测试其他功能
	} else {
		fmt.Printf("SUCCESS: 成功创建新用户，ID: %d\n", newUserID)
	}

	// 5. 测试列出角色
	fmt.Println("\nTEST: 测试列出角色...")
	roles, err := listRoles(token)
	if err != nil {
		fmt.Printf("FAILED: 列出角色失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: 成功列出角色，共 %d 个角色\n", len(roles))

	// 6. 测试列出权限
	fmt.Println("\nTEST: 测试列出权限...")
	permissions, err := listPermissions(token)
	if err != nil {
		fmt.Printf("FAILED: 列出权限失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: 成功列出权限，共 %d 个权限\n", len(permissions))

	// 7. 测试认证源功能
	fmt.Println("\nTEST: 测试认证源功能...")
	authSources, err := listAuthSources(token)
	if err != nil {
		fmt.Printf("WARNING: 列出认证源失败: %v\n", err)
		// 这可能不是致命错误，继续测试
	} else {
		fmt.Printf("SUCCESS: 成功列出认证源，共 %d 个认证源\n", len(authSources))
	}

	fmt.Println("\nCOMPLETE: 所有端到端测试完成！")
}

// isServerRunning 检查服务器是否运行
func isServerRunning() bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
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

// getCurrentUser 获取当前用户信息
func getCurrentUser(token string) (interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/auth/current-user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取当前用户失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var userInfo interface{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// listUsers 列出用户
func listUsers(token string) ([]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出用户失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	users, ok := result["items"].([]interface{})
	if !ok {
		// 如果返回格式不同，尝试直接解析为数组
		var usersArray []interface{}
		err = json.Unmarshal(body, &usersArray)
		if err != nil {
			return nil, fmt.Errorf("无法解析用户列表响应: %v", err)
		}
		return usersArray, nil
	}

	return users, nil
}

// createUser 创建用户
func createUser(token string) (int, error) {
	userData := map[string]interface{}{
		"name":         "testuser_e2e",
		"display_name": "E2E Test User",
		"email":        "testuser_e2e@example.com",
		"password":     "Password123!",
		"domain_id":    1,
	}

	reqBody, _ := json.Marshal(userData)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", baseURL+"/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("创建用户失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	// 尝试获取新用户的ID
	if idFloat, ok := result["id"].(float64); ok {
		return int(idFloat), nil
	}

	return 0, nil
}

// listRoles 列出角色
func listRoles(token string) ([]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/roles", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出角色失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	roles, ok := result["items"].([]interface{})
	if !ok {
		// 如果返回格式不同，尝试直接解析为数组
		var rolesArray []interface{}
		err = json.Unmarshal(body, &rolesArray)
		if err != nil {
			return nil, fmt.Errorf("无法解析角色列表响应: %v", err)
		}
		return rolesArray, nil
	}

	return roles, nil
}

// listPermissions 列出权限
func listPermissions(token string) ([]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/permissions", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出权限失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	permissions, ok := result["items"].([]interface{})
	if !ok {
		// 如果返回格式不同，尝试直接解析为数组
		var permsArray []interface{}
		err = json.Unmarshal(body, &permsArray)
		if err != nil {
			return nil, fmt.Errorf("无法解析权限列表响应: %v", err)
		}
		return permsArray, nil
	}

	return permissions, nil
}

// listAuthSources 列出认证源
func listAuthSources(token string) ([]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/auth-sources", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出认证源失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	authSources, ok := result["items"].([]interface{})
	if !ok {
		// 如果返回格式不同，尝试直接解析为数组
		var sourcesArray []interface{}
		err = json.Unmarshal(body, &sourcesArray)
		if err != nil {
			return nil, fmt.Errorf("无法解析认证源列表响应: %v", err)
		}
		return sourcesArray, nil
	}

	return authSources, nil
}
