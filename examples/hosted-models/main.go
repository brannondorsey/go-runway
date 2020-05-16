package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	runway "github.com/brannondorsey/go-runway"
	flag "github.com/spf13/pflag"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	args := parseArgs()

	model, err := runway.NewHostedModel(args.Url, args.Token)
	if err != nil {
		printAndExitIfError(fmt.Errorf("Error creating HostedModel: %s", err))
	}

	if args.Command == "info" {
		info, err := model.Info()
		if err != nil {
			printAndExitIfError(fmt.Errorf("Error in model.Info(): %s", err))
		}
		fmt.Fprintln(os.Stdout, jsonObjectToPretty(info))
		return
	}

	if args.Command == "query" {
		input, err := queryArgumentToJSONObject(args.Arguments[1])
		if err != nil {
			printAndExitIfError(fmt.Errorf("Error converting query argument to JSONObject: %s", err))
		}
		output, err := model.Query(input)
		if err != nil {
			printAndExitIfError(fmt.Errorf("Error in model.Query(): %s", err))
		}
		fmt.Fprintln(os.Stdout, jsonObjectToPretty(output))
	}
}

type Args struct {
	Url       string
	Token     string
	Command   string
	Arguments []string
}

func parseArgs() Args {

	url := flag.StringP("url", "u", "", "The hosted model url (e.g. https://my-text-model.hosted-models.runwayml.cloud/v1)")
	token := flag.StringP("token", "t", "", "The hosted model token. Required if model is private.")
	help := flag.BoolP("help", "h", false, "Print this help screen.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <COMMAND> [ARGUMENTS] ...\n", os.Args[0])
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  info: Print the input/output spec of the model as JSON. Equivalent to GET /v1/info on the hosted model URL.")
		fmt.Fprintln(os.Stderr, "  query <file-or-json-literal>: Query the model, using a JSON argument as input. The JSON argument can be a file or a JSON string.")
		fmt.Fprintln(os.Stderr)
	}

	flag.Parse()

	if *help || *url == "" {
		usageAndExit("")
	}

	if flag.NArg() < 1 {
		usageAndExit("Error: A command is required.")
	}

	command := flag.Args()[0]
	if !(command == "info" || command == "query") {
		usageAndExit(fmt.Sprintf("Error: Invalid command \"%s\".", command))
	}

	if command == "query" && flag.NArg() != 2 {
		usageAndExit("Error: The query command takes a input single argument. It must be a path to a JSON file or a JSON literal.")
	}

	if command == "info" && flag.NArg() != 1 {
		usageAndExit(fmt.Sprintf("Error: The %s command does not take an argument.", command))
	}

	return Args{
		Url:       *url,
		Token:     *token,
		Command:   command,
		Arguments: flag.Args(),
	}
}

func usageAndExit(optionalMessage string) {
	flag.Usage()
	if optionalMessage != "" {
		fmt.Fprintln(os.Stderr, optionalMessage)
	}
	os.Exit(1)
}

func jsonObjectToPretty(object runway.JSONObject) string {
	pretty, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		printAndExitIfError(fmt.Errorf("Error in jsonObjectToPretty(): %s", err))
	}
	return string(pretty)
}

func printAndExitIfError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func queryArgumentToJSONObject(argument string) (runway.JSONObject, error) {
	jsonLiteral := argument
	var err error
	// Ignore false negatives that may occur if checking for a file returns an error
	// as we probably couldn't open it anyway
	if exists, _ := fileExists(argument); exists {
		jsonLiteral, err = getFileContents(argument)
		printAndExitIfError(err)
	}
	var object runway.JSONObject
	err = json.Unmarshal([]byte(jsonLiteral), &object)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON: %s", err)
	}
	return object, nil
}

func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return !info.IsDir(), nil
}

func getFileContents(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error opening %v: %s", filename, err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading contents of %v: %s", filename, err)
	}
	return string(contents), nil
}
