<div align="center">
    <h1>💶 → 💵</h1>
rates – Currency converter lib + API with a simple GUI
</div>

---

* Rates-source agnostic currency converter (make your own source by implementing the `conv.RatesSource` interface)
* Simple GUI to convert currencies as a example of how to use the API
* Cache for the rates to avoid unnecessary requests to the source
* Custom TTL for the cache
* [TODO] Dockerfile
* [TODO] Tests

> [!WARNING]  
> Stable version is not released yet. The API and the lib are under development.

---

## Usage

### Standalone

You can use the API as a standalone service.

| URL                                | Description                                                             |
|------------------------------------|-------------------------------------------------------------------------|
| [localhost:80/](localhost:80/)     | Simple GUI with htmx, just to show how to init and use `conv.Converter` |
| [localhost:8080/](localhost:8080/) | http API                                                                |

GUI example:
![img](https://github.com/egregors/rates/assets/2153895/ddff6b77-175a-48bc-828d-c25933cf6921)

API support methods:

| Method | URL             | Description                                                                                        |
|--------|-----------------|----------------------------------------------------------------------------------------------------|
| POST   | /api/v0/convert | Expects: `{amount: 123.45, from: "usd", to: "eur"}`<br/> Returns: `{result": 114.48554492400001 }` |

API request example:

```shell
http localhost:8080/api/v0/convert from=usd to=eur amount=123.45

HTTP/1.1 200 OK
Content-Length: 30
Content-Type: application/json
Date: Sun, 11 Feb 2024 14:02:17 GMT

{
  "result": 114.48554492400001
}
```

### As a library

You can use the `conv` package as a library in your own project.

```go
package main

import (
	"fmt"

	"github.com/egregors/rates/conv"
	"github.com/egregors/rates/conv/backends"
)

func main() {
	// Create a new converter
	c := conv.New(
		backends.NewCurrencyAPI(),
	)

	// Convert 123.45 USD to EUR
	result, err := c.Conv(123.45, "usd", "eur")
	if err != nil {
		panic(err)
	}

	// Print the result
	fmt.Println(result)
}

```

`Conv` supports constructor options.

* `conv.WithCache` – Enable cache for the rates
* `conv.WithLogger` – Use custom logger for the requests

### Make your own rates source

You can make your own rates source by implementing the `conv.RatesSource` interface.

```go
package conv

type RatesSource interface {
	Rate(from, to string) (float64, error)
}

```

`Rate` method should return the rate for the given currencies.

As an example, you can check the `backends.CurrencyAPI` implementation.
At this moment, it uses the [currency-api](https://github.com/fawazahmed0/currency-api) as the source.

## Development

Check the Makefile for the available commands.

```shell
➜  rates git:(main) ✗ make help
Usage: make [task]

task                 help
------               ----
                     
lint                 Lint the files
test                 Run tests
run                  Run the application with watcher (go, gohtml)
                     
help                 Show help message
```

## Contributing

Bug reports, bug fixes and new features are always welcome.
Please open issues and submit pull requests for any new code.
