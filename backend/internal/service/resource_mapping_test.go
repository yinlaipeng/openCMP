package service

import (
	"testing"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:test?cache=shared&mode=memory"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate
	db.AutoMigrate(&model.Rule{}, &model.RuleTag{}, &model.SyncPolicy{}, &model.CloudAccount{}, &model.Domain{}, &model.Project{})

	return db
}

func TestMatchRule_AllMatch(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	// Create a rule with all_match condition
	rule := model.Rule{
		ConditionType: "all_match",
		Tags: []model.RuleTag{
			{TagKey: "Project", TagValue: "project-alpha"},
			{TagKey: "Environment", TagValue: "prod"},
		},
	}

	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name: "all tags match",
			tags: map[string]string{
				"Project":     "project-alpha",
				"Environment": "prod",
			},
			expected: true,
		},
		{
			name: "partial match - missing one tag",
			tags: map[string]string{
				"Project": "project-alpha",
			},
			expected: false,
		},
		{
			name: "wrong value",
			tags: map[string]string{
				"Project":     "project-beta",
				"Environment": "prod",
			},
			expected: false,
		},
		{
			name: "no tags",
			tags: map[string]string{},
			expected: false,
		},
		{
			name: "extra tags but all required match",
			tags: map[string]string{
				"Project":     "project-alpha",
				"Environment": "prod",
				"Team":        "devops",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchRule(rule, tt.tags)
			if result != tt.expected {
				t.Errorf("matchRule() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMatchRule_AnyMatch(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	// Create a rule with any_match condition
	rule := model.Rule{
		ConditionType: "any_match",
		Tags: []model.RuleTag{
			{TagKey: "Project", TagValue: "project-alpha"},
			{TagKey: "Project", TagValue: "project-beta"},
		},
	}

	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name: "first option matches",
			tags: map[string]string{
				"Project": "project-alpha",
			},
			expected: true,
		},
		{
			name: "second option matches",
			tags: map[string]string{
				"Project": "project-beta",
			},
			expected: true,
		},
		{
			name: "no match",
			tags: map[string]string{
				"Project": "project-gamma",
			},
			expected: false,
		},
		{
			name: "empty tags",
			tags: map[string]string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchRule(rule, tt.tags)
			if result != tt.expected {
				t.Errorf("matchRule() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMatchRule_KeyMatch(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	// Create a rule with key_match condition
	rule := model.Rule{
		ConditionType: "key_match",
		Tags: []model.RuleTag{
			{TagKey: "Project", TagValue: ""}, // Value not checked
			{TagKey: "Environment", TagValue: ""},
		},
	}

	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name: "key exists - Project",
			tags: map[string]string{
				"Project": "any-value",
			},
			expected: true,
		},
		{
			name: "key exists - Environment",
			tags: map[string]string{
				"Environment": "prod",
			},
			expected: true,
		},
		{
			name: "both keys exist",
			tags: map[string]string{
				"Project":     "alpha",
				"Environment": "prod",
			},
			expected: true,
		},
		{
			name: "key not found",
			tags: map[string]string{
				"Team": "devops",
			},
			expected: false,
		},
		{
			name: "empty tags",
			tags: map[string]string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchRule(rule, tt.tags)
			if result != tt.expected {
				t.Errorf("matchRule() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMatchTag_ExactMatch(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	ruleTag := model.RuleTag{
		TagKey:   "Project",
		TagValue: "project-alpha",
	}

	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name:     "exact match",
			tags:     map[string]string{"Project": "project-alpha"},
			expected: true,
		},
		{
			name:     "case sensitive mismatch",
			tags:     map[string]string{"Project": "PROJECT-ALPHA"},
			expected: false,
		},
		{
			name:     "different value",
			tags:     map[string]string{"Project": "project-beta"},
			expected: false,
		},
		{
			name:     "key missing",
			tags:     map[string]string{"Team": "devops"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchTag(ruleTag, tt.tags)
			if result != tt.expected {
				t.Errorf("matchTag() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMatchTag_RegexMatch(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	ruleTag := model.RuleTag{
		TagKey:   "Project",
		TagValue: "project-[a-z]+", // regex pattern
	}

	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name:     "matches pattern",
			tags:     map[string]string{"Project": "project-alpha"},
			expected: true,
		},
		{
			name:     "matches pattern - beta",
			tags:     map[string]string{"Project": "project-beta"},
			expected: true,
		},
		{
			name:     "does not match - numbers",
			tags:     map[string]string{"Project": "project-123"},
			expected: false,
		},
		{
			name:     "does not match - wrong prefix",
			tags:     map[string]string{"Project": "prod-alpha"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchTag(ruleTag, tt.tags)
			if result != tt.expected {
				t.Errorf("matchTag() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestParseResourceTags(t *testing.T) {
	db := setupTestDB(t)
	logger := zap.NewNop()
	service := NewResourceMappingService(db, logger)

	// Test with VM
	vm := &cloudprovider.VirtualMachine{
		ID:   "vm-123",
		Name: "test-vm",
		Tags: map[string]string{
			"Project":     "project-alpha",
			"Environment": "prod",
		},
	}

	tags := service.ParseResourceTags(vm)
	if tags["Project"] != "project-alpha" {
		t.Errorf("Expected Project tag to be 'project-alpha', got '%s'", tags["Project"])
	}
	if tags["Environment"] != "prod" {
		t.Errorf("Expected Environment tag to be 'prod', got '%s'", tags["Environment"])
	}

	// Test with nil/empty tags
	vpc := &cloudprovider.VPC{
		ID:   "vpc-123",
		Name: "test-vpc",
		Tags: map[string]string{},
	}

	tags = service.ParseResourceTags(vpc)
	if len(tags) != 0 {
		t.Errorf("Expected empty tags, got %d tags", len(tags))
	}
}