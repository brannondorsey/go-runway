package runway

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPermissionDenied = errors.New("Permission denied, this model is private. Did you include the correct token?")
	ErrNotFound         = errors.New("Model not found. Make sure the url is correct and that the model is \"active\".")
	ErrInvalidURL       = errors.New("The hosted model url you've provided is invalid. It must be in the format https://my-model.hosted-models.runwayml.cloud/v1.")
	ErrModelError       = errors.New("The model experienced an error while processing your input. Double-check that you are sending properly formed input parameters in HostedModel.Query(). You can use the HostedModel.Info() method to check the input parameters the model expects. If the error persists, contact support (https://support.runwayml.com).")
)

type ErrInvalidArgument struct {
	ArgumentName string
	Details      string
	Err          error
}

func (err *ErrInvalidArgument) Error() string {
	message := ""
	if err.ArgumentName == "" {
		message += fmt.Sprintf("Invalid argument: %s. ", err.ArgumentName)
	}
	if err.Details == "" {
		message += fmt.Sprintf("%s. ", err.Details)
	}
	if err.Err == nil {
		message += fmt.Sprintf("%v. ", err.Err)
	}
	return strings.TrimSpace(message)
}

type ErrNetworkError struct {
	Err error
}

func (err *ErrNetworkError) Error() string {
	return fmt.Sprintf("A network error has ocurred. Please check your internet connection is working properly and try again. Network error details: %s", err.Unwrap())
}

func (err *ErrNetworkError) Unwrap() error {
	return err.Err
}

type ErrUnexpectedError struct {
	Err error
}

func (err *ErrUnexpectedError) Error() string {
	return fmt.Sprintf("An unexpected error has ocurred. Please try again later or contact support (https://support.runwayml.com). Error details: %s", err.Unwrap())
}

func (err *ErrUnexpectedError) Unwrap() error {
	return err.Err
}
