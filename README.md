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
$ go get github.com/volatile/cors
```

## Usage [![GoDoc](https://godoc.org/github.com/volatile/cors?status.svg)](https://godoc.org/github.com/volatile/cors)

When using *CORS* (globally or locally), there is always a parameter of type [`OriginsMap`](https://godoc.org/github.com/volatile/cors#OriginsMap).  
It can contain a map of allowed origins and their specific options.

- Use `nil` as [`OriginsMap`](https://godoc.org/github.com/volatile/cors#OriginsMap) to allow all headers, methods and origins.
  ```Go
  cors.Use(nil)
  ```

- Use `nil` as origin's [`*Options`](https://godoc.org/github.com/volatile/cors#Options) to allow all headers and methods for this origin.
  ```Go
  cors.Use(cors.OriginsMap{
  	"example.com": nil,
  })
  ```

- Use [`AllOrigins`](https://godoc.org/github.com/volatile/cors#AllOrigins) as an [`OriginsMap`](https://godoc.org/github.com/volatile/cors#OriginsMap) key to set options for all origins.
  ```Go
  cors.Use(cors.OriginsMap{
  	"example.com": nil, // All is allowed for this origin.
  	cors.AllOrigins: &cors.Options{
  		AllowedMethods: []string{"GET"}, // Only the GET method is allowed for the others.
  	},
  })
  ```

### Global

[`Use`](https://godoc.org/github.com/volatile/cors#Use) sets a global *CORS* configuration for all the handlers downstream.

```Go
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
```

### Local

[`LocalUse`](https://godoc.org/github.com/volatile/cors#LocalUse) can be used to set *CORS* locally, for a single handler.  
The global *CORS* options (if used) are overwritten in this situation.

The last `func` parameter is called after the *CORS* headers are set, but only if it's not a [preflight request](http://www.w3.org/TR/cors/#resource-preflight-requests).

```Go
package main

import (
	"fmt"
	"net/http"

	"github.com/volatile/core"
	"github.com/volatile/cors"
	"github.com/volatile/response"
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
```
