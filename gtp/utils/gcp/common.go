package gcp_util

import "os"

func GCPDeployed() (bool, error) {
	gcp_deploy := os.Getenv("GCP_DEPLOY")
	if gcp_deploy == "" {
		return false, nil
	}
	if gcp_deploy == "false" {
		return false, nil
	}
	return true, nil
}
