package api

import (
	"net/http"

	"github.com/manifoldco/go-sample-provider/primitives"
)

// ProvisionRequest is the expected request from Manifold when it's trying to provision a new resource.
type ProvisionRequest struct {
	ResourceID string `json:"id"`
	Product    string
	Plan       string
	Region     string
	Features   FeatureMap
}

func provision(w http.ResponseWriter, r *http.Request, db primitives.Database) {
}

func deprovision(w http.ResponseWriter, r *http.Request, db primitives.Database) {
}
