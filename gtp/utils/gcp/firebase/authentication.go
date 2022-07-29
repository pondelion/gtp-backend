package firebase_util

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
)

var client *auth.Client

func AuthClient() (*auth.Client, error) {
	if client != nil {
		return client, nil
	}
	app, err := App()
	ctx := context.Background()
	client, err = app.Auth(ctx)
	if err != nil {
		fmt.Errorf("error initializing firebase auth: %v", err)
	}
	return client, nil
}
