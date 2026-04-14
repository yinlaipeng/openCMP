package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

func TestHostTemplateService(t *testing.T) {
	// Mock DB connection for testing
	db := &gorm.DB{}
	service := NewHostTemplateService(db)

	t.Run("CreateHostTemplate", func(t *testing.T) {
		// Since we don't have a real DB, just test that the function exists
		assert.NotNil(t, service)
		assert.NotNil(t, service.CreateHostTemplate)
	})

	t.Run("ValidateHostTemplateConfig", func(t *testing.T) {
		req := &cloudprovider.HostTemplate{
			Name:        "",
			InstanceType: "",
			ImageID:     "",
			ProjectID:   "",
			Platform:    "",
		}

		err := service.ValidateHostTemplateConfig(context.Background(), req)
		assert.Error(t, err)

		req.Name = "test-template"
		req.InstanceType = "ecs.t5-lc1m1.small"
		req.ImageID = "img-12345"
		req.ProjectID = "project-12345"
		req.Platform = "alibaba"

		err = service.ValidateHostTemplateConfig(context.Background(), req)
		assert.NoError(t, err)
	})
}