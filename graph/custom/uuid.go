package custom

import (
	"context"
	"fmt"
	"io"
	"musicd/internal/errors"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrInvalidUUID struct {
	Field string
}

func (e ErrInvalidUUID) Code() errors.ErrorCode {
	return errors.ErrorCode(strings.Replace(is.ErrUUID.Code(), "_", ".", 1))
}

func (e ErrInvalidUUID) Message() string {
	return "Invalid uuid"
}

func (e ErrInvalidUUID) Error() string {
	return fmt.Sprintf("invalid uuid for field %s", e.Field)
}

func MarshalUUID(id uuid.UUID) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		_, _ = w.Write([]byte(fmt.Sprintf("%q", id.String())))
		return nil
	})
}

func UnmarshalUUID(ctx context.Context, value interface{}) (uuid.UUID, error) {
	rawUuid, ok := value.(string)
	if ok {
		id, err := uuid.FromString(rawUuid)
		if err == nil {
			return id, nil
		}
	}

	path := graphql.GetPathContext(ctx).Path()
	return uuid.Nil, gqlerror.WrapPath(path, ErrInvalidUUID{Field: path.String()})
}
