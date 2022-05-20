package main

import (
	cs "ARS_projekat/configstore"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	store, err := cs.New()
	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		store: store,
	}
	router.HandleFunc("/config", server.createSingleConfigHandler).Methods("POST") //Create Single
	router.HandleFunc("/configs", server.createGroupConfigHandler).Methods("POST") //Create Group

	router.HandleFunc("/config", server.getAllSingleConfigHandler).Methods("GET") //All Single
	router.HandleFunc("/configs", server.getAllGroupConfigHandler).Methods("GET") //All Group

	router.HandleFunc("/config/{id}", server.getSingleConfigHandler).Methods("GET") //Find One Single {id}
	router.HandleFunc("/configs/{id}", server.getGroupConfigHandler).Methods("GET") //Find One Group {id}

	router.HandleFunc("/config/{id}", server.deleteSingleConfigHandler).Methods("DELETE") //Delete Single {id}
	router.HandleFunc("/configs/{id}", server.deleteGroupConfigHandler).Methods("DELETE") //Delete Group {id}

	router.HandleFunc("/config/{id}", server.updateConfigHandler).Methods("PUT") //Update /{id}/

	router.HandleFunc("/config/{id}/{version}", server.getSingleVersionConfigHandler).Methods("GET")       //Find One Single {id}/{version}
	router.HandleFunc("/configs/{id}/{version}", server.getGroupVersionConfigHandler).Methods("GET")       //Find One Group {id}/{version}
	router.HandleFunc("/config/{id}/{version}", server.deleteSingleVersionConfigHandler).Methods("DELETE") //Delete Single {id}/{version}

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("Server Starting-----")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("Service Shutting Down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server Stopped-----")
}
