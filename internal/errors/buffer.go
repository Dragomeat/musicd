package errors

import "context"

type contextKey string

const (
	bufferKey contextKey = "errorBuffer"
)

type buffer struct {
	errors []error
}

func Init(ctx context.Context) context.Context {
	b := &buffer{}
	return withBuffer(ctx, b)
}

func PushError(ctx context.Context, err error) context.Context {
	b := getBuffer(ctx)
	if b == nil {
		b = &buffer{}
		ctx = withBuffer(ctx, b)
	}

	b.errors = append(b.errors, err)

	return ctx
}

func PopErrors(ctx context.Context) []error {
	if buf := getBuffer(ctx); buf != nil {
		errors := buf.errors
		buf.errors = nil
		return errors
	}

	return nil
}

func getBuffer(ctx context.Context) *buffer {
	if buf, ok := ctx.Value(bufferKey).(*buffer); ok {
		return buf
	}
	return nil
}

func withBuffer(ctx context.Context, buffer *buffer) context.Context {
	return context.WithValue(ctx, bufferKey, buffer)
}
