package main

import (
	cs "ARS_projekat/configstore"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"net/url"
)

type Service struct {
	store *cs.ConfigStore
}

//Create Single
func (ts *Service) CreateSingleConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

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

	rt, err := decodeBodySingle(req.Body)
	if err != nil || rt.Version == "" || rt.Entries == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Already Sent", http.StatusBadRequest)
		return
	}

	singleConfig, err := ts.store.CreateSingleConfig(rt)
	reqId := ""
	if err == nil {
		reqId = ts.store.SaveRequestId()
	}
	w.Write([]byte(singleConfig.Id))
	w.Write([]byte("\n\nIdempotency Key: " + reqId))
}

//Put New {id}
func (ts *Service) PutNewSingleConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodySingle(req.Body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	rt.Id = id
	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Already Sent", http.StatusBadRequest)
	}
	singleConfig, err := ts.store.PutNewSingleConfigVersion(rt)

	if err != nil {
		http.Error(w, "Version Already Exists", http.StatusBadRequest)
		return
	}

	reqId := ""
	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	w.Write([]byte(singleConfig.Id))
	w.Write([]byte("\n\nIdempotency Key: " + reqId))
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
	renderJSON(w, task, "")
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
	renderJSON(w, task, "")
}

//Delete Single
func (ts *Service) DeleteSingleConfigHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]
	_, err := ts.store.DeleteSingleConfig(id, ver)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
	}
}

//Create Group
func (ts *Service) CreateGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

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

	rt, err := decodeBodyGroup(req.Body)
	if err != nil || rt.Version == "" || rt.GroupConfig == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Already Sent", http.StatusBadRequest)
		return
	}

	groupConfig, err := ts.store.CreateGroupConfig(rt)

	reqId := ""
	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	w.Write([]byte(groupConfig.Id))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

//Put new{id}
func (ts *Service) PutNewGroupConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyGroup(req.Body)
	if err != nil || rt.Version == "" || rt.GroupConfig == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Already Sent", http.StatusBadRequest)
		return
	}

	rt.Id = id
	groupConfig, err := ts.store.PutNewGroupConfigVersion(rt)

	reqId := ""
	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	if err != nil {
		http.Error(w, "Already Exists", http.StatusBadRequest)
		return
	}

	w.Write([]byte(groupConfig.Id))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

//Find One from Group {id}/{version}
func (ts *Service) GetGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	task, ok := ts.store.GetGroupConfig(id, ver)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task, "")
}

//Find Single Key:Value from version of a group
func (ts *Service) GetSingleConfigFromGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	req.ParseForm()
	params := url.Values.Encode(req.Form)
	label, err := ts.store.FindSingleInGroup(id, ver, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, label, "")
}

//Delete {id}/{ver}
func (ts *Service) DeleteGroupConfigHandler(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	ver := mux.Vars(request)["ver"]
	err := ts.store.DeleteGroupConfig(id, ver)
	if err != nil {
		http.Error(writer, "error", http.StatusBadRequest)
	}
}
