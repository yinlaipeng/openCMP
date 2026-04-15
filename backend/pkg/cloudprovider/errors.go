package cloudprovider

import "fmt"

// ErrorCode 错误码
type ErrorCode string

const (
	ErrProviderNotFound      ErrorCode = "ProviderNotFound"
	ErrInvalidCredentials    ErrorCode = "InvalidCredentials"
	ErrResourceNotFound      ErrorCode = "ResourceNotFound"
	ErrResourceAlreadyExists ErrorCode = "ResourceAlreadyExists"
	ErrOperationFailed       ErrorCode = "OperationFailed"
	ErrTimeout               ErrorCode = "Timeout"
	ErrUnsupportedOperation  ErrorCode = "UnsupportedOperation"
	ErrProviderError         ErrorCode = "ProviderError"
)

// CloudError 云操作错误
type CloudError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

func (e *CloudError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewCloudError 创建云操作错误
func NewCloudError(code ErrorCode, message string, detail ...string) *CloudError {
	err := &CloudError{
		Code:    code,
		Message: message,
	}
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	return err
}
