package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/manifoldco/go-sample-provider/primitives"
	"golang.org/x/oauth2"
)

func sso(w http.ResponseWriter, r *http.Request, db primitives.Database, config oauth2.Config) {
	query := r.URL.Query()

	code := query.Get("code")
	id := query.Get("resource_id")

	ctx := context.Background()

	token, err := config.Exchange(ctx, code)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid() {
		writeResponse(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var bear primitives.Bear

	err = db.FindBy("manifold_id", id, &bear)
	if err != nil {
		writeResponse(w, "bear was not found", http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("/dashboard?id=%d", bear.ID)

	http.Redirect(w, r, url, http.StatusFound)
}
