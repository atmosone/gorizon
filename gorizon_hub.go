package gorizon

type Store interface {
	Reset()
	Count() int
	Create(*Session)
	Delete(*Session)
	ForEach(func(*Session))
}

type hub struct {
	store      Store
	connect    chan *Session
	disconnect chan *Session
	broadcast  chan *Message
	close      chan *Message
}

func (h *hub) Connect(session *Session) {
	h.connect <- session
}

func (h *hub) Disconnect(session *Session) {
	h.disconnect <- session
}

func (h *hub) Broadcast(message *Message) {
	h.broadcast <- message
}

func (h *hub) Open() {
loop:
	for {
		select {
		case s := <-h.connect:
			h.store.Create(s)
		case s := <-h.disconnect:
			s.Close()
			h.store.Delete(s)
		case m := <-h.broadcast:
			h.store.ForEach(func(s *Session) {
				s.Write(m)
			})
		case m := <-h.close:
			h.store.ForEach(func(s *Session) {
				s.Write(m)
				s.Close()
			})
			h.store.Reset()
			break loop
		}
	}
}

func (h *hub) Close() {
	h.close <- &Message{
		Topic:   "close",
		Payload: []byte{},
	}
}

func NewHub(store Store) *hub {
	return &hub{
		store:      store,
		connect:    make(chan *Session),
		disconnect: make(chan *Session),
		broadcast:  make(chan *Message),
		close:      make(chan *Message),
	}
}
