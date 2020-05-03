package main

import (
	"fmt"

	runway "github.com/brannondorsey/go-runway"
)

func main() {

	// Instantiate
	url := "https://lotr.hosted-models.runwayml.cloud/v1/"
	token := "9M1i6hDO/74cJrrYb2KCMg=="
	model, err := runway.NewHostedModel(url, token)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Hosted model created with url:", url)
	}

	// Info
	info, err := model.Info()
	if err != nil {
		fmt.Println("Error getting model info")
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", info)
	}

	// Query
	input := map[string]string{
		"prompt": "Four score and seven years ago",
	}
	output, err := model.Query(input)
	if err != nil {
		fmt.Println("Error querying model")
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", output)
	}
}
