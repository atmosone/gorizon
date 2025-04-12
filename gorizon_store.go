package gorizon

import (
	"sync"
)

type store struct {
	mu       sync.RWMutex
	sessions map[*Session]struct{}
}

func (s *store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions = make(map[*Session]struct{})
}

func (s *store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

func (s *store) Create(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session] = struct{}{}
}

func (s *store) Delete(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, session)
}

func (s *store) ForEach(callback func(session *Session)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for session := range s.sessions {
		callback(session)
	}
}

func NewStore() *store {
	return &store{
		sessions: make(map[*Session]struct{}),
	}
}
