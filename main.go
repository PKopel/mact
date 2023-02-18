package main

import (
	"flag"

	conf "github.com/PKopel/mact/internal/config"
	"github.com/PKopel/mact/internal/routes"
	"github.com/gin-gonic/gin"
)

var configFile = flag.String("config", "config.yaml", "Path to config file")
var serverPort = flag.String("port", "8000", "Port to listen on")
var debugMode = flag.Bool("debug", false, "Enable debug mode")

func main() {
	flag.Parse()
	if !*debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	config := conf.ReadConfig(*configFile)
	routes.SetupRouter(router, config)
	router.Run(":" + *serverPort)
}
