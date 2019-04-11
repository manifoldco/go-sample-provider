package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-sample-provider/api"
	"github.com/manifoldco/go-sample-provider/db"
	"github.com/manifoldco/go-sample-provider/primitives"
)

var test bool
var graftonPath string

const (
	clientID     = "21jtaatqj8y5t0kctb2ejr6jev5w8"
	clientSecret = "3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw"
	connectorURL = "http://localhost:3001/v1/oauth/tokens"
)

var features = manifold.FeatureMap{
	"age":       2,
	"ready":     true,
	"hat_color": "red",
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.BoolVar(&test, "test", false, "run grafton test")
	flag.StringVar(&graftonPath, "grafton-path", "./grafton", "path of grafton bin")
}

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv, err := startServer(port, "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	defer srv.Close()

	if test {
		flags := map[string]interface{}{
			"product":        primitives.Product,
			"plan":           primitives.Plans[0],
			"new-plan":       primitives.Plans[1],
			"region":         primitives.Region,
			"client-id":      clientID,
			"client-secret":  clientSecret,
			"connector-port": "3001",
			"exclude":        []string{"plan-change", "resource-measures"},
			//"log":            "verbose",
			"plan-features": features,
		}

		err = testGrafton(graftonPath, port, flags)
		if err != nil {
			log.Print(err)
		}

	} else {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		log.Printf("Listening bear on port %s", port)

		<-quit
	}

	os.Remove("database.sqlite")
}

func startServer(port, database string) (*http.Server, error) {
	pkey, err := parsePubKey("masterkey.json")
	if err != nil {
		return nil, err
	}

	db, err := db.New("database.sqlite")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %s", err)
	}

	err = db.Register(&primitives.Bear{}, &primitives.Credential{})
	if err != nil {
		return nil, fmt.Errorf("failed to register database records: %s", err)
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
		return nil, fmt.Errorf("failed to initialize server: %s", err)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatalf("unable to start server %s", err)
		}
	}()

	return srv, nil
}

func testGrafton(path, port string, flags map[string]interface{}) error {
	args := []string{"test"}

	args = append(args, parseFlags(flags)...)

	args = append(args, "http://localhost:"+port)

	cmd := exec.Command(path, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func parseFlags(flags map[string]interface{}) []string {
	var args []string

	for k, value := range flags {
		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.Slice:
			values := value.([]string)
			for _, val := range values {
				args = append(args, fmt.Sprintf("--%s=%s", k, val))
			}
		case reflect.String:
			args = append(args, fmt.Sprintf("--%s=%s", k, value))
		case reflect.Map:
			values := value.(manifold.FeatureMap)
			if len(values) == 0 {
				continue
			}

			b, err := json.Marshal(values)
			if err != nil {
				panic(err)
			}

			arg := fmt.Sprintf(`--%s=%s`, k, b)

			args = append(args, arg)
		}
	}

	return args
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
