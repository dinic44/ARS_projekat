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
	w.Write([]byte(singleConfig.Id))
}

//Put New {id}
func (ts *Service) PutNewSingleConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]

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
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	rt.Id = id
	singleConfig, err := ts.store.PutNewSingleConfigVersion(rt)

	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Write([]byte(singleConfig.Id))
}

//Get Single {id}
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
/*func (ts *Service) GetAllSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllSingleConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}*/

//Delete Single
func (ts *Service) DeleteSingleConfigHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]
	_, err := ts.store.DeleteSingleConfig(id, ver)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
	}
}
