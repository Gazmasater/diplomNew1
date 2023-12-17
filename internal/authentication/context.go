package authentication

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyDeleteCaller = contextKey("UserID")
)

func GetUserIDFromCtx(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(ContextKeyDeleteCaller).(string)
	return caller, ok
}
