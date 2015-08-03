<p align="center"><img src="http://volatile.whitedevops.com/images/repositories/cors/logo.png" alt="Volatile CORS" title="Volatile CORS"><br><br></p>

Volatile CORS is a handler for the [Core](https://github.com/volatile/core).  
It enables *Cross-Origin Resource Sharing* support.

Make sure to include the handler above any other handler that alter the response body.

Documentation about *CORS*:
- [W3C official specification](http://www.w3.org/TR/cors/)
- [Mozilla Developer Network](https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS)
- [HTML5 Rocks](http://www.html5rocks.com/en/tutorials/cors/)

## Installation

```Shell
$ go get -u github.com/volatile/cors
```

## Usage

When using CORS (globally or locally), there is always a parameter of type `cors.OriginsMap`.  
It can contain a map of allowed origins and their specific options.  

- Use `nil` as `cors.OriginsMap` to allow all headers, methods and origins.  

- Use `nil` as origin's `*Options` to allow all headers and methods for this origin.

### Global

`cors.Use(cors.OriginsMap)` sets a global CORS configuration for all the handlers.

```Go
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
```

### Local

`cors.LocalUse(*core.Context, cors.OriginsMap, func())` allows to set CORS locally, for a single handler.  
The global CORS options (if used) are overwritten in this situation.

The last `func()` parameter is called after the CORS headers are set, but only if it's not a [preflight request](http://www.w3.org/TR/cors/#resource-preflight-requests).

```Go
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
			AllowedOrigins:     []string{"http://example.com", "http://another-example.com"},
			CredentialsAllowed: true,
			ExposedHeaders:     []string{"X-Header-Example", "X-Another-Header-Example"},
			MaxAge:             1 * time.Hour,
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
```

[![GoDoc](https://godoc.org/github.com/volatile/cors?status.svg)](https://godoc.org/github.com/volatile/cors)
