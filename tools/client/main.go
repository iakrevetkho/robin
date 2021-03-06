package main

import (
	// External
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	// Internal
	"github.com/iakrevetkho/robin/internal/config"
	"github.com/iakrevetkho/robin/internal/helpers"
)

func main() {
	// Load app configuration
	config, err := config.LoadConfig("../../config.yml")
	if err != nil {
		log.Fatalf("Couldn't load app config. %v", err)
	}
	log.Infof("Loaded config: %+v", config)

	// Connect to a server
	nc, _ := nats.Connect(config.NATS.Hostname)
	defer nc.Close()

	// Init gin router
	r := gin.Default()
	RegisterRoutes(config, nc, r)

	// Load HTML templates
	r.LoadHTMLFiles("auth.html")

	go func() {
		err = r.Run(":9000")
		if err != nil {
			log.Fatalf("Couldn't run web server. %v", err)
		}
	}()

	// Example of error response from `robin`
	err = sendFailedRequest(config, nc)
	if err != nil {
		log.Fatalf("Couldn't send error auth request. %v", err)
	}

	// Wait any terminate signal
	signal := helpers.WaitTermSignals()
	log.Infof("Exit. Catched signal '%s'", signal.String())
}
