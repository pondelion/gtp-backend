package gcp_util

import (
	"context"
	"fmt"
	firebase_util "gtp/utils/gcp/firebase"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var app *firebase.App
var client *auth.Client

func App() (*firebase.App, error) {
	if app != nil {
		return app, nil
	}
	gcp_sa_filepath := os.Getenv("GCP_SA_CREDENTIAL_FILEPATH")
	if gcp_sa_filepath == "" {
		panic("GCP_SA_CREDENTIAL_FILEPATH must be set")
	}
	var err error
	app, err = firebase_util.CreateApp(gcp_sa_filepath)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}

func IPAuthClient() (*auth.Client, error) {
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
