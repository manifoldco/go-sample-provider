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
	var payload CredentialsRequest

	err := parseBody(r, &payload)
	if err != nil {
		writeResponse(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var cred primitives.Credential
	err = db.FindBy("manifold_id", payload.CredentialID, &cred)
	if err == nil {
		res := CredentialResponse{
			Message: "credentials created",
			Credentials: map[string]string{
				"SECRET": cred.Secret,
			},
		}
		writeResponse(w, res, http.StatusCreated)
		return
	}

	var bear primitives.Bear

	err = db.FindBy("manifold_id", payload.ResourceID, &bear)
	if err != nil {
		writeResponse(w, "bear not found", http.StatusNotFound)
		return
	}

	cred = primitives.Credential{
		BearID:     bear.ID,
		Secret:     primitives.CredentialSecret(),
		ManifoldID: payload.CredentialID,
	}

	err = db.Create(&cred)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CredentialResponse{
		Message: "credentials created",
		Credentials: map[string]string{
			"SECRET": cred.Secret,
		},
	}

	writeResponse(w, res, http.StatusCreated)
}

func deleteCredentials(w http.ResponseWriter, r *http.Request, db primitives.Database) {
	id := parseID(r)

	var cred primitives.Credential

	err := db.FindBy("manifold_id", id, &cred)
	if err != nil {
		writeResponse(w, "credentials were not found", http.StatusNotFound)
		return
	}

	err = db.Delete(&cred)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, "credentials were deleted", http.StatusNoContent)
}
