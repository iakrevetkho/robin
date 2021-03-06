package apicontrollers

import (
	// External

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	// Internal
	helpers "github.com/iakrevetkho/robin/internal/api/helpers"
	apiresources "github.com/iakrevetkho/robin/internal/api/resources"
	apirouters "github.com/iakrevetkho/robin/internal/api/routers"
	proto_resources "github.com/iakrevetkho/robin/internal/proto_resources"
)

// Method for processing messages from NATS broker
func NatsMessage(controllerData apiresources.ControllerData, msg *nats.Msg) (err error) {
	// Check that msg is request
	if msg.Reply == "" {
		err = helpers.NatsErrorResponse(msg, &proto_resources.UUID{}, "Await request, but receive msg without reply: %+v", msg)
		return
	}

	// Parse message protobuf
	protoMsg := proto_resources.Msg{}
	if err = proto.Unmarshal(msg.Data, &protoMsg); err != nil {
		err = helpers.NatsErrorResponse(msg, &proto_resources.UUID{}, "Failed to parse proto msg '%s': %v", msg.Data, err)
		return
	}

	// Go to msg router
	response, err := apirouters.RouteMsg(controllerData, &protoMsg)
	if err != nil {
		err = helpers.NatsErrorResponse(msg, &proto_resources.UUID{}, "Couldn't process msg '%+v': %v", protoMsg.GetPayload(), err)
		return
	}

	// Serialize response
	responseProto, err := proto.Marshal(response)
	if err != nil {
		err = helpers.NatsErrorResponse(msg, &proto_resources.UUID{}, "Couldn't serialize proto response. %v", err)
		return
	}

	// Send response
	err = msg.Respond(responseProto)

	return
}
