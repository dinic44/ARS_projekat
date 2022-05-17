package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
)

type Service struct {
	Data map[string][]*Config `json:"data"`
}

//Create Single
func (ts *Service) createSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil || len(rt) > 1 {
		http.Error(w, "Invalid Format is > 1", http.StatusBadRequest)
		return
	}

	id := createId()
	ts.Data[id] = rt
	renderJSON(w, rt)
	w.Write([]byte(id))
}

//Create Group
func (ts *Service) createGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil || len(rt) == 1 {
		http.Error(w, "Invalid Format is == 1", http.StatusBadRequest)
		return
	}

	id := createId()
	ts.Data[id] = rt
	renderJSON(w, rt)
	w.Write([]byte(id))
}

//Get All Single
func (ts *Service) getAllSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.Data {
		if len(v) == 1 {
			allTasks = append(allTasks, v...)
		}
	}

	renderJSON(w, allTasks)
}

//Get All Group
func (ts *Service) getAllGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.Data {
		if len(v) > 1 {
			allTasks = append(allTasks, v...)
		}
	}

	renderJSON(w, allTasks)
}

//Get/{id} Single
func (ts *Service) getSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.Data[id]
	if !ok || len(task) > 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Get/{id} Group
func (ts *Service) getGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.Data[id]
	if !ok || len(task) == 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Delete/{id} Group
func (ts *Service) deleteGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if value, ok := ts.Data[id]; ok && len(value) > 1 {
		delete(ts.Data, id)
		renderJSON(w, value)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

//Delete/{id} Single
func (ts *Service) deleteSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if value, ok := ts.Data[id]; ok && len(value) == 1 {
		delete(ts.Data, id)
		renderJSON(w, value)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

//Put/{id}
func (ts *Service) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(req)["id"]
	task, err2 := ts.Data[id]
	if !err2 {
		err2 := errors.New("key not found")
		http.Error(w, err2.Error(), http.StatusNotFound)
		return
	}

	for _, config := range rt {
		task = append(task, config)
	}

	ts.Data[id] = task
	renderJSON(w, task)
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

func decodeBody(r io.Reader) ([]*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt []*Config

	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func createId() string {
	return uuid.New().String()
}
