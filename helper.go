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

	var singleConfig *cs.SingleConfig
	if err := dec.Decode(&singleConfig); err != nil {
		return nil, err
	}
	return singleConfig, nil
}

func decodeBodyGroup(r io.Reader) (*cs.GroupConfig, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var groupConfig *cs.GroupConfig
	if err := dec.Decode(&groupConfig); err != nil {
		return nil, err
	}
	return groupConfig, nil
}

func renderJSON(w http.ResponseWriter, v interface{}, id string) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
