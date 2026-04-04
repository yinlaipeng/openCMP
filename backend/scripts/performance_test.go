//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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

// PerformanceResult 性能测试结果
type PerformanceResult struct {
	TestName     string
	AvgTime      time.Duration
	MinTime      time.Duration
	MaxTime      time.Duration
	TotalTime    time.Duration
	RequestCount int
	ErrorCount   int
	TPS          float64 // Transactions Per Second
}

func main() {
	fmt.Println("🚀 开始性能测试...")

	// 首先获取一个有效的令牌
	token, err := login("admin", "admin123")
	if err != nil {
		fmt.Printf("❌ 无法登录以进行性能测试: %v
", err)
		return
	}
	fmt.Printf("✅ 成功获取令牌，开始性能测试
")

	// 1. 并发用户登录测试
	fmt.Println("
👥 测试并发用户登录...")
	loginResults := testConcurrentLogins()
	printPerformanceResult(loginResults)

	// 2. API响应时间测试
	fmt.Println("
⏱️  测试API响应时间...")
	apiResults := testAPIResponseTime(token)
	printPerformanceResult(apiResults)

	// 3. 并发用户操作测试
	fmt.Println("
👥 测试并发用户操作...")
	concurrentResults := testConcurrentUserOperations(token)
	printPerformanceResult(concurrentResults)

	// 4. 持续负载测试
	fmt.Println("
持久负载测试...")
	loadResults := testSustainedLoad(token)
	printPerformanceResult(loadResults)

	fmt.Println("
📊 性能测试完成！")
}

// testConcurrentLogins 测试并发登录
func testConcurrentLogins() PerformanceResult {
	startTime := time.Now()
	var wg sync.WaitGroup
	const numConcurrent = 10
	errors := make(chan error, numConcurrent)
	
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := login("admin", "admin123")
			if err != nil {
				errors <- err
			}
		}()
	}
	
	wg.Wait()
	close(errors)
	
	errorCount := 0
	for range errors {
		errorCount++
	}
	
	totalTime := time.Since(startTime)
	avgTime := totalTime / time.Duration(numConcurrent)
	tps := float64(numConcurrent) / totalTime.Seconds()
	
	return PerformanceResult{
		TestName:     "并发登录测试",
		AvgTime:      avgTime,
		MinTime:      avgTime, // 简化处理
		MaxTime:      avgTime, // 简化处理
		TotalTime:    totalTime,
		RequestCount: numConcurrent,
		ErrorCount:   errorCount,
		TPS:          tps,
	}
}

// testAPIResponseTime 测试API响应时间
func testAPIResponseTime(token string) PerformanceResult {
	client := &http.Client{Timeout: 30 * time.Second}
	
	const numRequests = 20
	times := make([]time.Duration, numRequests)
	errors := 0
	
	for i := 0; i < numRequests; i++ {
		start := time.Now()
		
		req, _ := http.NewRequest("GET", baseURL+"/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := client.Do(req)
		if err != nil {
			errors++
			continue
		}
		io.ReadAll(resp.Body) // 读取响应体
		resp.Body.Close()
		
		times[i] = time.Since(start)
		time.Sleep(100 * time.Millisecond) // 避免过于频繁的请求
	}
	
	if len(times) == 0 {
		return PerformanceResult{TestName: "API响应时间测试"}
	}
	
	var totalTime time.Duration
	var minTime, maxTime time.Duration = times[0], times[0]
	
	for _, t := range times {
		totalTime += t
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
	}
	
	avgTime := totalTime / time.Duration(len(times))
	tps := float64(numRequests) / (totalTime.Seconds())
	
	return PerformanceResult{
		TestName:     "API响应时间测试",
		AvgTime:      avgTime,
		MinTime:      minTime,
		MaxTime:      maxTime,
		TotalTime:    totalTime,
		RequestCount: numRequests,
		ErrorCount:   errors,
		TPS:          tps,
	}
}

// testConcurrentUserOperations 测试并发用户操作
func testConcurrentUserOperations(token string) PerformanceResult {
	client := &http.Client{Timeout: 30 * time.Second}
	
	const numConcurrent = 5
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalDuration time.Duration
	var errorCount int
	var requestCount int
	
	startTime := time.Now()
	
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(userNum int) {
			defer wg.Done()
			
			for j := 0; j < 3; j++ { // 每个并发用户执行3次操作
				start := time.Now()
				
				req, _ := http.NewRequest("GET", baseURL+"/users", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				
				resp, err := client.Do(req)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}
				io.ReadAll(resp.Body) // 读取响应体
				resp.Body.Close()
				
				duration := time.Since(start)
				mu.Lock()
				totalDuration += duration
				requestCount++
				mu.Unlock()
				
				time.Sleep(200 * time.Millisecond) // 请求间隔
			}
		}(i)
	}
	
	wg.Wait()
	totalTestTime := time.Since(startTime)
	
	avgTime := time.Duration(0)
	if requestCount > 0 {
		avgTime = totalDuration / time.Duration(requestCount)
	}
	tps := float64(requestCount) / totalTestTime.Seconds()
	
	return PerformanceResult{
		TestName:     "并发用户操作测试",
		AvgTime:      avgTime,
		MinTime:      avgTime, // 简化处理
		MaxTime:      avgTime, // 简化处理
		TotalTime:    totalTestTime,
		RequestCount: requestCount,
		ErrorCount:   errorCount,
		TPS:          tps,
	}
}

// testSustainedLoad 持续负载测试
func testSustainedLoad(token string) PerformanceResult {
	client := &http.Client{Timeout: 30 * time.Second}
	
	const testDuration = 10 * time.Second // 持续10秒
	const requestInterval = 500 * time.Millisecond // 每500毫秒一个请求
	
	var totalDuration time.Duration
	var errorCount int
	var requestCount int
	
	startTime := time.Now()
	endTime := startTime.Add(testDuration)
	
	for time.Now().Before(endTime) {
		start := time.Now()
		
		req, _ := http.NewRequest("GET", baseURL+"/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := client.Do(req)
		if err != nil {
			errorCount++
		} else {
			io.ReadAll(resp.Body) // 读取响应体
			resp.Body.Close()
		}
		
		duration := time.Since(start)
		totalDuration += duration
		requestCount++
		
		time.Sleep(requestInterval)
	}
	
	totalTestTime := time.Since(startTime)
	avgTime := time.Duration(0)
	if requestCount > 0 {
		avgTime = totalDuration / time.Duration(requestCount)
	}
	tps := float64(requestCount) / testDuration.Seconds()
	
	return PerformanceResult{
		TestName:     "持续负载测试",
		AvgTime:      avgTime,
		MinTime:      avgTime, // 简化处理
		MaxTime:      avgTime, // 简化处理
		TotalTime:    totalTestTime,
		RequestCount: requestCount,
		ErrorCount:   errorCount,
		TPS:          tps,
	}
}

// printPerformanceResult 打印性能测试结果
func printPerformanceResult(result PerformanceResult) {
	fmt.Printf("  测试名称: %s
", result.TestName)
	fmt.Printf("  平均响应时间: %v
", result.AvgTime)
	fmt.Printf("  最小响应时间: %v
", result.MinTime)
	fmt.Printf("  最大响应时间: %v
", result.MaxTime)
	fmt.Printf("  总请求数: %d
", result.RequestCount)
	fmt.Printf("  错误数: %d
", result.ErrorCount)
	fmt.Printf("  TPS (每秒事务数): %.2f
", result.TPS)
	
	if result.ErrorCount > 0 {
		fmt.Printf("  ⚠️  发现 %d 个错误
", result.ErrorCount)
	}
	
	if result.AvgTime > 2*time.Second {
		fmt.Printf("  ⚠️  平均响应时间较长
")
	} else if result.AvgTime > 1*time.Second {
		fmt.Printf("  ⚠️  平均响应时间偏长
")
	} else {
		fmt.Printf("  ✅ 响应时间良好
")
	}
	
	fmt.Println()
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
