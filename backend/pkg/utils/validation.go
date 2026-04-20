package utils

import (
	"strings"
)

// EscapeLikePattern 转义SQL LIKE模式中的特殊字符
// 防止用户输入中的 % 和 _ 被解释为通配符
func EscapeLikePattern(s string) string {
	// 转义顺序很重要：先转义 % 再转义 _
	result := strings.ReplaceAll(s, "%", "\\%")
	result = strings.ReplaceAll(result, "_", "\\_")
	return result
}

// Contains 检查字符串切片是否包含指定元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidProviderTypes 返回有效的云厂商类型列表
func ValidProviderTypes() []string {
	return []string{"alibaba", "tencent", "aws", "azure", "huawei", "google", "openstack", "vmware"}
}

// IsValidProviderType 检查云厂商类型是否有效
func IsValidProviderType(providerType string) bool {
	return Contains(ValidProviderTypes(), providerType)
}