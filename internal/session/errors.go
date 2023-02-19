package session

type CacheError struct {
	Message string
}

func (e CacheError) Error() string {
	return e.Message
}

type SessionNotFoundError struct {
	SessionID string
}

func (e SessionNotFoundError) Error() string {
	return e.SessionID
}
