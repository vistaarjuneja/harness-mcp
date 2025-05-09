package dto

// Scope represents a scope in the system
type Scope struct {
	AccountID string `json:"accountIdentifier"`
	OrgID     string `json:"orgIdentifier"`
	ProjectID string `json:"projectIdentifier"`
}
