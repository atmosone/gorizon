package gorizon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub interface {
	Open()
	Connect(*Session)
	Disconnect(*Session)
	Broadcast(*Message)
	Close()
}

type HandlerFunc func(*Session, *Message)

type Gorizon struct {
	hub      Hub
	upgrader websocket.Upgrader
	routes   map[string]HandlerFunc
}

func (g *Gorizon) Open() {
	go g.hub.Open()
}

func (g *Gorizon) Close() {
	g.hub.Close()
}

func (g *Gorizon) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := g.upgrader.Upgrade(res, req, nil)
	if err != nil {
		return
	}

	session := NewSession(context.Background(), conn)
	g.hub.Connect(session)
	go g.route(session)
}

func (g *Gorizon) route(session *Session) {
loop:
	for {
		select {
		case <-session.Context().Done():
			g.hub.Disconnect(session)
			break loop
		case message := <-session.Read():
			fmt.Println(message)
			if h, ok := g.routes[message.Topic]; ok {
				h(session, message)
			}
			continue loop
		}
	}
}

func (g *Gorizon) OnMessage(topic string, h HandlerFunc) {
	g.routes[topic] = h
}

func (g *Gorizon) Broadcast(message *Message) {
	g.hub.Broadcast(message)
}

func New() *Gorizon {
	store := NewStore()
	hub := NewHub(store)
	return &Gorizon{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		routes: make(map[string]HandlerFunc),
	}
}
