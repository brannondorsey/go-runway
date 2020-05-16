// runway is a Go library for interfacing with RunwayML Hosted Models.
package runway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// JSONObject is used to represent JSON structures which are variable or unknown
// at compile time.
type JSONObject map[string]interface{}

// HostedModel represents a RunwayML Hosted Model.
// See learn.runwayml.com/#/how-to/hosted-models for more details.
type HostedModel struct {
	url   string
	token string
}

// NewHostedModel instantiates a HostedModel object and returns a pointer to it.
// "url" is the full URL of your hosted model in the format "https://my-model.hosted-models.runwayml.cloud/v1".
// "token" is the secret token associated with this model, if the model is private.
// Use an empty string "" if the model has no token.
func NewHostedModel(url, token string) (*HostedModel, error) {
	if !isValidHostedModelsV1URL(url) {
		return nil, ErrInvalidURL
	}
	return &HostedModel{
		url:   url,
		token: token,
	}, nil
}

// Info returns a JSONObject containing the input/output spec provided by the model.
// It makes a GET request to the /v1/info route of a hosted model under the hood.
func (model *HostedModel) Info() (JSONObject, error) {
	return model.requestHostedModel("GET", model.GetURL()+"/info", nil)
}

// Query runs the model on your input and produce an output. This is how you "run" the model.
// "input" is an object containing input parameters to be sent to the model.
// Use the HostedModel.Info() method to get the correct format for this JSONobject,
// as each model expects different inputs.
func (model *HostedModel) Query(input JSONObject) (JSONObject, error) {
	return model.requestHostedModel("POST", model.GetURL()+"/query", input)
}

// IsAwake returns true if this model is awake, and false if it is still waking up.
// See the Awake, Awakening, and Awake in the Hosted Models docs for more info:
// https://learn.runwayml.com/#/how-to/hosted-models?id=asleep-awakening-and-awake-states.
func (model *HostedModel) IsAwake() (bool, error) {
	var meta JSONObject
	meta, err := model.root()
	if err != nil {
		return false, err
	}
	status, ok := meta["status"]
	if !ok {
		return false, &ErrUnexpectedError{
			fmt.Errorf("\"status\" is missing from model root response"),
		}
	}
	return status == "running", nil
}

// WaitUntilAwake returns once the model is awake. This method is never required, as
// HostedModel.Info() and HostedModel.Query() will always return results eventually, but
// it can be useful for managing UI if you want to postpone making Info() and Query()
// requests until you know that they will resolve more quickly.
// pollIntervalMillis controls the frequency this method will make HTTP requests to the
// underlying Hosted Model to check if it is awake yet.
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

// GetURL returns the hosted model's url with trailing slashes removed if they were
// present during creation
func (model *HostedModel) GetURL() string {
	return stripTrailingSlashIfExists(model.url)
}

func (model *HostedModel) root() (JSONObject, error) {
	return model.requestHostedModel("GET", model.GetURL(), nil)
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
			return nil, &ErrInvalidArgument{
				ArgumentName: "body",
				Details:      "Argument is not valid JSON",
			}
		}
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, &ErrInvalidArgument{
			Err: err,
		}
	}

	model.addRequestHeaders(request)
	response, err := doRequestWithRetry([]int{429, 502}, request)
	if err != nil {
		var networkError *ErrNetworkError
		if errors.As(err, &networkError) {
			return nil, err
		} else {
			return nil, &ErrUnexpectedError{
				fmt.Errorf("Error doing request with retry: %w", err),
			}
		}
	}
	defer response.Body.Close()

	if isHostedModelResponseError(response) {
		if response.StatusCode == 401 {
			return nil, ErrPermissionDenied
		} else if response.StatusCode == 404 {
			return nil, ErrNotFound
		} else if response.StatusCode == 500 {
			return nil, ErrModelError
		}
		return nil, &ErrUnexpectedError{
			fmt.Errorf("Unexpected HTTP response status code %v:", response.StatusCode),
		}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &ErrUnexpectedError{
			fmt.Errorf("Unable to read response body: %w", err),
		}
	}

	var output JSONObject
	if err := json.Unmarshal(responseBody, &output); err != nil {
		return nil, &ErrUnexpectedError{
			fmt.Errorf("Response body is not json: %w", err),
		}
	}
	return output, nil
}
