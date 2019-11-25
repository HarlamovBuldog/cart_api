package config

import (
	"github.com/kelseyhightower/envconfig"
)

// SERVICENAME is an environment variables prefix
const SERVICENAME = "CMGWAREHOUSE"

// AppConfig contains all of the config variables for the application
type AppConfig struct {
	DbConfig              DatabaseConfig
	TokenSigningSecret    string `split_words:"true" required:"true"`
	PasswordSigningSecret string `split_words:"true" required:"true"`
	IpInfoToken           string `split_words:"true" required:"true"`
	AmazonRegion          string `split_words:"true" required:"true"`
	AmazonEmailRegion     string `split_words:"true" required:"true"`
	ReportsEmail          string `split_words:"true" required:"true"`
}

// DatabaseConfig contains variables, that are required for a database connection
type DatabaseConfig struct {
	Dialect          string `split_words:"true" required:"true"`
	ConnectionString string `split_words:"true" required:"true"`
}

// Load settles environment variables into AppConfig structure
func (c *AppConfig) Load(serviceName string) error {
	return envconfig.Process(serviceName, c)
}
