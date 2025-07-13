package httpcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds the configuration for HTTP checking
type Config struct {
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`
	Method  string        `json:"method"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		URL:     "http://localhost:8000/check",
		Timeout: 5 * time.Second,
		Method:  "GET",
	}
}

// Checker handles HTTP request checking
type Checker struct {
	config *Config
	client *http.Client
}

// NewChecker creates a new HTTP checker
func NewChecker(config *Config) *Checker {
	if config == nil {
		config = DefaultConfig()
	}
	
	return &Checker{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// CheckConnection makes an HTTP request to verify if connection should be accepted
func (c *Checker) CheckConnection(ctx context.Context, remoteAddr string) error {
	req, err := http.NewRequestWithContext(ctx, c.config.Method, c.config.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add remote address as header for the FastAPI app to use
	req.Header.Set("X-Remote-Addr", remoteAddr)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check if the response indicates the connection should be accepted
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection rejected by FastAPI app, status: %d", resp.StatusCode)
	}
	
	return nil
} 