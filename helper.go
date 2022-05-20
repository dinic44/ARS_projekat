package main

import (
	cs "ARS_projekat/configstore"
	"encoding/json"
	"io"
	"net/http"
)

func decodeBody(r io.Reader) (*cs.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}