package config

import (
	"github.com/kelseyhightower/envconfig"
)

// SERVICENAME is an environment variables prefix
const SERVICENAME = "CARTAPI"

// DatabaseConfig contains variables, that are required for a database connection
type DatabaseConfig struct {
	DBName           string `split_words:"true" required:"true"`
	ConnectionString string `split_words:"true" required:"true"`
}

// Load settles environment variables into AppConfig structure
func (c *DatabaseConfig) Load(serviceName string) error {
	return envconfig.Process(serviceName, c)
}
