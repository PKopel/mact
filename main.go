package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/PKopel/mact/internal/routes"
	"github.com/PKopel/mact/types"
	"github.com/gorilla/mux"
)

var configFile = flag.String("config", "config.yaml", "Path to config file")
var serverPort = flag.String("port", "8000", "Port to listen on")

func main() {
	flag.Parse()
	log.SetPrefix("[MA/CT] ")
	router := mux.NewRouter()
	config := types.ReadConfig(*configFile)
	routes.SetupRouter(router, config)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + *serverPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
