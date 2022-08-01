package firebase_util

import (
	"context"
	"fmt"
	gcp_util "gtp/utils/gcp"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var g_app *firebase.App

func App() (*firebase.App, error) {
	if g_app != nil {
		return g_app, nil
	}
	var err error
	gcp_deployed, _ := gcp_util.GCPDeployed()
	if gcp_deployed {
		g_app, err = CreateApp(nil)
	} else {
		firebase_sa_filepath := os.Getenv("FIREBASE_SA_CREDENTIAL_FILEPATH")
		if firebase_sa_filepath == "" {
			panic("FIREBASE_SA_CREDENTIAL_FILEPATH must be set")
		}
		g_app, err = CreateApp(&firebase_sa_filepath)
	}
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	return g_app, nil
}

func CreateApp(sa_credential_filepath *string) (*firebase.App, error) {
	ctx := context.Background()
	var err error
	var app *firebase.App
	if sa_credential_filepath == nil {
		app, err = firebase.NewApp(ctx, nil)
	} else {
		opt := option.WithCredentialsFile(*sa_credential_filepath)
		app, err = firebase.NewApp(ctx, nil, opt)
	}
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}
