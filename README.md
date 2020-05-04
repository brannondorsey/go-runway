# RunwayML Go

A small Go module for interfacing with RunwayML. It currently supports the new [Hosted Models](https://learn.runwayml.com/#/how-to/hosted-models) functionality, modeled after the official [Hosted Models JavaScript SDK](https://github.com/runwayml/hosted-models/).

```bash
go get github.com/brannondorsey/go-runway
```

## Example

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	runway "github.com/brannondorsey/go-runway"
)

func main() {

	// Replace this with the URL of your hosted model (https://learn.runwayml.com/#/how-to/hosted-models)
	url := "https://example-text-generator.hosted-models.runwayml.cloud/v1"

	// Paste your secret token in here. Leave as empty string if the model is public.
	token := ""

	model, err := runway.NewHostedModel(url, token)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	input := runway.JSONObject{
		"prompt":         "Four score and seven years ago",
		"seed":           rand.Intn(1000),
		"max_characters": 512,
	}

	fmt.Println("Querying model...")
	output, err := model.Query(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(output["generated_text"])
}

```

More examples live in the `examples/` directory.

## Docs

See the [Usage](https://github.com/runwayml/hosted-models/) section of the Hosted Models JavaScript SDK for general information about the methods available for the `HostedModel` object, as they are identical to the ones provided by this package.

## Dev

Running `make` will build all `examples/` and place their executables in `build/bin`.

```bash
git clone https://github.com/brannondorsey/go-runway
cd runway

make
```
