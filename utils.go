package runway

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func isValidHostedModelsV1URL(url string) bool {
	re := regexp.MustCompile(`^https{0,1}:\/\/.+\.runwayml\.cloud\/v1/{0,1}$`)
	return re.Match([]byte(url))
}

func isHostedModelResponseError(response *http.Response) bool {
	return !strings.Contains(response.Header.Get("Content-Type"), "application/json") ||
		!(response.StatusCode >= 200 && response.StatusCode < 300)
}

func doRequestWithRetry(responseCodesToRetry []int, request *http.Request) (*http.Response, error) {
	client := http.Client{}
	for {
		response, err := client.Do(request)
		if err != nil {
			return nil, &ErrNetworkError{err}
		}
		if !intSliceIncludes(responseCodesToRetry, response.StatusCode) {
			return response, nil
		}
		response.Body.Close()
		// Requests that are repeated must be cloned, as their body has already been read
		request, err = cloneRequest(request)
		if err != nil {
			return nil, err
		}
	}
}

func cloneRequest(request *http.Request) (*http.Request, error) {
	clone := request.Clone(request.Context())
	clonedBody, err := request.GetBody()
	if err != nil {
		return nil, fmt.Errorf("Error getting request body: %s", err)
	}
	clone.Body = clonedBody
	return clone, nil
}

func intSliceIncludes(haystack []int, needle int) bool {
	for i := 0; i < len(haystack); i++ {
		if needle == haystack[i] {
			return true
		}
	}
	return false
}

func stripTrailingSlashIfExists(input string) string {
	return strings.TrimRight(input, "/")
}
