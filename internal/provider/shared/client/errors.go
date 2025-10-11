package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an API error response
type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Details    string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error
func (e *Error) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func (e *Error) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden error
func (e *Error) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsValidationError returns true if the error is a 400 Bad Request error
func (e *Error) IsValidationError() bool {
	return e.StatusCode == http.StatusBadRequest
}

// IsConflict returns true if the error is a 409 Conflict error
func (e *Error) IsConflict() bool {
	return e.StatusCode == http.StatusConflict
}

// IsServerError returns true if the error is a 5xx server error
func (e *Error) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// CheckResponse checks the API response for errors
func CheckResponse(resp *Response) error {
	if resp == nil {
		return &Error{StatusCode: 0, Message: "nil response"}
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	apiError := &Error{
		StatusCode: resp.StatusCode,
		Message:    http.StatusText(resp.StatusCode),
	}

	// Try to parse error response body
	if len(resp.Body) > 0 {
		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
			Code    string `json:"code"`
			Details string `json:"details"`
		}

		if err := json.Unmarshal(resp.Body, &errorResp); err == nil {
			if errorResp.Message != "" {
				apiError.Message = errorResp.Message
			} else if errorResp.Error != "" {
				apiError.Message = errorResp.Error
			}
			apiError.Code = errorResp.Code
			apiError.Details = errorResp.Details
		}
	}

	return apiError
}

// IsNotFoundError returns true if the error is a 404 Not Found error
func IsNotFoundError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsNotFound()
	}
	return false
}

// IsUnauthorizedError returns true if the error is a 401 Unauthorized error
func IsUnauthorizedError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsUnauthorized()
	}
	return false
}

// IsForbiddenError returns true if the error is a 403 Forbidden error
func IsForbiddenError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsForbidden()
	}
	return false
}

// IsValidationError returns true if the error is a 400 Bad Request error
func IsValidationError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsValidationError()
	}
	return false
}

// IsConflictError returns true if the error is a 409 Conflict error
func IsConflictError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsConflict()
	}
	return false
}

// IsServerError returns true if the error is a 5xx server error
func IsServerError(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.IsServerError()
	}
	return false
}
