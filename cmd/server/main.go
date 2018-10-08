package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/manifoldco/go-sample-provider/api"
	"github.com/manifoldco/go-sample-provider/db"
	"github.com/manifoldco/go-sample-provider/primitives"
	"golang.org/x/oauth2"
)

const (
	clientID     = "21jtaatqj8y5t0kctb2ejr6jev5w8"
	clientSecret = "3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw"
	connectorURL = "http://localhost:3001/v1/oauth/tokens"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pkey, err := parsePubKey("masterkey.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Register(&primitives.Bear{}, &primitives.Credential{})
	if err != nil {
		log.Fatal(err)
	}

	oauth := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: connectorURL,
		},
	}

	mux, err := api.NewServeMux(db, pkey, oauth)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server %s", err)
	}
}

func parsePubKey(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var sf struct {
		PublicKey string `json:"public_key"`
	}

	err = json.Unmarshal(b, &sf)
	if err != nil {
		return "", err
	}

	return sf.PublicKey, nil
}
