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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var READ_LIMIT int
var PORT int
var CONNECTION_TOKEN string

func main() {
	app := fiber.New()

	app.Use("/websocket", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/websocket/:channel", websocket.New(WebSocket))

	flag.IntVar(&READ_LIMIT, "limit", 1024, "Set sock's read limit")
	flag.IntVar(&PORT, "port", 3000, "Set sock's sevre port")
	flag.StringVar(&CONNECTION_TOKEN, "token", "demo", "Set sock's connection token")
	flag.Parse()

	log.Fatal(app.Listen(fmt.Sprintf(":%d", PORT)))
}
