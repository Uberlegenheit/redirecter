package apperrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	ErrCode string
)

//TODO after bitbucket deploy setup import from everstake-common
const (
	ErrService       ErrCode = "ERR_SERVICE"
	ErrNotFound      ErrCode = "ERR_NOT_FOUND"
	ErrBadRequest    ErrCode = "ERR_BAD_REQUEST"
	ErrBadParam      ErrCode = "ERR_BAD_PARAM"
	ErrTimeIsUp      ErrCode = "ERR_TIME_IS_UP"
	ErrNotAllowed    ErrCode = "ERR_NOT_ALLOWED"
	ErrAlreadyExists ErrCode = "ERR_ALREADY_EXISTS"
	ErrBadAuth       ErrCode = "ERR_BAD_AUTH"
	ErrBadAuthCookie ErrCode = "ERR_BAD_AUTH_COOKIE"
)

type (
	ServiceError interface {
		error
		ErrorCode() ErrCode
		ToMap(http.ResponseWriter) map[string]interface{}
		GetHttpCode() int
	}

	Error struct {
		Code        ErrCode `json:"code"`
		Value       string  `json:"value,omitempty"`
		Description string  `json:"description,omitempty"`
	}
)

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", string(e.Code), e.Value)
}

func (e Error) ErrorCode() ErrCode {
	return e.Code
}

// ToMap converts Error object to map[string]interface{}
func (e Error) ToMap() map[string]interface{} {
	r := map[string]interface{}{
		"error": string(e.Code),
	}

	if string(e.Value) != "" {
		r["value"] = string(e.Value)
	}

	if string(e.Description) != "" {
		r["description"] = string(e.Description)
	}

	return r
}

// GetHttpCode return a Http error code
func (e Error) GetHttpCode() int {
	switch e.Code {
	case ErrService:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// New creates an Error object
func New(code ErrCode, value ...string) *Error {
	e := &Error{Code: code}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// NewWithDesc creates an Error object with description
func NewWithDesc(code ErrCode, desc string, value ...string) *Error {
	e := &Error{Code: code, Description: desc}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// FromError creates a new Error (ErrService) from common golang error
func FromError(err error) *Error {
	if err != nil {
		return &Error{
			Code:  ErrService,
			Value: "",
		}
	}

	return nil
}

func AppEncode(err error) error {
	switch err.(type) {
	case Error, *Error:
		text, mErr := json.Marshal(err)
		if mErr != nil {
			return err
		}
		return errors.New(string(text))
	}
	return err
}

func AppDecode(argErr error) *Error {
	if argErr == nil {
		return nil
	}
	var sErr Error
	err := json.Unmarshal([]byte(argErr.Error()), &sErr)
	if err != nil {
		return nil
	}
	return &sErr
}
