/*
Package cors is a handler for the core (https://godoc.org/github.com/volatile/core).
It provides Cross-Origin Resource Sharing support.

Usage

When using CORS (globally or locally), there is always a parameter of type OriginsMap.
It can contain a map of allowed origins and their specific options.

Use nil to allow all headers, methods and origins:

	cors.Use(nil)

Use nil for origin to allow all headers and methods for this origin.

	cors.Use(&cors.OriginsMap{
		"example.com": nil,
	})

Use AllOrigins to set options for all origins.

	cors.Use(&cors.OriginsMap{
		"example.com": nil, // All is allowed for this origin.
		cors.AllOrigins: &cors.Options{
			AllowedMethods: []string{"GET"}, // Only the GET method is allowed for the others.
		},
	})

Global usage

Use sets a global CORS configuration for all the handlers downstream:

	cors.Use(nil)

	// All is allowed for the root path.
	core.Use(func(c *core.Context) {
		if c.Request.URL.Path == "/" {
			fmt.Fprint(c.ResponseWriter, "Hello, World!")
		} else {
			c.Next()
		}
	})

	// Previous CORS options are overwritten.
	cors.Use(&cors.OriginsMap{
		cors.AllOrigins: &cors.Options{
			AllowedMethods: []string{"GET"},
		},
	})

	// Only the GET method is allowed for this handler.
	core.Use(func(c *core.Context) {
		fmt.Fprint(c.ResponseWriter, "Read only")
	})

Make sure to include the handler above any other handler that alter the response body.

Local usage

LocalUse sets CORS locally, inside a single handler.
This setting takes precedence over he global CORS options (if set).

	core.Use(func(c *core.Context) {
		cors.LocalUse(c, &cors.OriginsMap{
			cors.AllOrigins: &cors.Options{AllowedMethods: []string{"GET"}},
		}, func() {
			response.Status(c, http.StatusOK)
		})
	})

Documentation

For more information:

	W3C		http://www.w3.org/TR/cors/
	Mozilla		https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
	HTML5 Rocks	http://www.html5rocks.com/en/tutorials/cors/
*/
package cors
