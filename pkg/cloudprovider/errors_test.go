package cloudprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CloudError
		expected string
	}{
		{
			name: "with detail",
			err: &CloudError{
				Code:    ErrProviderNotFound,
				Message: "provider not found",
				Detail:  "aws",
			},
			expected: "[ProviderNotFound] provider not found: aws",
		},
		{
			name: "without detail",
			err: &CloudError{
				Code:    ErrProviderNotFound,
				Message: "provider not found",
			},
			expected: "[ProviderNotFound] provider not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestNewCloudError(t *testing.T) {
	err := NewCloudError(ErrProviderNotFound, "provider not found", "aws")
	assert.Equal(t, ErrProviderNotFound, err.Code)
	assert.Equal(t, "provider not found", err.Message)
	assert.Equal(t, "aws", err.Detail)
}

func TestNewCloudError_NoDetail(t *testing.T) {
	err := NewCloudError(ErrProviderNotFound, "provider not found")
	assert.Equal(t, ErrProviderNotFound, err.Code)
	assert.Equal(t, "provider not found", err.Message)
	assert.Empty(t, err.Detail)
}
