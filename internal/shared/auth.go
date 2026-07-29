package shared

import (
	"context"
	"net/http"
)

type userIDContextKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

func CurrentUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDContextKey{}).(string)
	return userID, ok && userID != ""
}
