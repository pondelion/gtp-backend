package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gcp_util "gtp/utils/gcp"
)

type GCIPUser struct{}

func GCIPAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Middleware auth")
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "Authorization: Bearer [TOKEN] not found in header", http.StatusForbidden)
				return
			}
			ctx := r.Context()
			idToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
			// fmt.Println(idToken)
			//validate jwt token
			auth_client, err := gcp_util.IPAuthClient()
			token, err := auth_client.VerifyIDToken(ctx, idToken)
			if err != nil {
				http.Error(w, "GCP Identity Provider auth verification failed. Invalid token", http.StatusForbidden)
				return
			}
			fmt.Println(token)
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, GCIPUser{}, token.UID)))
		})
	}
}
