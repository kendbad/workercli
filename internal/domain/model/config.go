package model

// Config holds business-level configuration for TUI display
type Config struct {
	MaxTasks     int // Maximum number of concurrent tasks
	ProxyTimeout int // Timeout for proxy checks (in seconds)
}
