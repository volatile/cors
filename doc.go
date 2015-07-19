/*
Package cors is a handler for the Core.
It enables Cross-Origin Resource Sharing support.

Make sure to include the handler above any other handler that alter the response body (even before the Compress handler, if you use it).

Documentation about CORS:
- http://www.w3.org/TR/cors/
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
- http://www.html5rocks.com/en/tutorials/cors/

Usage

cors.Use(nil) uses the default configuration and it allows all headers, methods and origins.
If you need a more control, give &cors.Options{} instead of nil.

	package main

	import (
		"fmt"

		"github.com/volatile/core"
		"github.com/volatile/cors"
	)

	func main() {
		cors.Use(&cors.Options{
			AllowedHeaders: []string{"X-Client-Header-Example", "X-Another-Client-Header-Example"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedOrigins: []string{"http://example.com", "http://example.com"},
			CredentialsAllowed: true,
			ExposedHeaders: []string{"X-Header-Example", "X-Another-Header-Example"},
			MaxAge: 365 * 24 * time.Hour,
		})

		core.Use(func(c *core.Context) {
			fmt.Fprint(c.ResponseWriter, "Hello, World!")
		})

		core.Run()
	}
*/
package cors
