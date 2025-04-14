package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atmosone/gorizon"
)

func main() {
	g := gorizon.New()
	g.Open()
	defer g.Close()

	g.OnMessage("hello", func(s *gorizon.Session, m *gorizon.Message) {
		var data string
		if err := json.NewDecoder(bytes.NewBuffer(m.Payload)).Decode(&data); err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
		s.Write(m)
	})

	g.OnMessage("broadcast", func(s *gorizon.Session, m *gorizon.Message) {
		g.Broadcast(m)
	})

	mux := http.NewServeMux()

	mux.Handle("/ws", g)

	http.ListenAndServe(":8080", mux)
}
