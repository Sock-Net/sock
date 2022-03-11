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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var READ_LIMIT int
var PING_SECOND int
var PORT string
var CONNECTION_TOKEN string

func main() {
	app := fiber.New()

	app.Use("/channel", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/channel/:channel", websocket.New(WebSocket))
	app.Get("/api/channel/:channel/:id", GetSock)
	app.Get("/api/channel/:channel", GetSocks)

	flag.IntVar(&READ_LIMIT, "limit", 1024, "Set sock's read limit")
	flag.IntVar(&PING_SECOND, "ping", 60, "Set sock's ping second")
	flag.StringVar(&PORT, "port", "", "Set sock's sevre port")
	flag.StringVar(&CONNECTION_TOKEN, "token", "demo", "Set sock's connection token")
	flag.Parse()

	// If port is not given
	if PORT == "" {
		// Get from env
		PORT = os.Getenv("PORT")

		// If still port is blank
		if PORT == "" {
			// Set port to 3000 by default
			PORT = "3000"
		}
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", PORT)))
}
