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
	var payload ProvisionRequest

	err := parseBody(r, &payload)
	if err != nil {
		writeResponse(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if payload.Product != primitives.Product {
		writeResponse(w, "invalid product", http.StatusBadRequest)
		return
	}

	if !primitives.ValidPlan(payload.Plan) {
		writeResponse(w, "invalid plan", http.StatusBadRequest)
		return
	}

	if payload.Region != primitives.Region {
		writeResponse(w, "invalid region", http.StatusBadRequest)
		return
	}

	var bear primitives.Bear

	err = db.FindBy("manifold_id", payload.ResourceID, &bear)
	if err == nil {
		if bear.Plan != payload.Plan {
			writeResponse(w, "conflict", http.StatusConflict)
			return
		}

		writeResponse(w, "bear was already created", http.StatusCreated)
		return
	}

	bear = primitives.Bear{
		Name:       "test",
		ManifoldID: payload.ResourceID,
		Plan:       payload.Plan,
		Age:        payload.Features.GetInt("age"),
		Ready:      payload.Features.GetBool("ready"),
		HatColor:   payload.Features.GetString("hat_color"),
	}

	err = db.Create(&bear)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, "bear was created", http.StatusCreated)
}

func deprovision(w http.ResponseWriter, r *http.Request, db primitives.Database) {
	id := parseID(r)

	var bear primitives.Bear

	err := db.FindBy("manifold_id", id, &bear)
	if err != nil {
		writeResponse(w, "bear was not found", http.StatusNotFound)
		return
	}

	err = db.Delete(&bear)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, "bear was deleted", http.StatusNoContent)
}
