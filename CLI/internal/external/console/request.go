package console

import "net/http"

type RequestSpec struct {
	Method   string
	Path     string
	Headers  http.Header
	Body     []byte
	ReadOnly bool
}
