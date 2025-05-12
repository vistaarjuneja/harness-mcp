package config

type Config struct {
	Version     string
	BaseURL     string
	AccountID   string
	OrgID       string
	ProjectID   string
	APIKey      string
	ReadOnly    bool
	Toolsets    []string
	LogFilePath string
	Debug       bool
}
