package harness

import (
	"os"
)

// Config holds the configuration for the Harness API client
type Config struct {
	BaseURL     string
	CodeBaseURL string
	AccountID   string
	OrgID       string
	ProjectID   string
}

// NewConfig creates a new Config with default values that can be overridden by environment variables
func NewConfig() *Config {
	// Default values
	defaultBaseURL := "https://harness0.harness.io/ng/api"
	defaultCodeBaseURL := "https://harness0.harness.io/code/api"
	defaultAccountID := "l7B_kbSEQD2wjrM7PShm5w"
	defaultOrgID := "PROD"
	defaultProjectID := "Harness_Commons"

	// Environment variables override defaults
	baseURL := getEnv("HARNESS_BASE_URL", defaultBaseURL)
	codeBaseURL := getEnv("HARNESS_CODE_BASE_URL", defaultCodeBaseURL)
	accountID := getEnv("HARNESS_ACCOUNT_ID", defaultAccountID)
	orgID := getEnv("HARNESS_ORG_ID", defaultOrgID)
	projectID := getEnv("HARNESS_PROJECT_ID", defaultProjectID)

	return &Config{
		BaseURL:     baseURL,
		CodeBaseURL: codeBaseURL,
		AccountID:   accountID,
		OrgID:       orgID,
		ProjectID:   projectID,
	}
}

// getEnv retrieves an environment variable value or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetCodeBaseURL returns the base URL for code-related APIs
func (c *Config) GetCodeBaseURL() string {
	return c.CodeBaseURL
}
