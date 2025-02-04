package models

import (
	"fmt"
	"strings"
)

// APIError represents an error returned by the SpaceTraders API
type APIError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// HasValidationErrors checks if there are any validation errors in the Data field
func (e *APIError) HasValidationErrors() bool {
	return e.Data != nil && len(e.Data) > 0
}

// GetFieldErrors returns validation errors for a specific field
func (e *APIError) GetFieldErrors(field string) []string {
	if e.Data == nil {
		return nil
	}

	if fieldErrors, exists := e.Data[field].([]interface{}); exists {
		errors := make([]string, len(fieldErrors))
		for i, err := range fieldErrors {
			errors[i] = fmt.Sprint(err)
		}
		return errors
	}
	return nil
}

// GetValidationFields returns all field names that have validation errors
func (e *APIError) GetValidationFields() []string {
	if e.Data == nil {
		return nil
	}

	var fields []string
	for field := range e.Data {
		fields = append(fields, field)
	}
	return fields
}

// FormattedMessage returns a formatted error message including validation errors
func (e *APIError) FormattedMessage() string {
	if !e.HasValidationErrors() {
		return fmt.Sprintf("[%d] %s", e.Code, e.Message)
	}

	var sb strings.Builder
	// Add the main error message with code
	fmt.Fprintf(&sb, "[%d] %s", e.Code, e.Message)

	// Add each validation error on its own line with proper indentation
	for field, messages := range e.Data {
		if msgArray, ok := messages.([]interface{}); ok {
			for _, msg := range msgArray {
				fmt.Fprintf(&sb, " | %s: %s", field, msg)
			}
		}
	}
	return sb.String()
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.FormattedMessage()
}

// AsError returns the APIError as a standard error interface
func (e APIError) AsError() error {
	return &e
}

// IsAPIError checks if an error is of type APIError
func IsAPIError(err error) bool {
	_, ok := err.(APIError)
	return ok
}
