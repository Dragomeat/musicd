package auth

import "musicd/internal/errors"

var ErrUnathenticated = errors.NewUserError("auth.unauthenticated", "Unauthenticated")
