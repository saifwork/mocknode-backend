package redis

import (
	"fmt"
)

// ---------------------------
// Redis Key Naming Patterns
// ---------------------------

// Each session maps to one project
func SessionMetaKey(sessionID string) string {
	return fmt.Sprintf("session:%s:meta", sessionID)
}

func SessionProjectKey(sessionID string) string {
	return fmt.Sprintf("session:%s:project", sessionID)
}

// Project-level keys
func ProjectMetaKey(projectID string) string {
	return fmt.Sprintf("project:%s:meta", projectID)
}

func ProjectCollectionsKey(projectID string) string {
	return fmt.Sprintf("project:%s:collections", projectID)
}

func ProjectDataKey(projectID, collection string) string {
	return fmt.Sprintf("project:%s:data:%s", projectID, collection)
}

// Track requests and usage
func ProjectRequestCountKey(projectID string) string {
	return fmt.Sprintf("project:%s:reqcount", projectID)
}

func ProjectRequestsLogKey(projectID string) string {
	return fmt.Sprintf("project:%s:requests", projectID)
}

// Optional reverse mapping to quickly find session by project
func ProjectSessionKey(projectID string) string {
	return fmt.Sprintf("project:%s:session", projectID)
}

func SessionCollectionsKey(sessionID string) string {
	return fmt.Sprintf("session:%s:collections", sessionID)
}

func SessionDataKey(sessionID, collection string) string {
	return fmt.Sprintf("session:%s:data:%s", sessionID, collection)
}

func SessionRequestCountKey(sessionID string) string {
	return fmt.Sprintf("session:%s:reqcount", sessionID)
}

func SessionRequestsLogKey(sessionID string) string {
	return fmt.Sprintf("session:%s:requests", sessionID)
}
