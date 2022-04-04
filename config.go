package dirtyhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Database relevant config
// Include port in DbUrl
type DbConfig struct {
	DbUrl      string
	DbUser     string
	DbPassword string
	DbName     string
}

// Cors relevant configuration
type CorsConfig struct {
	CorsAllowedOrigins map[string]bool
	CorsAllowedMethods map[string]bool
}

// User and Password for Basic Authentication
type AuthConfig struct {
	Enabled     bool
	ApiUser     string
	ApiPassword string
}

// Timeout configuration
type TimeoutConfig struct {
	Enabled bool
	Length  int
}

// Gzip configuration
type GzipConfig struct {
	Enabled bool
}

// Struct for all configuration
type Config struct {
	ApiPort        string
	Authentication AuthConfig
	Timeout        TimeoutConfig
	Gzip           GzipConfig
	Database       DbConfig
	Cors           CorsConfig
}

func newConfig() (*Config, error) {
	var c Config

	c.ApiPort = ":" + viper.GetString("port")

	c.Authentication.Enabled = viper.GetBool("middleware.auth.enabled")
	c.Authentication.ApiUser = viper.GetString("middleware.auth.user")
	c.Authentication.ApiPassword = viper.GetString("middleware.auth.password")

	c.Timeout.Enabled = viper.GetBool("middleware.timeout.enabled")
	c.Timeout.Length = viper.GetInt("middleware.timeout.length")

	c.Database.DbUrl = viper.GetString("db.url")
	c.Database.DbName = viper.GetString("db.name")
	c.Database.DbUser = viper.GetString("db.user")
	c.Database.DbUrl = viper.GetString("db.url")

	c.Cors.CorsAllowedOrigins = make(map[string]bool)
	for _, v := range viper.GetStringSlice("middleware.cors.allowed-origins") {
		if v == "*" {
			c.Cors.CorsAllowedOrigins = make(map[string]bool)
			c.Cors.CorsAllowedOrigins[v] = true
			break
		}
		c.Cors.CorsAllowedOrigins[v] = true
	}

	c.Cors.CorsAllowedMethods = make(map[string]bool)
	for _, v := range viper.GetStringSlice("middleware.cors.allowed-methods") {
		switch strings.ToUpper(v) {
		case http.MethodGet:
			c.Cors.CorsAllowedMethods[http.MethodGet] = true
		case http.MethodPut:
			c.Cors.CorsAllowedMethods[http.MethodPut] = true
		case http.MethodPost:
			c.Cors.CorsAllowedMethods[http.MethodPost] = true
		case http.MethodDelete:
			c.Cors.CorsAllowedMethods[http.MethodDelete] = true
		}
	}

	return &c, nil

}

var cfgFile string
var printConfig bool
var config *Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(_ *cobra.Command, _ []string) {
		if printConfig {
			prettyPrint(viper.AllSettings())
		}

		c, err := newConfig()
		if err != nil {
			os.Exit(1)
		}

		config = c
	},
}

func getConfig() {
	err := rootCmd.Execute()
	if err != nil {
		l := logger{}
		l.Fatal(fmt.Sprintf("Error in config: %v", err))
		// os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Use = os.Args[0]

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is ./.config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&printConfig, "print-config", false, "Print application configuration on startup")
	rootCmd.Flags().IntP("port", "p", 8080, "Port to run Application server on")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// set default values
	viper.SetDefault("port", 8080)

	viper.SetDefault("middleware.gzip.enabled", true)

	viper.SetDefault("middleware.auth.enabled", false)
	viper.SetDefault("middleware.auth.user", "")
	viper.SetDefault("middleware.auth.password", "")

	viper.BindEnv("middleware.auth.user", "API_USER")
	viper.BindEnv("middleware.auth.password", "API_PASSWORD")

	viper.SetDefault("middleware.timeout.enabled", true)
	viper.SetDefault("middleware.timeout.length", 30)

	viper.SetDefault("db.url", "")
	viper.SetDefault("db.name", "")
	viper.SetDefault("db.user", "")
	viper.SetDefault("db.password", "")

	viper.BindEnv("db.url", "DB_URL")
	viper.BindEnv("db.name", "DB_NAME")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")

	viper.SetDefault("middleware.cors.allowed-origins", []string{"*"})
	viper.SetDefault("middleware.cors.allowed-methods", []string{"GET"})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".config-play" (without extension).
		viper.AddConfigPath(".")
		// viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.SetEnvPrefix("API")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

//import (
//	"errors"
//	"os"
//	"strconv"

//	"github.com/iambenzo/dirtyhttp/middleware"
//)

//type EnvConfig struct {
//	ApiUser     string
//	ApiPassword string
//	ApiPort     string
//	DbUrl       string
//	DbUser      string
//	DbPassword  string
//	DbName      string
//	Cors        middleware.CorsConfig
//}

//// Pulls configuration from matching environment variables.
////
//// For example, ApiUser pulls from API_USER environment variable
//// and DbPassword pulls from DB_PASSWORD.
////
//// At the moment the API_USER and API_PASSWORD environment variables
//// are required.
////
//// The DB related variables are optional.
//func getEnvConfig() (*EnvConfig, error) {
//	var haveProblem = false
//	var apiPort string

//	if os.Getenv("API_USER") == "" {
//		haveProblem = true
//	}

//	if os.Getenv("API_PASSWORD") == "" {
//		haveProblem = true
//	}

//	if os.Getenv("API_PORT") == "" {
//		apiPort = ":8080"
//	} else {
//		_, err := strconv.Atoi(os.Getenv("API_PORT"))
//		if err != nil {
//			apiPort = ":8080"
//		} else {
//			apiPort = ":" + os.Getenv("API_PORT")
//		}
//	}

//	if os.Getenv("DB_URL") != "" {
//		if os.Getenv("DB_USER") == "" {
//			haveProblem = true
//		}
//		if os.Getenv("DB_PASSWORD") == "" {
//			haveProblem = true
//		}
//		if os.Getenv("DB_NAME") == "" {
//			haveProblem = true
//		}
//	}

//	if haveProblem {
//		return &EnvConfig{}, errors.New("not all environment variables are set")
//	} else {
//		return &EnvConfig{
//			ApiUser:     os.Getenv("API_USER"),
//			ApiPassword: os.Getenv("API_PASSWORD"),
//			ApiPort:     apiPort,
//			DbUrl:       os.Getenv("DB_URL"),
//			DbUser:      os.Getenv("DB_USER"),
//			DbPassword:  os.Getenv("DB_PASSWORD"),
//			DbName:      os.Getenv("DB_NAME"),
//			Cors:        *middleware.DefaultCorsConfig(),
//		}, nil
//	}
//}
