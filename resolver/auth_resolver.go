package resolver

type LoginResponse struct {
	ok    bool
	token string
}

func (l LoginResponse) OK() bool {
	return l.ok
}

func (l LoginResponse) TOKEN() string {
	return l.token
}
