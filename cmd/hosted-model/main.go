package main

import (
	"fmt"
	"math/rand"

	runway "github.com/brannondorsey/go-runway"
)

func main() {

	// Instantiate
	url := "https://lotr.hosted-models.runwayml.cloud/v1/"
	token := "9M1i6hDO/74cJrrYb2KCMg=="
	model, err := runway.NewHostedModel(url, token)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Hosted model created with url:", url)
	}

	fmt.Println("Waiting for model to wake up...")
	err = model.WaitUntilAwake(1000)
	if err != nil {
		panic(err)
	}
	fmt.Println("Model is now awake")

	if awake, err := model.IsAwake(); !awake || err != nil {
		if err != nil {
			panic(err)
		} else {
			panic("WaitUntilAwake() reported the model was awake, but IsAwake() said it wasn't")
		}
	}

	// Info
	info, err := model.Info()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Received description from model info: %+v\n", info["description"])
	}

	// Query
	input := runway.JSONObject{
		"prompt":         "Four score and seven years ago",
		"seed":           rand.Intn(1000),
		"max_characters": 256,
	}

	output, err := model.Query(input)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Received response from model: %+v\n", output["generated_text"])
	}
}
