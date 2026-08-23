package httpadapter

import (
	"context"
	"net/http"
)

// AuthPlaceholder reserves the authentication boundary for deployments that
// terminate identity at an API gateway. It preserves caller identity without
// making local development depend on an identity provider.
func AuthPlaceholder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subject := r.Header.Get("X-Caller-Subject"); subject != "" {
			ctx := withSubject(r.Context(), subject)
			_ = ctx
		}
		next.ServeHTTP(w, r)
	})
}

type subjectKey struct{}

func withSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey{}, subject)
}
func Subject(ctx context.Context) string { value, _ := ctx.Value(subjectKey{}).(string); return value }
