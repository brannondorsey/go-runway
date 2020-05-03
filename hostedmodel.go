package runway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type JSONObject map[string]interface{}
type HostedModel struct {
	url   string
	token string
}

func NewHostedModel(url, token string) (*HostedModel, error) {
	return &HostedModel{
		url:   url,
		token: token,
	}, nil
}

func (model *HostedModel) Info() (JSONObject, error) {
	return model.requestHostedModel("GET", model.url+"/info", nil)
}

func (model *HostedModel) Query(input JSONObject) (JSONObject, error) {
	return model.requestHostedModel("POST", model.url+"/query", input)
}

func (model *HostedModel) IsAwake() (bool, error) {
	var meta JSONObject
	meta, err := model.root()
	if err != nil {
		return false, nil
	}
	status, ok := meta["status"]
	if !ok {
		return false, UnexpectedError
	}
	return status == "running", nil
}

func (model *HostedModel) WaitUntilAwake(pollIntervalMillis int) error {
	pollIntervalMillis = int(math.Max(float64(pollIntervalMillis), float64(500)))
	for {
		awake, err := model.IsAwake()
		if err != nil {
			return err
		}
		if awake {
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (model *HostedModel) root() (JSONObject, error) {
	return model.requestHostedModel("GET", model.url, nil)
}

func (model *HostedModel) addRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if model.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", model.token))
	}
}

func (model *HostedModel) requestHostedModel(method, url string, body JSONObject) (JSONObject, error) {

	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, *NewInvalidArgumentError("input")
		}
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, *NewInvalidArgumentError("")
	}

	client := http.Client{}
	model.addRequestHeaders(request)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Received an error during request")
		fmt.Println(err)
		return nil, NetworkError
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, UnexpectedError
	}

	var output JSONObject
	if err := json.Unmarshal(responseBody, &output); err != nil {
		fmt.Println(err)

		return nil, UnexpectedError
	}
	return output, nil
}
