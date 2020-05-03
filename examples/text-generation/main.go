// This example shows how to use every public method of the HostedModels object
// using a text-generation model (GPT-2). You should create a Hosted Model on RunwayML
// before running this example. You can train your own text-generation model or use
// one of the ones publicly available on the platform.
// See https://learn.runwayml.com/#/how-to/hosted-models for details.
//
// Usage ./build/bin/text-generation
//   -prompt string
//     	An optional prompt to use when querying the model. (default "Four score and seven years ago")
//   -token string
//     	The hosted model token. Required if model is private.
//   -url string
//     	A text-generation (GPT-2) hosted model url (e.g. https://my-text-model.hosted-models.runwayml.cloud/v1)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	runway "github.com/brannondorsey/go-runway"
)

func main() {

	args := parseArgs()
	model, err := runway.NewHostedModel(args.Url, args.Token)
	if err != nil {
		panic(err)
	}

	fmt.Println("[INFO] Waiting for model to wake up...")
	pollIntervalMillis := 1000 // check if the model is awake every second
	err = model.WaitUntilAwake(pollIntervalMillis)
	if err != nil {
		panic(err)
	}

	// You can also check if the model is awake via the model.IsAwake() method. This is
	// unnecessary when combined with WaitUntilAwake() above and used here simply for
	// demonstration.
	awake, err := model.IsAwake()
	if err != nil {
		panic(err)
	}
	if awake {
		fmt.Println("[INFO] Model is awake")
	}

	fmt.Println("[INFO] Calling Model.Info()...")
	info, err := model.Info()
	if err != nil {
		panic(err)
	} else {
		pretty, _ := json.MarshalIndent(info, "", "    ")
		fmt.Printf("[INFO] Received response from model.info(): \n%v\n\n", string(pretty))
	}

	// Query
	input := runway.JSONObject{
		"prompt":         args.Prompt,
		"seed":           rand.Intn(1000),
		"max_characters": 512,
	}

	fmt.Println("[INFO] Calling Model.Query()...")
	output, err := model.Query(input)
	if err != nil {
		panic(err)
	} else {
		pretty, _ := json.MarshalIndent(output, "", "    ")
		fmt.Printf("[INFO] Received response from model.query():\n%v\n\n", string(pretty))
	}
}

type Args struct {
	Url    string
	Token  string
	Prompt string
}

func parseArgs() Args {
	url := flag.String("url", "", "A text-generation (GPT-2) hosted model url (e.g. https://my-text-model.hosted-models.runwayml.cloud/v1)")
	token := flag.String("token", "", "The hosted model token. Required if model is private.")
	prompt := flag.String("prompt", "Four score and seven years ago", "An optional prompt to use when querying the model.")
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	return Args{
		Url:    *url,
		Token:  *token,
		Prompt: *prompt,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
