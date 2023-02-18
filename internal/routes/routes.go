package routes

import (
	"encoding/json"
	"net/http"

	conf "github.com/PKopel/mact/internal/config"
	mact "github.com/PKopel/mact/internal/json"
	"github.com/PKopel/mact/internal/utils"

	"github.com/gin-gonic/gin"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func setupEndopint(endpoint conf.EndpointConfig, host string) func(*gin.Context) {
	return func(c *gin.Context) {
		request, err := http.NewRequest(string(endpoint.Verb), host+endpoint.Path, c.Request.Body)
		if err != nil {
			panic(err)
		}
		copyHeader(request.Header, c.Request.Header)

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		var body mact.JSON
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			panic(err)
		}

		if utils.Contains(endpoint.StatusCodes, resp.StatusCode) {
			body = mact.ApplyChanges(body, endpoint.Changes)
		}

		c.JSON(resp.StatusCode, body)
	}
}

func SetupRouter(router *gin.Engine, config conf.MactConfig) {
	for _, service := range config.Services {
		for _, endpoint := range service.Endpoints {
			handlerFunc := setupEndopint(endpoint, service.Host)
			switch endpoint.Verb {
			case conf.GET:
				router.GET(endpoint.Path, handlerFunc)
			case conf.PUT:
				router.PUT(endpoint.Path, handlerFunc)
			case conf.POST:
				router.POST(endpoint.Path, handlerFunc)
			case conf.PATCH:
				router.PATCH(endpoint.Path, handlerFunc)
			case conf.DELETE:
				router.DELETE(endpoint.Path, handlerFunc)
			}
		}
	}
}
