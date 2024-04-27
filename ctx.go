package authz

import "context"

type tokenKey struct{}
type sessionKey struct{}

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

func GetToken(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(tokenKey{}).(string)
	return value, ok
}

func SetSession[ID IdType](ctx context.Context, session *Session[ID]) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}

func GetSession[ID IdType](ctx context.Context) (*Session[ID], bool) {
	value, ok := ctx.Value(sessionKey{}).(*Session[ID])
	return value, ok
}
