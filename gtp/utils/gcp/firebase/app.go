package firebase_util

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App

func App() (*firebase.App, error) {
	if app != nil {
		return app, nil
	}
	firebase_sa_filepath := os.Getenv("FIREBASE_SA_CREDENTIAL_FILEPATH")
	if firebase_sa_filepath == "" {
		panic("FIREBASE_SA_CREDENTIAL_FILEPATH must be set")
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(firebase_sa_filepath)
	fmt.Println(opt)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}
