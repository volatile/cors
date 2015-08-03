package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/volatile/core"
)

const (
	// AllOrigins is the wildcard.
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

// OriginsMap is a list of allowed origins and their respective options.
type OriginsMap map[string]*Options
type formattedOriginsMap map[string]*formattedOptions

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
func Use(origins OriginsMap) {
	fmtOrigins := formatCORS(origins)
	core.Use(func(c *core.Context) {
		setCORS(c, fmtOrigins, c.Next)
	})
}

// LocalUse sets CORS with the provided options locally, for a single handler.
// The global options (if used) are overwritten in this situation.
func LocalUse(c *core.Context, origins OriginsMap, handler func()) {
	setCORS(c, formatCORS(origins), handler)
}

func formatCORS(origins OriginsMap) formattedOriginsMap {
	// If no origins map is set, all are allowed.
	if origins == nil || len(origins) == 0 {
		origins = OriginsMap{AllOrigins: nil}
	}

	fmtOrigins := make(formattedOriginsMap)

	for origin, opts := range origins {
		// If no options are set, all are allowed.
		if opts == nil {
			opts = new(Options)
		}

		fmtOpts := new(formattedOptions)

		if len(opts.AllowedHeaders) > 0 {
			*fmtOpts.AllowedHeaders = strings.Join(opts.AllowedHeaders, ", ")
		}
		if len(opts.AllowedMethods) > 0 {
			*fmtOpts.AllowedMethods = strings.Join(opts.AllowedMethods, ", ")
		}
		if opts.CredentialsAllowed {
			*fmtOpts.CredentialsAllowed = "true"
		}
		if len(opts.ExposedHeaders) > 0 {
			*fmtOpts.ExposedHeaders = strings.Join(opts.ExposedHeaders, ", ")
		}
		if opts.MaxAge.Seconds() != 0 {
			*fmtOpts.MaxAge = fmt.Sprintf("%.f", opts.MaxAge.Seconds())
		}

		fmtOrigins[origin] = fmtOpts
	}

	return fmtOrigins
}

// setCORS sets the response headers and continues downstream if it's not a preflight request.
func setCORS(c *core.Context, fmtOrigins formattedOriginsMap, handler func()) {
	origin := c.Request.Header.Get("Origin")

	// Use CORS only if an Origin header is defined for the request.
	if origin != "" {
		fmtOpts, knownOrigin := fmtOrigins[origin]
		allOriginsAllowed := false

		// If origin is unknown, see for wildcard.
		if !knownOrigin {
			fmtOpts, allOriginsAllowed = fmtOrigins[AllOrigins]
		}

		// If origin is unknown and wildcard isn't set, reject the request.
		if !knownOrigin && !allOriginsAllowed {
			http.Error(c.ResponseWriter, "Invalid CORS request", http.StatusForbidden)
			return
		}

		c.ResponseWriter.Header().Set(headerAllowOrigin, origin)
		c.ResponseWriter.Header().Set("Vary", "Origin")

		// Set credentials header only if they are allowed.
		if fmtOpts.CredentialsAllowed != nil {
			c.ResponseWriter.Header().Set(headerAllowCredentials, *fmtOpts.CredentialsAllowed)
		}

		if fmtOpts.ExposedHeaders != nil {
			c.ResponseWriter.Header().Set(headerExposeHeaders, *fmtOpts.ExposedHeaders)
		}

		if fmtOpts.MaxAge != nil {
			c.ResponseWriter.Header().Set(headerMaxAge, *fmtOpts.MaxAge)
		}

		// OPTIONS method is used for a preflight request.
		// In this case, other CORS headers still need to be set before sending all of them, without any other work downstream.
		if c.Request.Method == "OPTIONS" {
			// If no allowed headers are set, all are allowed.
			if fmtOpts.AllowedHeaders == nil {
				c.ResponseWriter.Header().Set(headerAllowHeaders, c.Request.Header.Get(headerRequestHeaders))
			} else {
				c.ResponseWriter.Header().Set(headerAllowHeaders, *fmtOpts.AllowedHeaders)
			}

			// If no allowed methods are set, all are allowed.
			if fmtOpts.AllowedMethods == nil {
				c.ResponseWriter.Header().Set(headerAllowMethods, c.Request.Header.Get(headerRequestMethod))
			} else {
				c.ResponseWriter.Header().Set(headerAllowMethods, *fmtOpts.AllowedMethods)
			}

			c.ResponseWriter.WriteHeader(http.StatusOK)
			return
		}
	}

	handler()
}
