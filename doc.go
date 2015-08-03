/*
Package cors is a handler for the Core.
It enables Cross-Origin Resource Sharing support.

Make sure to include the handler above any other handler that alter the response body.

Documentation about CORS:

- http://www.w3.org/TR/cors/

- https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS

- http://www.html5rocks.com/en/tutorials/cors/

Usage

When using CORS (globally or locally), there is always a parameter of type cors.OriginsMap.
It can contain a map of allowed origins and their specific options.

- Use nil as cors.OriginsMap to allow all headers, methods and origins.

	cors.Use(nil)

- Use nil as origin's *Options to allow all headers and methods for this origin.

	cors.Use(cors.OriginsMap{
		"example.com": nil,
	})

- Use cors.AllOrigins as an cors.OriginsMap key to set options for all origins.

	cors.Use(cors.OriginsMap{
		"example.com": nil, // All is allowed for this origin.
		cors.AllOrigins: &cors.Options{
			AllowedMethods: []string{"GET"}, // Only the GET method is allowed for the others.
		},
	})


Global

cors.Use(cors.OriginsMap) sets a global CORS configuration for all the handlers downstream.

	package main

	import (
		"fmt"

		"github.com/volatile/core"
		"github.com/volatile/cors"
	)

	func main() {
		cors.Use(nil)

		// All is allowed for the "/" path.
		core.Use(func(c *core.Context) {
			if c.Request.URL.Path == "/" {
				fmt.Fprint(c.ResponseWriter, "Hello, World!")
			}
		})

		// The previous CORS options are overwritten.
		cors.Use(cors.OriginsMap{
			cors.AllOrigins: &cors.Options{
				AllowedMethods: []string{"GET"},
			},
		})

		// Only the GET method is allowed for this handler.
		core.Use(func(c *core.Context) {
			fmt.Fprint(c.ResponseWriter, "Read only")
		})

		core.Run()
	}

Local

cors.LocalUse(*core.Context, cors.OriginsMap, func()) can be used to set CORS locally, for a single handler.
The global CORS options (if used) are overwritten in this situation.

The last func() parameter is called after the CORS headers are set, but only if it's not a preflight request (http://www.w3.org/TR/cors/#resource-preflight-requests).

	package main

	import (
		"fmt"

		"github.com/volatile/core"
		"github.com/volatile/cors"
	)

	func main() {
		// Global use
		cors.Use(nil)

		// Local use for the "/hook" path.
		core.Use(func(c *core.Context) {
			if c.Request.URL.Path == "/hook" {
				cors.LocalUse(c, cors.OriginsMap{
					cors.AllOrigins: &cors.Options{AllowedMethods: []string{"GET"}},
				}, func() {
					response.Status(c, http.StatusOK)
				})
			}
			c.Next()
		})

		// No local CORS are set, so the global CORS options are used.
		core.Use(func(c *core.Context) {
			fmt.Fprint(c.ResponseWriter, "Hello, World!")
		})

		core.Run()
	}
*/
package cors
