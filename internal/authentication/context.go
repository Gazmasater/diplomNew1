package authentication

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyDeleteCaller = contextKey("UserID")
)

// GetUserIDFromCtx gets the caller value from the context.
func GetUserIDFromCtx(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(ContextKeyDeleteCaller).(string)
	return caller, ok
}
