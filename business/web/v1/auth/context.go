package auth

import (
	"context"

	"github.com/google/uuid"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claims value from a context.Context.
// const claimKey ctxKey = 1

// key is used to store/retrieve a user value from a context.Context.
const userKey ctxKey = 1

const idTokenKey ctxKey = 2

const audienceKey ctxKey = 3

// =============================================================================

// SetUserID stores the user id from the request in the context.
func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return v
}

func SetIDToken(ctx context.Context, idToken string) context.Context {
	return context.WithValue(ctx, idTokenKey, idToken)
}

func GetIDToken(ctx context.Context) string {
	v, ok := ctx.Value(idTokenKey).(string)
	if !ok {
		return ""
	}
	return v
}

func SetAudience(ctx context.Context, audience string) context.Context {
	return context.WithValue(ctx, audienceKey, audience)
}

func GetAudience(ctx context.Context) string {
	v, ok := ctx.Value(audienceKey).(string)
	if !ok {
		return ""
	}
	return v
}
