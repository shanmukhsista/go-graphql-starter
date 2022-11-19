package apperrors

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/go-multierror"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	ExtensionFieldNameKey = "field"
	ExtensionErrorStatus  = "status"
)

const defaultErrorMessage = "An error occoured. Please contact support."

//go:embed error_messages.json
var errorMessagesContent []byte
var ErrorMessagesMap map[string]string

func init() {
	err := json.Unmarshal(errorMessagesContent, &ErrorMessagesMap)
	if err != nil {
		panic("Messages file required for initializing App. ")
	}
}

func getErrorMessageAndFieldForKey(key string) string {
	msg := ErrorMessagesMap[key]
	if msg == "" {
		return defaultErrorMessage
	}
	return msg
}

type appError struct {
	Status     int
	Field      string
	Key        string
	Underlying error
}

func (r appError) Error() string {
	return fmt.Sprintf("status %d: key %s ", r.Status, r.Key)
}

func (r appError) SetError(err error) {
	r.Underlying = err
}

func (r appError) AddErrorString(err string) {
	r.Underlying = errors.New(err)
}

func (r appError) AddErrorStringForField(err string, field string) {
	r.AddErrorString(err)
	r.Field = field
}

func NewInternalErrorWithUnderlying(key string, err error) *appError {
	return &appError{Underlying: err, Key: key, Status: http.StatusInternalServerError}
}

func NewErrorWithUnderlyingAndStatus(key string, err error, status int) *appError {
	return &appError{Underlying: err, Key: key, Status: status}
}

func NewErrorWithFieldAndStatus(key string, field string, status int) *appError {
	return &appError{Key: key, Status: status, Field: field}
}

func GetAppErrorObject(err error) *appError {
	if apperror, ok := err.(*appError); ok {
		return apperror
	}
	if apperror, ok := err.(appError); ok {
		return &apperror
	}
	appError := appError{Key: err.Error()}
	return &appError
}

func TranslateAppErrorsToGraphqlResponse(ctx context.Context, err error, fieldsMap map[string]string) bool {
	if len(graphql.GetErrors(ctx)) == 0 && err == nil {
		return false
	}
	var resErrors error
	resErrors = multierror.Append(resErrors, err)
	if merr, ok := resErrors.(*multierror.Error); ok {
		// Use merr.Errors
		for _, er := range merr.Errors {
			// if error is of type api error.
			appError := GetAppErrorObject(er)
			if fieldsMap[appError.Field] != "" {
				appError.Field = fieldsMap[appError.Field]
			}
			AppendAppErrorToGraphqlContext(ctx, appError)
		}
		return true
	} else {
		appError := GetAppErrorObject(err)
		AppendAppErrorToGraphqlContext(ctx, appError)
	}
	return true
}

func AppendAppErrorToGraphqlContext(ctx context.Context, err *appError) {
	gqlError := &gqlerror.Error{
		Path:       graphql.GetPath(ctx),
		Message:    getErrorMessageAndFieldForKey(err.Key),
		Extensions: getExtensionsMap(err.Status, err.Field),
	}
	graphql.AddError(ctx, gqlError)
}

func getExtensionsMap(status int, field string) map[string]interface{} {
	extensionsMap := make(map[string]interface{}, 0)
	updateMapIfNonEmpty(field, ExtensionFieldNameKey, extensionsMap)
	extensionsMap[ExtensionErrorStatus] = status
	return extensionsMap
}

func updateMapIfNonEmpty(field string, extensionKey string,
	extensionsMap map[string]interface{}) {
	if strings.TrimSpace(field) != "" {
		extensionsMap[extensionKey] = field
	}
}
