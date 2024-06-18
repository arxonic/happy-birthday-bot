package states

import "sync"

// Conditions
const (
	StateAuthMiddleware = iota
	StateEmailWait
	StateEmailSent

	StateMenu

	StateFind
	StateSubscribe
)

type States struct {
	mx         sync.RWMutex
	UserStates map[int64]*UserState
}

type UserState struct {
	State  int
	Finder FindState
}

type FindState struct {
	Organization string
	City         string
	Office       string
	Department   string
	UserID       int64
}

func NewStates() *States {
	return &States{
		UserStates: make(map[int64]*UserState),
	}
}

func (s *States) Load(key int64) (*UserState, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.UserStates[key]

	return val, ok
}

func (s *States) Store(key int64, value *UserState) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.UserStates[key] = value
}
