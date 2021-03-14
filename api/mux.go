package api

import (
	"net/http"

	"github.com/manifoldco/go-sample-provider/primitives"
	"github.com/manifoldco/go-signature"
	"golang.org/x/oauth2"
)

// NewServeMux returns a new http ServeMux that handles api calls.
func NewServeMux(db primitives.Database, pkey string, oauth oauth2.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	verifier, err := signature.NewVerifier(pkey)
	if err != nil {
		return nil, err
	}

	mux.Handle("/v1/resources/", verifier.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			provision(w, r, db)
		case "PATCH":
			changePlan(w, r, db)
		case "DELETE":
			deprovision(w, r, db)
		default:
			writeResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/v1/credentials/", verifier.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			addCredentials(w, r, db)
		case "DELETE":
			deleteCredentials(w, r, db)
		default:
			writeResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/v1/sso", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			sso(w, r, db, oauth)
		default:
			writeResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux, nil
}
