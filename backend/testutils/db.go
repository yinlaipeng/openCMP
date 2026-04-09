package testutils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the schema
	_ = db.AutoMigrate(
		// Core models
		&model.CloudAccount{},
		&model.Domain{},
		&model.Project{},
		&model.User{},
		&model.Group{},
		&model.UserGroup{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.ProjectUserRole{},
		&model.GroupRole{},
		&model.AuthSource{},
		&model.SecurityAlert{},
		&model.MessageType{},
		&model.Message{},
		&model.NotificationChannel{},
		&model.MessageSubscription{},
		&model.Robot{},
		&model.Receiver{},
		&model.Policy{},
		&model.PolicyStatement{},
		&model.RolePolicy{},
	)

	return db
}

// TeardownTestDB cleans up the test database
// For in-memory SQLite, cleanup is automatic
func TeardownTestDB(db *gorm.DB) {
	// No explicit cleanup needed for in-memory SQLite
}