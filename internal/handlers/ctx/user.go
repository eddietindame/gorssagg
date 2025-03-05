package ctx

import (
	"context"

	"github.com/google/uuid"
)

type UserContext struct {
	UserID   uuid.UUID
	Username string
	Email    string
}

type key int

var userKey key

func NewContextWithUser(ctx context.Context, u UserContext) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func GetUserFromContext(ctx context.Context) (UserContext, bool) {
	user, ok := ctx.Value(userKey).(UserContext)
	return user, ok
}
