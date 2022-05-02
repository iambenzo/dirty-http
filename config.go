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
	Enabled            bool
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

	c.Cors.Enabled = viper.GetBool("middleware.cors.enabled")
	c.Cors.CorsAllowedOrigins = make(map[string]bool)

	if c.Cors.Enabled {
		for _, v := range viper.GetStringSlice("middleware.cors.allowed-origins") {
			if v == "*" {
				c.Cors.CorsAllowedOrigins = make(map[string]bool)
				c.Cors.CorsAllowedOrigins[v] = true
				break
			}
			c.Cors.CorsAllowedOrigins[v] = true
		}

		c.Cors.CorsAllowedMethods = make(map[string]bool)
		c.Cors.CorsAllowedMethods[http.MethodOptions] = true
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
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Use = os.Args[0]

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is ./.config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&printConfig, "print-config", false, "Print application configuration on startup")
	rootCmd.Flags().IntP("port", "p", 8080, "Port to run Application server on")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	rootCmd.Flags().Bool("cors", false, "Enable CORS processing")
	viper.BindPFlag("middleware.cors.enabled", rootCmd.Flags().Lookup("cors"))
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

	viper.SetDefault("middleware.cors.enabled", false)
	viper.SetDefault("middleware.cors.allowed-origins", []string{"*"})
	viper.SetDefault("middleware.cors.allowed-methods", []string{http.MethodOptions, http.MethodGet})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in local directory with name ".config" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
