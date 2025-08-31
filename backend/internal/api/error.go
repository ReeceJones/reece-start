package api

// Custom type of error
type ApiError struct {
	Message string              `json:"message"`
}

func (e ApiError) Error() string {
	return e.Message
}
