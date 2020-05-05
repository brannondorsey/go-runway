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
