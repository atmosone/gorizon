package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atmosone/gorizon"
)

func main() {

	g := gorizon.NewRouter(gorizon.NewHub(gorizon.NewStore()))

	g.HandlerFunc("hello", func(s *gorizon.Session, m *gorizon.Message) {
		var data string
		if err := json.NewDecoder(bytes.NewBuffer(m.Payload)).Decode(&data); err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
		s.Write(m)
	})

	g.HandlerFunc("broadcast", func(s *gorizon.Session, m *gorizon.Message) {
		g.Broadcast(m)
	})

	http.ListenAndServe(":8080", g)
}
