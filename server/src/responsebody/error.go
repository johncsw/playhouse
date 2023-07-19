package responsebody

type ResponseErr struct {
	Code    int
	ErrBody error
}

func (h ResponseErr) Error() string {
	return h.ErrBody.Error()
}
