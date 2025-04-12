package gorizon

import (
	"context"
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

type router struct {
	hub      Hub
	upgrader websocket.Upgrader
	routes   map[string]HandlerFunc
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := r.upgrader.Upgrade(res, req, nil)
	if err != nil {
		return
	}

	session := NewSession(context.Background(), conn)
	r.hub.Connect(session)
	go r.route(session)
}

func (r *router) route(session *Session) {
loop:
	for {
		select {
		case <-session.Context().Done():
			r.hub.Disconnect(session)
			break loop
		case message, ok := <-session.Read():
			if !ok {
				r.hub.Disconnect(session)
				break loop
			}
			if h, ok := r.routes[message.Topic]; ok {
				h(session, message)
			}
			continue loop
		}
	}
}

func (r *router) HandlerFunc(topic string, h HandlerFunc) {
	r.routes[topic] = h
}

func (r *router) Broadcast(message *Message) {
	r.hub.Broadcast(message)
}

func NewRouter(hub Hub) *router {
	return &router{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		routes: make(map[string]HandlerFunc),
	}
}
