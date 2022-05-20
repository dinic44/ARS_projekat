package main

import (
	cs "ARS_projekat/configstore"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type Service struct {
	store *cs.ConfigStore
}

func (ts *Service) createPostHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := ts.store.Post(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
}

func (ts *Service) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

func (ts *Service) getPostHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, err := ts.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) delPostHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	msg, err := ts.store.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}