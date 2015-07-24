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

// Options represents access control options.
type Options struct {
	AllowedHeaders     []string      // AllowedHeaders indicates, in the case of a preflight request, which headers can be used during the actual request.
	AllowedMethods     []string      // AllowedMethods indicates, in the case of a preflight request, which methods can be used during the actual request.
	AllowedOrigins     []string      // AllowedOrigins indicates which origins can make a request.
	CredentialsAllowed bool          // CredentialsAllowed indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
	ExposedHeaders     []string      // ExposedHeaders whitelists headers that browsers are allowed to access.
	MaxAge             time.Duration // MaxAge indicates how long the results of a preflight request can be cached.
}

type formattedOptions struct {
	AllowedHeaders     string
	AllowedMethods     string
	AllowedOrigins     string
	CredentialsAllowed string
	ExposedHeaders     string
	MaxAge             string
}

// Use tells the core to use this handler with the provided options.
func Use(options *Options) {
	fmtOpt := formatCORS(options)
	core.Use(func(c *core.Context) {
		setCORS(c, fmtOpt, func() {
			c.Next()
		})
	})
}

// LocalUse allows to set CORS locally, for a single handler.
// Remember you can't set headers after Write or WriteHeader has been called.
func LocalUse(c *core.Context, options *Options, handler func()) {
	setCORS(c, formatCORS(options), handler)
}

func formatCORS(opt *Options) *formattedOptions {
	if opt == nil {
		opt = new(Options)
	}

	fmtOpt := &formattedOptions{
		AllowedHeaders:     strings.Join(opt.AllowedHeaders, ", "),
		AllowedMethods:     strings.Join(opt.AllowedMethods, ", "),
		CredentialsAllowed: strconv.FormatBool(opt.CredentialsAllowed),
		ExposedHeaders:     strings.Join(opt.ExposedHeaders, ", "),
		MaxAge:             fmt.Sprintf("%.f", opt.MaxAge.Seconds()),
	}

	// If no allowed origins set, all origins are allowed.
	if opt.AllowedOrigins == nil {
		fmtOpt.AllowedOrigins = "*"
	} else {
		fmtOpt.AllowedOrigins = strings.Join(opt.AllowedOrigins, ", ")
	}

	return fmtOpt
}

// setCORS set the response headers and continue if it's not a preflight request.
func setCORS(c *core.Context, fmtOpt *formattedOptions, handler func()) {
	c.ResponseWriter.Header().Set(headerAllowOrigin, fmtOpt.AllowedOrigins)
	if fmtOpt.AllowedOrigins != "*" {
		c.ResponseWriter.Header().Set("Vary", "Origin")
	}

	if fmtOpt.CredentialsAllowed == "true" {
		c.ResponseWriter.Header().Set(headerAllowCredentials, fmtOpt.CredentialsAllowed)
	}

	if fmtOpt.ExposedHeaders != "" {
		c.ResponseWriter.Header().Set(headerExposeHeaders, fmtOpt.ExposedHeaders)
	}

	if fmtOpt.MaxAge != "0" {
		c.ResponseWriter.Header().Set(headerMaxAge, fmtOpt.MaxAge)
	}

	// OPTIONS is used for preflight requests.
	// For that, only the CORS handler must respond, so the handlers chain is broken.
	if c.Request.Method == "OPTIONS" {
		// If no allowed headers are set, accept all from the real request.
		if fmtOpt.AllowedHeaders == "" {
			c.ResponseWriter.Header().Set(headerAllowHeaders, c.Request.Header.Get(headerRequestHeaders))
		} else {
			c.ResponseWriter.Header().Set(headerAllowHeaders, fmtOpt.AllowedHeaders)
		}

		// If no allowed methods are set, accept the method of the real request.
		if fmtOpt.AllowedMethods == "" {
			c.ResponseWriter.Header().Set(headerAllowMethods, c.Request.Header.Get(headerRequestMethod))
		} else {
			c.ResponseWriter.Header().Set(headerAllowMethods, fmtOpt.AllowedMethods)
		}

		c.ResponseWriter.WriteHeader(http.StatusOK)
	} else {
		handler()
	}
}
