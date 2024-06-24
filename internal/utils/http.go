package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ParseRequest(w http.ResponseWriter, r *http.Request, target any) error {
	var buf bytes.Buffer
	defer r.Body.Close()

	if _, err := buf.ReadFrom(r.Body); err != nil {
		Log.Infof("Could not get body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &target); err != nil {
		Log.Infof("Could not parse request JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	return nil
}
