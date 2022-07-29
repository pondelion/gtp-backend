package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserKey struct{}
type App struct {
	certs map[string]string
	aud   string
}

var g_app *App

func IAPAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if g_app == nil {
				var err error
				g_app, err = NewApp()
				if err != nil {
					log.Println(err)
					fmt.Fprintln(w, "Failed to initialize IAP app.")
					return
				}
			}

			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
			if assertion == "" {
				fmt.Fprintln(w, "No Cloud IAP header found.")
				return
			}
			email, userID, err := ValidateAssertion(assertion, g_app.certs, g_app.aud)
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, "Could not validate assertion. Check app logs.")
				return
			}

			fmt.Fprintf(w, "Hello %s\n", email)
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, UserKey{}, userID)))
		})
	}
}

func ValidateAssertion(assertion string, certs map[string]string, aud string) (email string, userID string, err error) {
	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		keyID := token.Header["kid"].(string)

		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %q", token.Header["alg"])
		}

		cert := certs[keyID]
		return jwt.ParseECPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("could not extract claims (%T): %+v", token.Claims, token.Claims)
	}

	if claims["aud"].(string) != aud {
		return "", "", fmt.Errorf("mismatched audience. aud field %q does not match %q", claims["aud"], aud)
	}
	return claims["email"].(string), claims["sub"].(string), nil
}

func NewApp() (*App, error) {
	certs, err := Certificates()
	if err != nil {
		return nil, err
	}

	aud, err := Audience()
	if err != nil {
		return nil, err
	}

	a := &App{
		certs: certs,
		aud:   aud,
	}
	return a, nil
}

func Audience() (string, error) {
	gcp_project_number := os.Getenv("GCP_PROJECT_NUMBER")
	if gcp_project_number == "" {
		return "", fmt.Errorf("Environment variable GCP_PROJECT_NUMBER must be set")
	}

	gcp_project_id := os.Getenv("GCP_PROJECT_ID")
	if gcp_project_id == "" {
		return "", fmt.Errorf("Environment variable GCP_PROJECT_ID must be set")
	}

	return "/projects/" + gcp_project_number + "/apps/" + gcp_project_id, nil
}

func Certificates() (map[string]string, error) {
	const url = "https://www.gstatic.com/iap/verify/public_key"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Get: %v", err)
	}

	var certs map[string]string
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&certs); err != nil {
		return nil, fmt.Errorf("Decode: %v", err)
	}

	return certs, nil
}
