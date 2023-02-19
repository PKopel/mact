package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	conf "github.com/PKopel/mact/internal/config"
	"github.com/PKopel/mact/internal/routes"
	"github.com/gorilla/mux"
)

var configFile = flag.String("config", "config.yaml", "Path to config file")
var serverPort = flag.String("port", "8000", "Port to listen on")

func main() {
	flag.Parse()
	log.SetPrefix("[MA/CT]")
	router := mux.NewRouter()
	config := conf.ReadConfig(*configFile)
	routes.SetupRouter(router, config)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + *serverPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
