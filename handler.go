package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/volatile/core"
)

const (
	// AllOrigins contains the wildcard symbol.
	AllOrigins = "*"

	// Response headers
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"

	// Request headers
	headerRequestHeaders = "Access-Control-Request-Headers"
	headerRequestMethod  = "Access-Control-Request-Method"
)

// Options represents access control options for an origin.
type Options struct {
	AllowedHeaders     []string      // AllowedHeaders indicates, in the case of a preflight request, which headers can be used during the actual request. If none are set, all are allowed.
	AllowedMethods     []string      // AllowedMethods indicates, in the case of a preflight request, which methods can be used during the actual request. If none are set, all are allowed.
	CredentialsAllowed bool          // CredentialsAllowed indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
	ExposedHeaders     []string      // ExposedHeaders whitelists headers that browsers are allowed to access.
	MaxAge             time.Duration // MaxAge indicates how long the results of a preflight request can be cached.
}

type formattedOptions struct {
	AllowedHeaders     *string
	AllowedMethods     *string
	CredentialsAllowed *string
	ExposedHeaders     *string
	MaxAge             *string
}

// Use adds a handler that sets CORS with the provided options for all dowstream requests.
func Use(options map[string]Options) {
	fmtOpt := formatCORS(options)
	core.Use(func(c *core.Context) {
		setCORS(c, fmtOpt, c.Next)
	})
}

// LocalUse sets CORS with the provided options locally, for a single handler.
// The global options (if used) are overwritten in this situation.
func LocalUse(c *core.Context, options map[string]Options, handler func()) {
	setCORS(c, formatCORS(options), handler)
}

func formatCORS(opt map[string]Options) map[string]formattedOptions {
	fmtOpt := make(map[string]formattedOptions, len(opt))

	for origin, item := range opt {
		result := formattedOptions{}
		if len(item.AllowedHeaders) > 0 {
			*result.AllowedHeaders = strings.Join(item.AllowedHeaders, ", ")
		}
		if len(item.AllowedMethods) > 0 {
			*result.AllowedMethods = strings.Join(item.AllowedMethods, ", ")
		}
		if item.CredentialsAllowed {
			*result.CredentialsAllowed = "true"
		}
		if len(item.ExposedHeaders) > 0 {
			*result.ExposedHeaders = strings.Join(item.ExposedHeaders, ", ")
		}
		if item.MaxAge.Seconds() > 0.5 {
			*result.MaxAge = fmt.Sprintf("%.f", item.MaxAge.Seconds())
		}

		fmtOpt[origin] = result
	}

	return fmtOpt
}

// setCORS sets the response headers and continues downstream if it's not a preflight request.
func setCORS(c *core.Context, fmtOpts map[string]formattedOptions, handler func()) {
	origin := c.Request.Header.Get("Origin")

	// Use CORS only if an Origin header is defined for the request.
	if origin != "" {
		fmtOpt, knownOrigin := fmtOpts[origin]
		allOriginsAllowed := false

		// Unknown origin: check for wildcard.
		if !knownOrigin {
			fmtOpt, allOriginsAllowed = fmtOpts[AllOrigins]
		}

		// No origin matched and wildcard not accepted: reject the request.
		if !knownOrigin && !allOriginsAllowed {
			http.Error(c.ResponseWriter, "Invalid CORS request", http.StatusForbidden)
			return
		}

		c.ResponseWriter.Header().Set(headerAllowOrigin, origin)
		c.ResponseWriter.Header().Set("Vary", "Origin")

		// Set credentials header only if they are allowed.
		if fmtOpt.CredentialsAllowed != nil {
			c.ResponseWriter.Header().Set(headerAllowCredentials, *fmtOpt.CredentialsAllowed)
		}

		if fmtOpt.ExposedHeaders != nil {
			c.ResponseWriter.Header().Set(headerExposeHeaders, *fmtOpt.ExposedHeaders)
		}

		if fmtOpt.MaxAge != nil {
			c.ResponseWriter.Header().Set(headerMaxAge, *fmtOpt.MaxAge)
		}

		// OPTIONS method is used for a preflight request.
		// In this case, other CORS headers still need to be set before sending all of them, without any other work.
		if c.Request.Method == "OPTIONS" {
			// If no allowed headers are set, all are allowed.
			if fmtOpt.AllowedHeaders == nil {
				c.ResponseWriter.Header().Set(headerAllowHeaders, c.Request.Header.Get(headerRequestHeaders))
			} else {
				c.ResponseWriter.Header().Set(headerAllowHeaders, *fmtOpt.AllowedHeaders)
			}

			// If no allowed methods are set, all are allowed.
			if fmtOpt.AllowedMethods == nil {
				c.ResponseWriter.Header().Set(headerAllowMethods, c.Request.Header.Get(headerRequestMethod))
			} else {
				c.ResponseWriter.Header().Set(headerAllowMethods, *fmtOpt.AllowedMethods)
			}

			c.ResponseWriter.WriteHeader(http.StatusOK)
			return
		}
	}

	handler()
}
