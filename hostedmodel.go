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
	if !isValidHostedModelsV1URL(url) {
		return nil, InvlaidURLError
	}
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
		return false, err
	}
	status, ok := meta["status"]
	if !ok {
		return false, UnexpectedError
	}
	return status == "running", nil
}

func (model *HostedModel) WaitUntilAwake(pollIntervalMillis int) error {
	intervalMillis := time.Duration(math.Max(float64(pollIntervalMillis), float64(500)))
	for {
		awake, err := model.IsAwake()
		if err != nil {
			return err
		}
		if awake {
			return nil
		}
		time.Sleep(intervalMillis * time.Millisecond)
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

	model.addRequestHeaders(request)
	response, err := doRequestWithRetry([]int{425, 502}, request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if isHostedModelResponseError(response) {
		if response.StatusCode == 401 {
			return nil, PermissionDeniedError
		} else if response.StatusCode == 404 {
			return nil, NotFoundError
		} else if response.StatusCode == 500 {
			return nil, ModelError
		}
		return nil, UnexpectedError
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, UnexpectedError
	}

	var output JSONObject
	if err := json.Unmarshal(responseBody, &output); err != nil {
		return nil, UnexpectedError
	}
	return output, nil
}
