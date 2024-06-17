package states

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
