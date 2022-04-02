package dirtyhttp

import (
	"errors"
	"os"
	"strconv"

	"github.com/iambenzo/dirtyhttp/middleware"
)

type EnvConfig struct {
	ApiUser     string
	ApiPassword string
	ApiPort     string
	DbUrl       string
	DbUser      string
	DbPassword  string
	DbName      string
	Cors        middleware.CorsConfig
}

// Pulls configuration from matching environment variables.
//
// For example, ApiUser pulls from API_USER environment variable
// and DbPassword pulls from DB_PASSWORD.
//
// At the moment the API_USER and API_PASSWORD environment variables
// are required.
//
// The DB related variables are optional.
func getEnvConfig() (*EnvConfig, error) {
	var haveProblem = false
	var apiPort string

	if os.Getenv("API_USER") == "" {
		haveProblem = true
	}

	if os.Getenv("API_PASSWORD") == "" {
		haveProblem = true
	}

	if os.Getenv("API_PORT") == "" {
		apiPort = ":8080"
	} else {
		_, err := strconv.Atoi(os.Getenv("API_PORT"))
		if err != nil {
			apiPort = ":8080"
		} else {
			apiPort = ":" + os.Getenv("API_PORT")
		}
	}

	if os.Getenv("DB_URL") != "" {
		if os.Getenv("DB_USER") == "" {
			haveProblem = true
		}
		if os.Getenv("DB_PASSWORD") == "" {
			haveProblem = true
		}
		if os.Getenv("DB_NAME") == "" {
			haveProblem = true
		}
	}

	if haveProblem {
		return &EnvConfig{}, errors.New("not all environment variables are set")
	} else {
		return &EnvConfig{
			ApiUser:     os.Getenv("API_USER"),
			ApiPassword: os.Getenv("API_PASSWORD"),
			ApiPort:     apiPort,
			DbUrl:       os.Getenv("DB_URL"),
			DbUser:      os.Getenv("DB_USER"),
			DbPassword:  os.Getenv("DB_PASSWORD"),
			DbName:      os.Getenv("DB_NAME"),
			Cors:        *middleware.DefaultCorsConfig(),
		}, nil
	}
}
