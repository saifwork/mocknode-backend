package dtos

// Generic model to hold user’s mock object
type MockData struct {
	ID      string `json:"id"`
	Payload any    `json:"payload"`
}

// DTO for create/update requests
type MockCreateRequest struct {
	Data map[string]any `json:"data" binding:"required"`
}

// DTO for list or single data response
type MockDataResponse struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
}

// Response when performing any CRUD operation
type MockActionResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
