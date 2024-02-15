// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var hubs = make(map[string]*Hub)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		valor, existe := hubs[id]
		if existe {
			hub := valor
			serveWs(hub, w, r)
		} else {
			hub := newHub()
			hub.players = make(map[string]*Player)
			go hub.run()
			hubs[id] = hub
			serveWs(hub, w, r)
		}
	})
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
