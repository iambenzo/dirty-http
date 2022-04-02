package middleware

import "net/http"

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func DefaultCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
}

func (c *CorsConfig) setAllowOrigins(origins []string) {
	c.AllowOrigins = origins
}

func (c *CorsConfig) setAllowMethods(methods []string) {
	c.AllowMethods = methods
}

func (c *CorsConfig) setAllowHeaders(headers []string) {
	c.AllowHeaders = headers
}

func (c *CorsConfig) setExposeHeaders(headers []string) {
	c.ExposeHeaders = headers
}

func (c *CorsConfig) setAllowCredentials(allow bool) {
	c.AllowCredentials = allow
}

func (c *CorsConfig) setMaxAge(age int) {
	c.MaxAge = age
}
