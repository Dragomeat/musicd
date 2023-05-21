package http

import (
	"context"
	"encoding/json"
	"fmt"
	"musicd/internal/errors"
	oHttp "net/http"
)

type ErrorTransformer func(ctx context.Context, err error) error

type ErrorHandler struct {
	errorHandler *errors.ErrorHandler
	transformers []ErrorTransformer
}

func NewErrorHandler(errorHandler *errors.ErrorHandler) *ErrorHandler {
	return &ErrorHandler{
		errorHandler: errorHandler,
	}
}

func (m *ErrorHandler) AttachTransformer(transformer ErrorTransformer) {
	m.transformers = append(m.transformers, transformer)
}

func (m *ErrorHandler) Handle(w oHttp.ResponseWriter, req *oHttp.Request, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.Join(err, fmt.Errorf("panic: %v", p))
			m.errorHandler.Handle(req.Context(), err)
			m.sendErrors(w, req, []map[string]interface{}{m.formatError(err)})
		}
	}()

	ctx := req.Context()
	for _, transformer := range m.transformers {
		err = transformer(ctx, err)
	}

	// if !ignoreHandling {
	m.errorHandler.Handle(ctx, err)
	// }

	m.sendErrors(w, req, []map[string]interface{}{m.formatError(err)})
}

func (m *ErrorHandler) sendErrors(w oHttp.ResponseWriter, req *oHttp.Request, errors []map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(oHttp.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(
		map[string]interface{}{
			"data":   nil,
			"errors": errors,
		},
	)

	if err == nil {
		return
	}

	m.errorHandler.Handle(
		req.Context(),
		fmt.Errorf("can`t send error to client: %w", err),
	)
}

func (m *ErrorHandler) formatError(err error) map[string]interface{} {
	code := errors.Code(err)
	message := errors.Message(err)
	extra := errors.Extra(err)

	return map[string]interface{}{
		"code":    code,
		"message": message,
		"extra":   extra,
	}
}

func (m *ErrorHandler) Wrap(handler RequestHandler) oHttp.Handler {
	return oHttp.HandlerFunc(
		func(w oHttp.ResponseWriter, req *oHttp.Request) {
			defer func() {
				if p := recover(); p != nil {
					m.Handle(w, req, fmt.Errorf("panic: %v", p))
				}
			}()

			err := handler.Handle(w, req)
			if err != nil {
				m.Handle(w, req, err)
			}
		},
	)
}
