package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func parseID(r *http.Request) string {
	paths := strings.Split(r.URL.Path, "/")

	return paths[len(paths)-1]
}

func parseBody(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}

	return nil
}

func writeResponse(w http.ResponseWriter, content interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")

	msg, ok := content.(string)
	if ok {
		content = struct {
			Message string `json:"message"`
		}{
			Message: msg,
		}
	}

	b, err := json.Marshal(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(b)
}
