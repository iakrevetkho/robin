package main

import (
	// External
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	// Internal
	"github.com/iakrevetkho/robin/internal/config"
	resources "github.com/iakrevetkho/robin/internal/resources"
)

func sendFailedRequest(config config.Config, nc *nats.Conn) {
	log.Info("Send failed auth request")

	// Send request
	responseProto, err := nc.Request(config.NATS.Request.Subject, []byte("blabla"), 1*time.Second)
	if err != nil {
		log.Fatalf("Couldn't send request. %v", err)
	}

	// Parse response
	response := resources.Msg{}
	err = proto.Unmarshal(responseProto.Data, &response)
	if err != nil {
		log.Fatalf("Couldn't deserialize response. %v", err)
	}

	log.Infof("Response: %+v", response.GetPayload())
}