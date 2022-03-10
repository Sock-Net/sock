// Copyright (C) 2022 aiocat
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
)

var CONNECTIONS []*Sock // List all avaible sock connections

// Sock struct
type Sock struct {
	Connection *websocket.Conn
	Pinged     bool
	Channel    string
	Id         string
	CreatedAt  time.Time
}

// WebsocketMessage struct
type WebSocketMessage struct {
	Type    int    `json:"type"`
	Message string `json:"data"`
	From    string `json:"from"`
}

// Websocket handler
func WebSocket(c *websocket.Conn) {
	// c.Locals is added to the *websocket.Conn
	if !c.Locals("allowed").(bool) {
		c.Close()
	}

	channel := c.Params("channel", "default")

	if !IsChannelFormat(channel) {
		c.Close()
	}

	// New sock instance
	sock := Sock{
		Connection: c,
		Pinged:     true,
		Channel:    channel,
		Id:         RandomId(),
		CreatedAt:  time.Now(),
	}

	// Add to connections
	CONNECTIONS = append(CONNECTIONS, &sock)

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		messageType int
		message     []byte
		wsError     error
	)
	for {
		// Read message
		if messageType, message, wsError = c.ReadMessage(); wsError != nil {
			log.Println("Read error:", wsError)
			break
		}

		// Decode message
		wsMessage := new(WebSocketMessage)
		err := json.Unmarshal(message, &wsMessage)

		if err != nil {
			// Write error message
			if wsError = sock.WriteMessage(messageType, []byte("{\"error\":\"Invalid json format\"}")); wsError != nil {
				log.Println("Write error:", wsError)
				break
			}
		}

		// Prepare websocket message
		wsMessage.From = sock.Id
		marshalledMessage, err := json.Marshal(wsMessage)

		if err != nil {
			// Write error message
			if wsError = sock.WriteMessage(messageType, []byte("{\"error\":\"Unknown error\"}")); wsError != nil {
				log.Println("Write error:", wsError)
				break
			}
		}

		// Send message to the all instances in the channel
		connectedInstances := FindConnections(sock.Channel)
		for _, instance := range connectedInstances {
			if wsError = instance.WriteMessage(messageType, marshalledMessage); wsError != nil {
				log.Println("Write error:", wsError)
				break
			}
		}
	}
}
