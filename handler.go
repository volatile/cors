package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/volatile/core"
)

const (
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"

	headerRequestHeaders = "Access-Control-Request-Headers"
	headerRequestMethod  = "Access-Control-Request-Method"
)

// Options represents access control options.
type Options struct {
	// AllowedHeaders indicates, in the case of a preflight request, which headers can be used during the actual request.
	AllowedHeaders []string
	// AllowedMethods indicates, in the case of a preflight request, which methods can be used during the actual request.
	AllowedMethods []string
	// AllowedOrigins indicates which origins can make a request.
	AllowedOrigins []string
	// CredentialsAllowed indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
	CredentialsAllowed bool
	// ExposedHeaders whitelists headers that browsers are allowed to access.
	ExposedHeaders []string
	// MaxAge indicates how long the results of a preflight request can be cached.
	MaxAge time.Duration
}

// Use tells the core to use this handler with the provided options.
func Use(options *Options) {
	if options == nil {
		options = new(Options)
	}

	// First, format all the potential data to not do that on each request.

	headerAllowHeadersValue := strings.Join(options.AllowedHeaders, ", ")

	headerAllowMethodsValue := strings.Join(options.AllowedMethods, ", ")

	var headerAllowOriginValue string
	// If no allowed origins set, all origins are allowed.
	if options.AllowedOrigins == nil {
		headerAllowOriginValue = "*"
	} else {
		headerAllowOriginValue = strings.Join(options.AllowedOrigins, ", ")
	}

	headerExposeHeadersValue := strings.Join(options.ExposedHeaders, ", ")

	headerMaxAgeValue := fmt.Sprintf("%.f", options.MaxAge.Seconds())

	// Finally, add handler to stack.

	core.Use(func(c *core.Context) {
		c.ResponseWriter.Header().Set(headerAllowOrigin, headerAllowOriginValue)

		if options.CredentialsAllowed {
			c.ResponseWriter.Header().Set(headerAllowCredentials, "true")
		}

		if headerExposeHeadersValue != "" {
			c.ResponseWriter.Header().Set(headerExposeHeaders, headerExposeHeadersValue)
		}

		if options.MaxAge != 0 {
			c.ResponseWriter.Header().Set(headerMaxAge, headerMaxAgeValue)
		}

		// OPTIONS is used for preflight requests.
		// For that, only the CORS handler must respond, so the handlers chain is broken.
		if c.Request.Method == "OPTIONS" {
			// If no allowed headers are set, accept all from the real request.
			if headerAllowHeadersValue == "" {
				c.ResponseWriter.Header().Set(headerAllowHeaders, c.Request.Header.Get(headerRequestHeaders))
			} else {
				c.ResponseWriter.Header().Set(headerAllowHeaders, headerAllowHeadersValue)
			}

			// If no allowed methods are set, accept the method of the real request.
			if headerAllowMethodsValue == "" {
				c.ResponseWriter.Header().Set(headerAllowMethods, c.Request.Header.Get(headerRequestMethod))
			} else {
				c.ResponseWriter.Header().Set(headerAllowMethods, headerAllowMethodsValue)
			}

			c.ResponseWriter.WriteHeader(http.StatusOK)
		} else {
			c.Next()
		}
	})
}
