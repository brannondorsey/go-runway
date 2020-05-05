package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	runway "github.com/brannondorsey/go-runway"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	args := parseArgs()

	model, err := runway.NewHostedModel(args.Url, args.Token)
	if err != nil {
		panic(fmt.Errorf("Error creating HostedModel: %w", err))
	}

	if args.Command == "info" {
		info, err := model.Info()
		if err != nil {
			panic(fmt.Errorf("Error in model.Info(): %w", err))
		}
		fmt.Println(jsonObjectToPretty(info))
		return
	}

	if args.Command == "query" {
		input, err := queryArgumentToJSONObject(args.Arguments[1])
		if err != nil {
			panic(fmt.Errorf("Error converting query argument to JSONObject: %w", err))
		}
		output, err := model.Query(input)
		if err != nil {
			panic(fmt.Errorf("Error in model.Query(): %w", err))
		}
		fmt.Println(jsonObjectToPretty(output))
	}
}

type Args struct {
	Url       string
	Token     string
	Command   string
	Arguments []string
}

func parseArgs() Args {

	url := flag.String("url", "", "A text-generation (GPT-2) hosted model url (e.g. https://my-text-model.hosted-models.runwayml.cloud/v1)")
	token := flag.String("token", "", "The hosted model token. Required if model is private.")

	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <info|query> [json-input-file-or-literal] ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *url == "" || flag.NArg() < 1 {
		usageAndExit()
	}

	command := flag.Args()[0]
	if !(command == "info" || command == "query") {
		usageAndExit()
	}

	if command == "query" && flag.NArg() != 2 {
		fmt.Println("The query command takes a input single argument. It must be a path to a JSON file or a JSON literal.")
		os.Exit(1)
	}

	if command == "info" && flag.NArg() != 1 {
		fmt.Println("The info command does not take an argument.")
		os.Exit(1)
	}

	return Args{
		Url:       *url,
		Token:     *token,
		Command:   command,
		Arguments: flag.Args(),
	}
}

func usageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func jsonObjectToPretty(object runway.JSONObject) string {
	pretty, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		panic(fmt.Errorf("Error in jsonObjectToPretty: %w", err))
	}
	return string(pretty)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func queryArgumentToJSONObject(argument string) (runway.JSONObject, error) {
	jsonLiteral := argument
	var err error
	if fileExists(argument) {
		jsonLiteral, err = getFileContents(argument)
		panicOnError(err)
	}
	var object runway.JSONObject
	err = json.Unmarshal([]byte(jsonLiteral), &object)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON: %w", err)
	}
	return object, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getFileContents(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error opening %v: %w", filename, err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading contents of %v: %w", filename, err)
	}
	return string(contents), nil
}