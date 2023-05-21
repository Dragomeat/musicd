package validator

import (
	"context"
	"fmt"
	"musicd/internal/errors"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrInvalidInput struct {
	Fields []InvalidField
}

type InvalidField struct {
	Name    string
	Message string
	Code    string
}

func (e ErrInvalidInput) Code() errors.ErrorCode {
	return "app.invalidInput"
}

func (e ErrInvalidInput) Message() string {
	return "Invalid input"
}

func (e ErrInvalidInput) Extra() map[string]any {
	fields := make(map[string]map[string]string)
	for _, field := range e.Fields {
		fields[field.Name] = map[string]string{
			"code":    field.Code,
			"message": field.Message,
		}
	}

	return map[string]any{
		"fields": fields,
	}
}

func (e ErrInvalidInput) Error() string {
	var fields []string
	for _, field := range e.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s: %s", field.Name, field.Code, field.Message))
	}
	return fmt.Sprintf("invalid input (%s)", strings.Join(fields, ", "))
}

func NewErrInvalidInputFromOzzo(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var fieldErrors validation.Errors
	if !errors.As(err, &fieldErrors) {
		return err
	}

	basePath := graphql.GetPath(ctx).String()
	fields := make([]InvalidField, 0, len(fieldErrors))
	for key, fieldErr := range fieldErrors {
		fieldErr, ok := fieldErr.(validation.ErrorObject)
		if !ok {
			continue
		}

		field := basePath + "." + strings.ToLower(key)
		fields = append(fields, NewInvalidFieldFromOzzo(field, fieldErr))
	}

	if len(fields) == 0 {
		return err
	}

	return ErrInvalidInput{Fields: fields}
}

func NewInvalidFieldFromOzzo(field string, err validation.Error) InvalidField {
	return InvalidField{
		Name:    field,
		Code:    strings.Replace(err.Code(), "_", ".", 1),
		Message: err.Error(),
	}
}

func Path(ctx context.Context, path ...string) string {
	path = append([]string{graphql.GetPath(ctx).String()}, path...)
	return strings.Join(path, ".")
}

type Validatable interface {
	Validate(ctx context.Context) error
}
