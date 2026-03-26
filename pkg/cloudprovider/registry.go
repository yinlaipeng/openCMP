package cloudprovider

import (
	"fmt"
	"sync"
)

// ProviderFactory 云提供商工厂函数
type ProviderFactory func(config CloudAccountConfig) (ICloudProvider, error)

var (
	providerRegistry = make(map[string]ProviderFactory)
	registryMutex    sync.RWMutex
)

// RegisterProvider 注册云厂商适配器
func RegisterProvider(providerType string, factory ProviderFactory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	providerRegistry[providerType] = factory
}

// GetProvider 获取云提供商实例
func GetProvider(providerType string, config CloudAccountConfig) (ICloudProvider, error) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	factory, ok := providerRegistry[providerType]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}
	return factory(config)
}

// ListProviders 列出所有已注册的云厂商
func ListProviders() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	providers := make([]string, 0, len(providerRegistry))
	for p := range providerRegistry {
		providers = append(providers, p)
	}
	return providers
}
