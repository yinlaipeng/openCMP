package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// 1. 测试登录
	fmt.Println("=== Testing Login ===")
	loginResp, err := login("admin", "admin123")
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Printf("Login successful, token: %s\n\n", loginResp.Token)

	// 2. 测试获取当前用户信息
	fmt.Println("=== Testing Get Current User ===")
	userInfo, err := getCurrentUser(loginResp.Token)
	if err != nil {
		fmt.Printf("Get current user failed: %v\n", err)
		return
	}
	fmt.Printf("Current user: %+v\n\n", userInfo)

	// 3. 测试权限检查
	fmt.Println("=== Testing Permission Check ===")
	hasPermission, err := checkPermission(loginResp.Token, "user", "list")
	if err != nil {
		fmt.Printf("Check permission failed: %v\n", err)
		return
	}
	fmt.Printf("Has permission to user:list: %t\n\n", hasPermission)

	// 4. 测试列出用户（带角色信息）
	fmt.Println("=== Testing List Users With Roles ===")
	usersWithRoles, err := listUsersWithRoles(loginResp.Token)
	if err != nil {
		fmt.Printf("List users with roles failed: %v\n", err)
		return
	}
	fmt.Printf("Number of users with roles: %d\n\n", len(usersWithRoles))

	// 5. 测试获取特定用户权限
	fmt.Println("=== Testing Get User Permissions ===")
	userPermissions, err := getUserPermissions(loginResp.Token, 1) // 假设admin用户ID为1
	if err != nil {
		fmt.Printf("Get user permissions failed: %v\n", err)
		return
	}
	fmt.Printf("Number of permissions for user 1: %d\n\n", len(userPermissions))
}

func login(username, password string) (*LoginResponse, error) {
	reqBody, _ := json.Marshal(LoginRequest{
		Username: username,
		Password: password,
	})

	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func getCurrentUser(token string) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/auth/current-user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get current user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo interface{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func checkPermission(token, resource, action string) (bool, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"resource": resource,
		"action":   action,
	})

	client := &http.Client{}
	req, _ := http.NewRequest("POST", baseURL+"/iam/check-permission", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("check permission failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]bool
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, err
	}

	return result["allowed"], nil
}

func listUsersWithRoles(token string) ([]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/iam/users-with-roles?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list users with roles failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	users, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return users, nil
}

func getUserPermissions(token string, userID int) ([]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/iam/users/%d/permissions", baseURL, userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user permissions failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	permissions, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return permissions, nil
}