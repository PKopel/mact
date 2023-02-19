package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	conf "github.com/PKopel/mact/internal/config"
	mact "github.com/PKopel/mact/internal/json"
	"github.com/PKopel/mact/internal/utils"
	"github.com/gorilla/mux"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func setupEndopint(endpoint conf.EndpointConfig, host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server := host + r.URL.Path
		request, err := http.NewRequest(string(endpoint.Verb), server, r.Body)
		if err != nil {
			log.Fatalf("Error while creating request: %v", err)
		}
		copyHeader(request.Header, r.Header)

		log.Printf("Sending request to %v", server)
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatalf("Error while making request: %v", err)
		}

		var body mact.JSON
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			log.Fatalf("Error while decoding response: %v", err)
		}

		if utils.Contains(endpoint.StatusCodes, resp.StatusCode) {
			body = mact.ApplyChanges(body, endpoint.Changes)
		}

		log.Printf("Writing response body: %v", body)

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(body)
	}
}

func setupIgnored(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server := host + r.URL.Path
		request, err := http.NewRequest(r.Method, host+r.URL.Path, r.Body)
		if err != nil {
			log.Fatalf("Error while creating request: %v", err)
		}
		copyHeader(request.Header, r.Header)

		log.Printf("Sending request to %v", server)

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatalf("Error while making request: %v", err)
		}

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func SetupRouter(router *mux.Router, config conf.MactConfig) {

	for _, service := range config.Services {
		log.Printf("Setting up service %v", service.Host)
		// setup endpoints to be processed
		for _, endpoint := range service.Endpoints {
			handlerFunc := setupEndopint(endpoint, service.Host)
			router.Path(endpoint.Path).HandlerFunc(handlerFunc).Methods(string(endpoint.Verb))
			log.Printf("Setting up endpoint %v", endpoint.Path)
		}
		// setup endpoints to be ignored
		handlerFunc := setupIgnored(service.Host)
		router.PathPrefix("/").HandlerFunc(handlerFunc)

	}
}
