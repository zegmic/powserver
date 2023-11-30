package challange

type Empty struct{}

func (e *Empty) Create() string {
	return "test"
}

func (e *Empty) Verify(challenge string, solution int) bool {
	return true
}
