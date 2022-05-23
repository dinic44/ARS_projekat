package main

import (
	cs "ARS_projekat/configstore"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type Service struct {
	store *cs.ConfigStore
}

//Create Single
func (ts *Service) CreateSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
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

	rt, err := decodeBodySingle(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	singleConfig, err := ts.store.CreateSingleConfig(rt)
	log.Default().Println(singleConfig)
	w.Write([]byte(singleConfig.Id))
}

func (ts *Service) GetSingleConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindSingleConfigVersion(id)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Find One Single/{id}/{version}
func (ts *Service) FindSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindSingleConfig(id, ver)
	if ok != nil {
		err := errors.New("not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Find All Single
func (ts *Service) GetAllSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllSingleConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}
