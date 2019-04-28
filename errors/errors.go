package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode struct {
	Kind             string `json:"kind"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
}

// APIError represents an error that can be sent in an error response.
type APIError struct {
	HTTPStatus int    `json:"-"`
	Err        string `json:"error,omitempty"`
	*ErrorCode
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode.Kind, e.ErrorCode.Message)
}

func newErrorCode(kind, message string) *ErrorCode {
	return &ErrorCode{
		Kind:    kind,
		Message: message,
	}
}

func newAPIError(httpStatus int, errCode *ErrorCode) *APIError {
	err := &APIError{
		HTTPStatus: httpStatus,
	}
	err.ErrorCode = errCode
	return err
}

func NewRecordNotFound() *APIError {
	return newAPIError(http.StatusNotFound, ErrCodeRecordNotFound)
}

func NewUnauthorizedError() *APIError {
	return newAPIError(http.StatusUnauthorized, ErrCodeUnauthorized)
}

func NewForbiddenError() *APIError {
	return newAPIError(http.StatusForbidden, ErrCodeForbidden)
}

type InputDataError struct {
	APIError
	errorFields []string
}

func (e *InputDataError) AppendFields(fields ...string) {
	e.errorFields = append(e.errorFields, fields...)
}

func (e *InputDataError) Valid() bool {
	if len(e.errorFields) == 0 {
		return true
	}
	e.Message = fmt.Sprintf("Invalid fields: %s", e.errorFields)
	return false
}

func NewInputDataErrorOld(fields ...string) *InputDataError {
	err := &InputDataError{}
	err.APIError.HTTPStatus = http.StatusBadRequest
	err.APIError.ErrorCode = ErrCodeInvalidInputData
	if len(fields) > 0 {
		err.AppendFields(fields...)
		err.Valid()
	}
	return err
}

func NewInputDataError(err error) *APIError {
	e := &APIError{
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  ErrCodeInvalidInputData,
	}
	e.DeveloperMessage = err.Error()
	return e
}

func NewIncorrectUnmarshal(errCode *ErrorCode) *APIError {
	return newAPIError(http.StatusBadRequest, errCode)
}

func NewBadRequest(errCode *ErrorCode) *APIError {
	return newAPIError(http.StatusBadRequest, errCode)
}

func NewConflict(errCode *ErrorCode) *APIError {
	return newAPIError(http.StatusConflict, errCode)
}

func NewInternalServerError() *APIError {
	e := newAPIError(http.StatusInternalServerError, ErrCodeInternal)
	return e
}

func IsRecordNotFoundError(err error) bool {
	if e, ok := err.(*APIError); ok == true {
		return e.Kind == ErrCodeRecordNotFound.Kind
	}
	return false
}
