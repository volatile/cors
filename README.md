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

If you give `nil` as the `*cors.Options` parameter, the default configuration is used and it allows all headers, methods and origins.  
If you need a more control, give custom options with `&cors.Options{}` instead.

### Global

`cors.Use(*Options)` sets a global CORS configuration for all the handlers.

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

`cors.LocalUse(*core.Context, *cors.Options, func())` allows to set CORS locally, for a single handler.  
 The global CORS options are overwritten in this situation.

The last `func()` parameter is called after the CORS headers are set, but only if it's not a [preflight request](http://www.w3.org/TR/cors/#resource-preflight-requests).

```Go
package main

import (
	"fmt"
	"time"

	"github.com/volatile/core"
	"github.com/volatile/cors"
)

func main() {
	// Global use
	cors.Use(nil)

	core.Use(func(c *core.Context) {
		// Custom CORS are set only for the "/localcors" path.
		if c.Request.URL.String() == "/localcors" {
			opt := &cors.Options{
				AllowedHeaders:     []string{"X-Client-Header-Example", "X-Another-Client-Header-Example"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowedOrigins:     []string{"http://example.com", "http://example.com"},
				CredentialsAllowed: true,
				ExposedHeaders:     []string{"X-Header-Example", "X-Another-Header-Example"},
				MaxAge:             365 * 24 * time.Hour,
			}

			// Local use
			cors.LocalUse(c, opt, func() {
				fmt.Fprint(c.ResponseWriter, "Hello, World!")
			})
		}

		c.Next()
	})

	core.Use(func(c *core.Context) {
		// No local CORS are set. Then, obviously, global CORS are used.
		fmt.Fprint(c.ResponseWriter, "Hello, World!")
	})

	core.Run()
}
```

[![GoDoc](https://godoc.org/github.com/volatile/cors?status.svg)](https://godoc.org/github.com/volatile/cors)
