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

	router.HandleFunc("/singleConfig/", server.CreateSingleConfigHandler).Methods("POST")            //Create Single
	router.HandleFunc("/singleConfig/{id}", server.PutNewSingleConfigVersionHandler).Methods("POST") //Update Single
	router.HandleFunc("/singleConfig/{id}/", server.GetSingleConfigVersionHandler).Methods("GET")    //Find One Single{id}
	router.HandleFunc("/singleConfig/{id}/{ver}", server.FindSingleConfigHandler).Methods("GET")     //Find One Single{id}/{version}
	/*	router.HandleFunc("/singleConfigs", server.GetAllSingleConfigHandler).Methods("GET")  */       //Find All Single
	router.HandleFunc("/singleConfig/{id}/{ver}", server.DeleteSingleConfigHandler).Methods("DELETE") //Delete Single

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
