package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/opencmp/opencmp/internal/model"
	"gorm.io/gorm"
)

// This is a simple test to demonstrate the mock LDAP functionality
func main() {
	// Example of how the mock LDAP functionality would work
	config := map[string]interface{}{
		"url": "mock://test-ldap-server",
		"base_dn": "dc=example,dc=com",
		"bind_dn": "cn=admin,dc=example,dc=com",
		"bind_password": "admin",
		"user_filter": "(objectClass=person)",
		"user_id_attr": "uid",
		"user_name_attr": "cn",
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Error marshaling config:", err)
	}

	authSource := &model.AuthSource{
		Name:   "Test Mock LDAP",
		Type:   "ldap",
		Scope:  "domain",
		Config: configJSON,
	}

	fmt.Println("Authentication Source Configured:")
	fmt.Printf("- Name: %s\n", authSource.Name)
	fmt.Printf("- Type: %s\n", authSource.Type)
	fmt.Printf("- Config URL includes 'mock': %t\n", containsMock(authSource))
	fmt.Println("\nThe mock LDAP functionality will simulate connection when URL contains 'mock'")
	fmt.Println("- Successful connections return true with no actual LDAP server needed")
	fmt.Println("- Connections with 'invalid' or 'error' in URL simulate failure")
	fmt.Println("- Normal LDAP URLs will attempt actual connection")
}

func containsMock(source *model.AuthSource) bool {
	return source.Config != nil &&
		   len(source.Config) > 0 &&
		   (json.RawMessage(source.Config)).Contains([]byte("mock"))
}