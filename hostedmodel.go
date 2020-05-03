package runway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func (model *HostedModel) Info() (interface{}, error) {
	request, err := http.NewRequest("GET", model.url+"/info", nil)
	if err != nil {
		return nil, err
	}
	response, err := model.requestHostedModel(request)
	if err != nil {
		return nil, err
	}
	var info interface{}
	if err := json.Unmarshal(response, &info); err != nil {
		return nil, UnexpectedError
	}
	return &info, nil
}

func (model *HostedModel) Query(input interface{}) (interface{}, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return InvalidArgumentError{
			ArgumentName: "input",
		}, nil
	}
	request, err := http.NewRequest("POST", model.url+"/query", bytes.NewReader(jsonInput))
	if err != nil {
		return err, nil
	}
	response, err := model.requestHostedModel(request)
	if err != nil {
		return err, nil
	}
	var output interface{}
	if err := json.Unmarshal(response, &output); err != nil {
		return UnexpectedError, nil
	}
	return output, nil
}

func (model *HostedModel) root() {

}

func (model *HostedModel) addRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if model.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", model.token))
	}
}

func (model *HostedModel) requestHostedModel(req *http.Request) ([]byte, error) {
	client := http.Client{}
	model.addRequestHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Received an error during request")
		fmt.Println(err)
		return nil, NetworkError
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, UnexpectedError
	}
	return response, nil
}
