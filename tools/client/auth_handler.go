package main

import (
	// External
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	// Internal
	"github.com/iakrevetkho/robin/internal/config"
	proto_resources "github.com/iakrevetkho/robin/internal/proto_resources"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func AuthHandler(config config.Config, nc *nats.Conn, c *gin.Context) {
	log.Info("Process auth request")

	// Create test message
	msg := proto_resources.Msg{
		Uuid: &proto_resources.UUID{
			Value: uuid.NewV4().Bytes(),
		},
		Ts: &proto_resources.Timestamp{
			Value: uint64(time.Now().Unix()),
		},
		Payload: &proto_resources.Msg_AuthRequest{
			AuthRequest: &proto_resources.AuthRequest{
				Provider:                proto_resources.AuthProviderEnum_google,
				AuthProviderUrlResponse: c.Request.URL.String(),
			},
		},
	}

	// Serialize message
	msgBytes, err := proto.Marshal(&msg)
	if err != nil {
		c.String(500, "Couldn't serialize msg. %v", err)
		c.Abort()
		return
	}

	// Send request
	responseProto, err := nc.Request(config.NATS.Request.Subject, msgBytes, 1*time.Second)
	if err != nil {
		c.String(500, "Couldn't send request. %v", err)
		c.Abort()
		return
	}

	// Parse response
	response := proto_resources.Msg{}
	err = proto.Unmarshal(responseProto.Data, &response)
	if err != nil {
		c.String(500, "Couldn't deserialize response. %v", err)
		c.Abort()
		return
	}

	// Get user profile from response
	userProfile := response.GetAuthResponse()

	c.HTML(http.StatusOK, "auth.html", map[string]interface{}{
		"FirstName": userProfile.FirstName,
		"LastName":  userProfile.LastName,
		"Email":     userProfile.Email,
		"Locale":    userProfile.Locale,
	})
}
