package custom

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalTimestamp(time time.Time) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		_, _ = w.Write([]byte(strconv.FormatInt(time.Unix(), 10)))
		return nil
	})
}

func UnmarshalTimestamp(ctx context.Context, value interface{}) (time.Time, error) {
	panic(fmt.Errorf("unmarshal timestamp not implemented"))
}

func MarshalMilliTimestamp(time time.Time) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		_, _ = w.Write([]byte(strconv.FormatInt(time.UnixMilli(), 10)))
		return nil
	})
}
