// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "os"

// IConfig is interface for application config
type IConfig interface {
	ServiceID() string
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
