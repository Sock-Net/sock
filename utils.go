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
	"math/rand"
	"strings"
	"time"
)

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Create random id for connections
func RandomId() string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, 12)

	for i := range bytes {
		bytes[i] = LETTERS[rand.Intn(len(LETTERS))]
	}

	return string(bytes)
}

// Check if string is right channel format
func IsChannelFormat(str string) bool {
	if len(str) < 3 && len(str) > 16 {
		return false
	}

	for _, char := range str {
		if !strings.Contains(LETTERS, string(char)) {
			return false
		}
	}
	return true
}

// Write a message to sock instance
func (s *Sock) WriteMessage(messageType int, message []byte) error {
	return s.Connection.WriteMessage(messageType, message)
}

// Remove sock from connections list
func (s *Sock) Destroy() {
	for index, sock := range CONNECTIONS {
		if sock.Id == s.Id {
			CONNECTIONS[index] = CONNECTIONS[len(CONNECTIONS)-1]
			CONNECTIONS = CONNECTIONS[:len(CONNECTIONS)-1]

			return
		}
	}
}

// Start checker (if sock is alive)
func (s *Sock) StartPingChecker() {
	go func() {
		for {
			time.Sleep(time.Duration(PING_SECOND) * time.Second)

			if s.Pinged {
				s.Pinged = false
			} else {
				s.Deleted = true
				return
			}
		}
	}()
}

// Find all connections by channel
func FindConnections(channel string) []*Sock {
	var socks []*Sock

	for _, sock := range CONNECTIONS {
		if sock.Channel == channel {
			socks = append(socks, sock)
		}
	}

	return socks
}

// Check if id exists in channel
func CheckIdExists(channel, id string) bool {
	if len(id) < 3 && len(id) > 16 {
		return false
	}

	for _, sock := range CONNECTIONS {
		if sock.Channel == channel && sock.Id == id {
			return true
		}
	}

	return false
}
