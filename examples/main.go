package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atmosone/gorizon"
)

func main() {
	store := gorizon.NewStore()
	hub := gorizon.NewHub(store)
	go hub.Open()

	router := gorizon.NewRouter(hub)

	router.OnMessage("hello", func(s *gorizon.Session, m *gorizon.Message) {
		var data string
		if err := json.NewDecoder(bytes.NewBuffer(m.Payload)).Decode(&data); err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
		s.Write(m)
	})

	router.OnMessage("broadcast", func(s *gorizon.Session, m *gorizon.Message) {
		router.Broadcast(m)
	})

	mux := http.NewServeMux()

	mux.Handle("/ws", router)

	http.ListenAndServe(":8080", mux)
}
