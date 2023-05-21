package pagination

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"musicd/internal/errors"
)

var ErrInvalidCursor = errors.New("invalid cursor")

func DecodeCursor[T any](ctx context.Context, rawCursor string) (T, error) {
	var cursor T
	rawJson, err := base64.StdEncoding.DecodeString(rawCursor)
	if err != nil {
		return cursor, ErrInvalidCursor
	}

	err = json.Unmarshal(rawJson, &cursor)
	if err != nil {
		return cursor, ErrInvalidCursor
	}

	return cursor, nil
}

func EncodeCursor[T any](cursor T) (string, error) {
	rawJson, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to encode cursor: %w", err)
	}
	return base64.StdEncoding.EncodeToString(rawJson), nil
}
