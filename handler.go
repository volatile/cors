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
)

// OriginsMap represents the allowed origins with their respective options.
type OriginsMap map[string]*Options

// Options represents access control options for an origin.
type Options struct {
	AllowedHeaders     []string      // AllowedHeaders indicates, in the case of a preflight request, which headers can be used during the actual request. If none are set, all are allowed.
	AllowedMethods     []string      // AllowedMethods indicates, in the case of a preflight request, which methods can be used during the actual request. If none are set, all are allowed.
	CredentialsAllowed bool          // CredentialsAllowed indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
	ExposedHeaders     []string      // ExposedHeaders whitelists headers that browsers are allowed to access.
	MaxAge             time.Duration // MaxAge indicates how long the results of a preflight request can be cached.
}

// Use adds a handler that sets CORS with the provided options for all handlers dowstream.
func Use(origins *OriginsMap) {
	core.Use(func(c *core.Context) {
		setCORS(c, origins, c.Next)
	})
}

// LocalUse sets CORS with the provided options locally, for a single handler.
// The global options (if used) are overwritten in this situation.
func LocalUse(c *core.Context, origins *OriginsMap, handler func()) {
	setCORS(c, origins, handler)
}

// setCORS sets the response headers and continues downstream if it's not a preflight request.
func setCORS(c *core.Context, origins *OriginsMap, handler func()) {
	origin := c.Request.Header.Get("Origin")

	// Don't use CORS without an origin.
	if origin == "" {
		handler()
		return
	}

	if origins == nil || len(*origins) == 0 {
		origins = &OriginsMap{AllOrigins: nil}
	}

	opts, knownOrigin := (*origins)[origin]

	// If origin is unknown, see for wildcard.
	var allOriginsAllowed bool
	if !knownOrigin {
		opts, allOriginsAllowed = (*origins)[AllOrigins]
	}

	// If origin is unknown and wildcard isn't set, reject the request.
	if !knownOrigin && !allOriginsAllowed {
		http.Error(c.ResponseWriter, "Invalid CORS request", http.StatusForbidden)
		return
	}

	c.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)
	c.ResponseWriter.Header().Set("Vary", "Origin")

	// Set credentials header only if they are allowed.
	if opts != nil && opts.CredentialsAllowed {
		c.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if opts != nil && len(opts.ExposedHeaders) > 0 {
		c.ResponseWriter.Header().Set("Access-Control-Expose-Headers", strings.Join(opts.ExposedHeaders, ", "))
	}

	if opts != nil && opts.MaxAge != 0 {
		c.ResponseWriter.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%.f", opts.MaxAge.Seconds()))
	}

	// OPTIONS method is used for a preflight request.
	// In this case, other CORS headers still need to be set before sending all of them, without any other work downstream.
	if c.Request.Method != "OPTIONS" {
		handler()
		return
	}

	// If no allowed headers are set, all are allowed.
	if opts != nil && len(opts.AllowedHeaders) > 0 {
		c.ResponseWriter.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowedHeaders, ", "))
	} else {
		c.ResponseWriter.Header().Set("Access-Control-Allow-Headers", c.Request.Header.Get("Access-Control-Request-Headers"))
	}

	// If no allowed methods are set, all are allowed.
	if opts != nil && len(opts.AllowedHeaders) > 0 {
		c.ResponseWriter.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ", "))
	} else {
		c.ResponseWriter.Header().Set("Access-Control-Allow-Methods", c.Request.Header.Get("Access-Control-Request-Method"))
	}

	// It was a preflight request so we just send the headers.
	c.ResponseWriter.WriteHeader(http.StatusOK)
}
