package http

type MetricHttpError struct {
	ErrorText string `json:"error"`
	Code      int    `json:"code"`
}

func (e *MetricHttpError) Error() string {
	return e.ErrorText
}
