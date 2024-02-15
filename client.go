// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Replace "http://example.com" with the origin you want to allow.
		//allowedOrigin := "https://0dfd-186-206-50-111.ngrok-free.app"
		return true
	},
}

type Message struct {
	Username string     `json:"username"`
	Message  [][]string `json:"message"`
}

type Player struct {
	Id       string     `json:"id"`
	Username string     `json:"username"`
	Shape    [][]string `json:"shape"`
	Bloco    Bloco      `json:"bloco"`
	Num      int        `json:"num"`
}

type Bloco struct {
	Position Position `json:"position"`
	Shape    Shape    `json:"shape"`
}

type Shape struct {
	Number int        `json:"number"`
	Format [][]string `json:"format"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	player *Player
	// Buffered channel of outbound messages.
	send chan map[string]*Player
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Player

		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var player = c.hub.players[msg.Id]
		player.Shape = msg.Shape
		player.Bloco = msg.Bloco
		c.hub.broadcast <- c.hub.players
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	id := r.URL.Query().Get("idPlayer")
	var player = &Player{
		Id:       id,
		Username: id,
		Shape:    [][]string{},
		Bloco:    Bloco{},
		Num:      len(hub.players),
	}
	client := &Client{player: player, hub: hub, conn: conn, send: make(chan map[string]*Player, 1024)}
	client.hub.register <- client
	hub.players[id] = player

	go client.writePump()
	go client.readPump()
}
