package cors

import (
	"fmt"
	"net/http"
	"strconv"
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

// Options represents access control options for an origin.
type Options struct {
	// AllowedHeaders indicates, in the case of a preflight request,
	// which headers can be used during the actual request.
	AllowedHeaders []string

	// AllowedMethods indicates, in the case of a preflight request,
	// which methods can be used during the actual request.
	//
	// If len(ExposedHeaders) == 0, all methods will be allowed.
	AllowedMethods []string

	// CredentialsAllowed indicates whether the request can include user
	// credentials like cookies, HTTP authentication or client side SSL certificates.
	CredentialsAllowed bool

	// ExposedHeaders whitelists headers that browsers are allowed to access.
	//
	// If len(ExposedHeaders) == 0, all headers on preflight requests will be exposed.
	ExposedHeaders []string

	// MaxAge indicates how long the results of a preflight request can be cached.
	MaxAge time.Duration
}

type formattedOptions struct {
	AllowedHeaders,
	AllowedMethods,
	CredentialsAllowed,
	ExposedHeaders,
	MaxAge *string
}

// Use tells the core to use this handler with the provided options.
func Use(options map[string]Options) {
	fmtOpt := formatCORS(options)
	core.Use(func(c *core.Context) {
		setCORS(c, fmtOpt, c.Next)
	})
}

// LocalUse allows to set CORS locally, for a single handler.
// Remember you can't set headers after Write or WriteHeader has been called.
func LocalUse(c *core.Context, options map[string]Options, handler func()) {
	setCORS(c, formatCORS(options), handler)
}

func formatCORS(opt map[string]Options) (fmtOpt map[string]formattedOptions) {
	// For requests without credentials, the server may specify "*" as a wildcard,
	// thereby allowing any origin to access the resource.
	if wild, ok := opt["*"]; ok && wild.CredentialsAllowed {
		panic("sending credentials via CORS is not permitted for wildcarded origins")
	}

	fmtOpt = make(map[string]formattedOptions, len(opt))

	for origin, item := range opt {
		result := formattedOptions{}
		if len(item.AllowedHeaders) > 0 {
			*result.AllowedHeaders = strings.Join(item.AllowedHeaders, ", ")
		}
		if len(item.AllowedMethods) > 0 {
			*result.AllowedMethods = strings.Join(item.AllowedMethods, ", ")
		}
		if item.CredentialsAllowed {
			*result.CredentialsAllowed = strconv.FormatBool(item.CredentialsAllowed)
		}
		if len(item.ExposedHeaders) > 0 {
			*result.ExposedHeaders = strings.Join(item.ExposedHeaders, ", ")
		}
		if item.MaxAge.Seconds() > 0.5 {
			*result.MaxAge = fmt.Sprintf("%.f", item.MaxAge.Seconds())
		}

		fmtOpt[origin] = result
	}

	return
}

// setCORS set the response headers and continue if it's not a preflight request.
func setCORS(c *core.Context, fmtOpts map[string]formattedOptions, handler func()) {
	origin := c.Request.Header.Get("Origin")

	// Request does not have an Origin header, therefore it is no CORS request
	// and control is passed to the handler.
	if origin == "" {
		handler()
		return
	}

	fmtOpt, ok := fmtOpts[origin]
	wildcard = false

	// Source origin is not known by name, so check for wildcard.
	if !ok {
		fmtOpt, wildcard = fmtOpts["*"]
	}

	// No origin matched and wildcard not accepted, reject the request.
	if !ok && !wildcard {
		http.Error(c.ResponseWriter, "Invalid CORS request", http.StatusForbidden)
		return
	}

	c.ResponseWriter.Header().Set(headerAllowOrigin, origin)

	if wildcard {
		c.ResponseWriter.Header().Set("Vary", "Origin")
	}

	if fmtOpt.CredentialsAllowed != nil {
		c.ResponseWriter.Header().Set(headerAllowCredentials, *fmtOpt.CredentialsAllowed)
	}

	if fmtOpt.ExposedHeaders != nil {
		c.ResponseWriter.Header().Set(headerExposeHeaders, *fmtOpt.ExposedHeaders)
	}

	if fmtOpt.MaxAge != nil {
		c.ResponseWriter.Header().Set(headerMaxAge, *fmtOpt.MaxAge)
	}

	if c.Request.Method != "OPTIONS" {
		handler()
		return
	}

	// OPTIONS is used for preflight requests.
	// For that, only the CORS handler must respond, so the handlers chain is broken.

	// If no allowed headers are set, accept all from the real request.
	if fmtOpt.AllowedHeaders == nil {
		c.ResponseWriter.Header().Set(headerAllowHeaders, c.Request.Header.Get(headerRequestHeaders))
	} else {
		c.ResponseWriter.Header().Set(headerAllowHeaders, *fmtOpt.AllowedHeaders)
	}

	// If no allowed methods are set, accept the method of the real request.
	if fmtOpt.AllowedMethods == nil {
		c.ResponseWriter.Header().Set(headerAllowMethods, c.Request.Header.Get(headerRequestMethod))
	} else {
		c.ResponseWriter.Header().Set(headerAllowMethods, *fmtOpt.AllowedMethods)
	}

	c.ResponseWriter.WriteHeader(http.StatusOK)
}
