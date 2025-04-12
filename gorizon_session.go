package gorizon

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
	input  chan *Message
	output chan *Message
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) Close() error {
	s.cancel()
	err := s.conn.Close()
	time.Sleep(100 * time.Millisecond)

	close(s.input)
	close(s.output)
	return err
}

func (s *Session) Write(message *Message) {
	s.output <- message
}

func (s *Session) Read() <-chan *Message {
	return s.input
}

func (s *Session) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in write:", r)
		}
		s.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case m := <-s.output:
			if err := s.conn.WriteJSON(m); err != nil {
				fmt.Printf("Write error: %v\n", err)
				return
			}
		}
	}
}

func (s *Session) read() {
	defer s.Close()

	for {
		message := &Message{}
		if err := s.conn.ReadJSON(message); err != nil {
			fmt.Printf("Read error: %v\n", err)
			return
		}

		select {
		case <-s.ctx.Done():
			return
		case s.input <- message:
		}
	}
}

func NewSession(ctx context.Context, conn *websocket.Conn) *Session {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
		input:  make(chan *Message, 100),
		output: make(chan *Message, 100),
	}

	go s.read()
	go s.write()

	return s
}
