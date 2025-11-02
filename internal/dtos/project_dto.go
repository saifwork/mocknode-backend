package dtos

// DTO for project creation (called internally after session creation)
type CreateProjectRequest struct {
	Name string `json:"name,omitempty"` // Optional name; can default to "My Project"
}

// Response when a project is created or fetched
type CreateProjectResponse struct {
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}

// Metadata stored in Redis for project details
type ProjectMeta struct {
	ProjectID  string `json:"projectId"`
	Name       string `json:"name"`
	SessionID  string `json:"sessionId"` // backlink to session
	CreatedAt  string `json:"createdAt"`
	Collection int    `json:"collectionCount"`
	ReqCount   int    `json:"reqCount"`
}
