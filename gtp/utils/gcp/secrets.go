package gcp_util

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetSecret(name string, projectId string, saFilepath *string) (string, error) {
	// Create the client.
	ctx := context.Background()
	var err error
	var client *secretmanager.Client
	if saFilepath == nil {
		client, err = secretmanager.NewClient(ctx)
	} else {
		client, err = secretmanager.NewClient(ctx, option.WithCredentialsFile(*saFilepath))
	}
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectId + "/secrets/" + name + "/versions/latest",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}
	fmt.Printf("retrieved payload for: %s\n", result.Name)
	fmt.Println(string(result.Payload.Data))
	return string(result.Payload.Data), nil
}
