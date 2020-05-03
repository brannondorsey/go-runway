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
	fmt.Println("[INFO] Model is awake")

	// You can also check if the model is awake via the model.IsAwake() method
	if awake, err := model.IsAwake(); !awake || err != nil {
		if err != nil {
			panic(err)
		} else {
			panic("Error: WaitUntilAwake() reported the model was awake, but IsAwake() said it wasn't")
		}
	}

	// Info
	info, err := model.Info()
	if err != nil {
		panic(err)
	} else {
		pretty, _ := json.Marshal(info)
		fmt.Printf("[INFO] Received response from model.info(): ", string(pretty))
	}

	// Query
	input := runway.JSONObject{
		"prompt":         args.Prompt,
		"seed":           rand.Intn(1000),
		"max_characters": 1000,
	}

	output, err := model.Query(input)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("[INFO] Received response from model.query(): %+v\n", output["generated_text"])
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
