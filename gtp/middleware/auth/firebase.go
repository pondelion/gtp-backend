package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase_util "gtp/utils/gcp/firebase"
)

type userKey struct{}

func FirebaseAuth() func(http.Handler) http.Handler {
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
			auth_client, err := firebase_util.AuthClient()
			token, err := auth_client.VerifyIDToken(ctx, idToken)
			if err != nil {
				http.Error(w, "Firebase auth verification failed. Invalid token", http.StatusForbidden)
				return
			}
			// fmt.Println(token)

			// user, err := db.GetUserByDigestUID(hash(token.UID))
			// if err != nil {
			// 	http.Error(w, "Invalid token", http.StatusForbidden)
			// 	return
			// }

			// and call the next with our new context
			// r = r.WithContext(ctx)
			// next.ServeHTTP(w, r)
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, userKey{}, token.UID)))
		})
	}
}
