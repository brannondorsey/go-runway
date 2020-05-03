package runway

import (
	"errors"
	"fmt"
)

var PermissionDeniedError = errors.New("Permission denied, this model is private. Did you include the correct token?")
var NotFoundError = errors.New("Model not found. Make sure the url is correct and that the model is \"active\".")
var InvlaidURLError = errors.New("The url you've provided is not valid. Your Hosted Model your must be in the format https://my-model.hosted-models.runwayml.cloud/v1.")
var UnexpectedError = errors.New("An unexpected error has ocurred. Please try again later or contact support (https://support.runwayml.com).")
var NetworkError = errors.New("A network error has ocurred. Please check your internet connection is working properly and try again.")
var ModelError = errors.New("The model experienced an error while processing your input. Double-check that you are sending properly formed input parameters in HostedModel.query(). You can use the HostedModel.info() method to check the input parameters the model expects. If the error persists, contact support (https://support.runwayml.com).")

type InvalidArgumentError struct {
	ArgumentName string
}

func (e InvalidArgumentError) Error() string {
	message := "Invalid argument."
	if e.ArgumentName != "" {
		message = fmt.Sprintf("The required argument \"%v\" is invalid.")
	}
	return message
}

func NewInvalidArgumentError(argumentName string) *InvalidArgumentError {
	return &InvalidArgumentError{
		ArgumentName: argumentName,
	}
}