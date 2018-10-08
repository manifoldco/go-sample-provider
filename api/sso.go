package api

import (
	"net/http"

	"github.com/manifoldco/go-sample-provider/primitives"
	"golang.org/x/oauth2"
)

func sso(w http.ResponseWriter, r *http.Request, db primitives.Database, config oauth2.Config) {
}
