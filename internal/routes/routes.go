package routes

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
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

func setupEndopint(endpoint conf.EndpointConfig, service conf.ServiceConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server := service.Host + r.URL.Path + "?" + r.URL.RawQuery
		request, err := http.NewRequest(string(endpoint.Verb), server, r.Body)
		if err != nil {
			log.Fatalf("Error while creating request: %v", err)
		}
		copyHeader(request.Header, r.Header)

		transport := http.DefaultTransport
		if service.TrustAllCerts {
			log.Printf("Trusting all certificates")
			transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client := &http.Client{
			Transport: transport,
		}

		log.Printf("Sending request to %v", server)
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

		bytes, err := json.Marshal(body)
		if err != nil {
			log.Fatalf("Error while encoding response: %v", err)
		}
		log.Printf("Writing response body: %v", string(bytes))

		copyHeader(w.Header(), resp.Header)
		w.Header().Set("Content-Length", fmt.Sprint(len(bytes)))
		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(bytes)
		if err != nil {
			log.Fatalf("Error while writing response: %v", err)
		}
	}
}

func setupIgnored(service conf.ServiceConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server := service.Host + r.URL.Path + "?" + r.URL.RawQuery
		request, err := http.NewRequest(r.Method, server, r.Body)
		if err != nil {
			log.Fatalf("Error while creating request: %v", err)
		}
		copyHeader(request.Header, r.Header)

		transport := http.DefaultTransport
		if service.TrustAllCerts {
			log.Printf("Trusting all certificates")
			transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client := &http.Client{
			Transport: transport,
		}

		log.Printf("Sending request to %v", server)
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
			handlerFunc := setupEndopint(endpoint, service)
			router.Path(endpoint.Path).HandlerFunc(handlerFunc).Methods(string(endpoint.Verb))
			log.Printf("Setting up endpoint %v", endpoint.Path)
		}
		// setup endpoints to be ignored
		handlerFunc := setupIgnored(service)
		router.PathPrefix("/").HandlerFunc(handlerFunc)

	}
}
