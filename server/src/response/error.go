package response

type Error struct {
	Code  int
	Cause error
}

func (e Error) Error() string {
	return e.Cause.Error()
}
