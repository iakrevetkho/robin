package connectorsnats

import (
	// External
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	// Internal
	apirouters "github.com/iakrevetkho/robin/internal/api/routers"
	resources "github.com/iakrevetkho/robin/internal/resources"
)

func processMsg(msg *nats.Msg) {
	// Check that msg is request
	if msg.Reply == "" {
		log.Errorf("Await request, but receive msg without reply: %+v", msg)
		// TODO Send error response
		return
	}

	// Parse message protobuf
	protoMsg := resources.Msg{}
	if err := proto.Unmarshal(msg.Data, &protoMsg); err != nil {
		log.Errorf("Failed to parse proto msg: %v", err)
		// TODO Send error response
		return
	}

	// Go to msg router
	apirouters.RouteMsg(&protoMsg)
}