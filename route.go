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

import "github.com/gofiber/fiber/v2"

// Error struct
type Error struct {
	Message string `json:"error"`
}

func GetSock(c *fiber.Ctx) error {
	channel := c.Params("channel")
	id := c.Params("id")

	for _, sock := range CONNECTIONS {
		if sock.Channel == channel && sock.Id == id {
			return c.Status(200).JSON(sock)
		}
	}

	return c.Status(404).JSON(Error{"Instance not found"})
}

func GetSocks(c *fiber.Ctx) error {
	channel := c.Params("channel")
	socks := FindConnections(channel)

	if socks != nil {
		return c.Status(200).JSON(socks)
	}

	return c.Status(404).JSON(Error{"Channel not found"})
}
