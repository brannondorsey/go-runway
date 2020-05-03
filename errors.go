package runway

import (
	"errors"
)

var (
	ErrPermissionDenied = errors.New("Permission denied, this model is private. Did you include the correct token?")
	ErrNotFound         = errors.New("Model not found. Make sure the url is correct and that the model is \"active\".")
	ErrInvlaidURL       = errors.New("The url you've provided is not valid. Your Hosted Model your must be in the format https://my-model.hosted-models.runwayml.cloud/v1.")
	ErrUnexpectedError  = errors.New("An unexpected error has ocurred. Please try again later or contact support (https://support.runwayml.com).")
	ErrNetworkError     = errors.New("A network error has ocurred. Please check your internet connection is working properly and try again.")
	ErrModelError       = errors.New("The model experienced an error while processing your input. Double-check that you are sending properly formed input parameters in HostedModel.query(). You can use the HostedModel.info() method to check the input parameters the model expects. If the error persists, contact support (https://support.runwayml.com).")
	ErrInvalidArgument  = errors.New("Invalid argument. Make sure the arguments to this function are the correct types and values.")
)
