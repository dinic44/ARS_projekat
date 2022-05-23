package main

import (
	cs "ARS_projekat/configstore"
	"encoding/json"
	"io"
	"net/http"
)

func decodeBodySingle(r io.Reader) (*cs.SingleConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *cs.SingleConfig
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeBodyGroup(r io.Reader) (*cs.GroupConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var group *cs.GroupConfig
	if err := dec.Decode(&group); err != nil {
		return nil, err
	}
	return group, nil
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
