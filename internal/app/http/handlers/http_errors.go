package handlers

type ErrorHTTP struct {
	ErrorText string `json:"error"`
	Code      int    `json:"code"`
}

func (e *ErrorHTTP) Error() string {
	return e.ErrorText
}
