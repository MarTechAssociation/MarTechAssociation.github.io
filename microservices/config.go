// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package microservices

import "os"

// IConfig is interface for application config
type IConfig interface {
	ServiceID() string
	GoogleToken() string
}

// Config implement IConfig
type Config struct{}

// NewConfig return new config instance
func NewConfig() *Config {
	return &Config{}
}

// ServiceID return ID of service
func (cfg *Config) ServiceID() string {
	return os.Getenv("SERVICE_ID")
}

func (cfg *Config) GoogleToken() string {
	return os.Getenv("GOOGLE_TOKEN")
}
