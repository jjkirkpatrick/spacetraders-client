package models

import (
	"fmt"
)

// APIError represents an error returned by the SpaceTraders API
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error returns the string representation of the APIError
func (e APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// IsAPIError checks if an error is of type APIError
func IsAPIError(err error) bool {
	_, ok := err.(APIError)
	return ok
}
