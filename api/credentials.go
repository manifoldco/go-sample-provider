package api

import (
	"net/http"

	"github.com/manifoldco/go-sample-provider/primitives"
)

// CredentialsRequest is the expected request when Manifold is trying to create a new credentials set
type CredentialsRequest struct {
	CredentialID string `json:"id"`
	ResourceID   string `json:"resource_id"`
}

// CredentialResponse is the expected response from us when Manifold is trying to createa new credentials set
type CredentialResponse struct {
	Message     string
	Credentials map[string]string
}

func addCredentials(w http.ResponseWriter, r *http.Request, db primitives.Database) {
}

func deleteCredentials(w http.ResponseWriter, r *http.Request, db primitives.Database) {
}
