// Package config loads the application's configuration values from environment variables.
// These values are used to configure the port, database connection, and session settings.
package config

import "os"

// EnvData contains the runtime values that are read from the environment.
type EnvData struct {
	Port          string
	DB            string
	SessionSecret string
}

// Load reads the configured environment variables and returns them as a structured config object.
func Load() EnvData {
	return EnvData{
		Port:          os.Getenv("PORT"),
		DB:   	         os.Getenv("DATABASE_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
	}
}
