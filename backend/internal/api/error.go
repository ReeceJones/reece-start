package api

import (
	"fmt"

	"reece.start/internal/constants"
)

// Custom type of error
type ApiError struct {
	Code    constants.ErrorCode `json:"code"`
	Message string              `json:"message"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Error: %s, Message: %s", e.Code, e.Message)
}
