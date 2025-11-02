package dtos

import "time"

// Request DTO to create a new session
type CreateSessionRequest struct {
	ExpirySeconds int `json:"expirySeconds,omitempty"` // Optional (default 86400)
}

// Response DTO when a session is created
type CreateSessionResponse struct {
	SessionID     string `json:"sessionId"`
	ProjectID     string `json:"projectId"` // 🔹 new: link session → project
	ExpirySeconds int    `json:"expirySeconds"`
	Message       string `json:"message"`
}

// Metadata stored in Redis for each session
type SessionMeta struct {
	SessionID string        `json:"sessionId"`
	ProjectID string        `json:"projectId"` // 🔹 new: associated project
	CreatedAt string        `json:"createdAt"`
	Expiry    time.Duration `json:"expiry"`   // in seconds
	ReqCount  int           `json:"reqCount"` // total requests so far
}
