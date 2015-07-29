/*
Package cors is a handler for the Core.
It enables Cross-Origin Resource Sharing support.

Make sure to include the handler above any other handler that alter the response body.

Documentation about CORS:
- http://www.w3.org/TR/cors/
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
- http://www.html5rocks.com/en/tutorials/cors/

Usage

When using CORS (globally or locally), there is always a parameter of type "*cors.Options".
If you set "nil" for this parameter, the default configuration is used: it allows all headers, methods and origins.
If you need a more control, give custom options with "&cors.Options{}".

Global usage

cors.Use(*cors.Options) sets a global CORS configuration for all the handlers.

	package main

	import (
		"fmt"

		"github.com/volatile/core"
		"github.com/volatile/cors"
	)

	func main() {
		cors.Use(nil)

		core.Use(func(c *core.Context) {
			fmt.Fprint(c.ResponseWriter, "Hello, World!")
		})

		core.Run()
	}

Local usage

cors.LocalUse(*core.Context, *cors.Options, func()) allows to set CORS locally, for a single handler.
The global CORS options are overwritten in this situation.

The last func() parameter is called after the CORS headers are set, but only if it's not a [preflight request](http://www.w3.org/TR/cors/#resource-preflight-requests).

	package main

	import (
		"fmt"
		"time"

		"github.com/volatile/core"
		"github.com/volatile/cors"
		"github.com/volatile/response"
		"github.com/volatile/route"
	)

	func main() {
		// Global use
		cors.Use(nil)

		// Local use for the "/hook" path.
		route.Get("^/hook$", func(c *core.Context) {
			opt := &cors.Options{
				AllowedHeaders:     []string{"X-Client-Header-Example", "X-Another-Client-Header-Example"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowedOrigins:     []string{"http://example.com", "http://example.com"},
				CredentialsAllowed: true,
				ExposedHeaders:     []string{"X-Header-Example", "X-Another-Header-Example"},
				MaxAge:             365 * 24 * time.Hour,
			}

			cors.LocalUse(c, opt, func() {
				response.Status(c, http.StatusOK)
			})
		})

		// No local CORS are set: global CORS options are used.
		core.Use(func(c *core.Context) {
			response.String(c, "Hello, World!")
		})

		core.Run()
	}
*/
package cors
