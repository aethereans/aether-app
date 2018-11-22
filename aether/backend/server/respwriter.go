package server

import (
	"net/http"
)

type CustomRespWriter struct {
	http.ResponseWriter
}

// We rewrite 404 Not founds to 204 No content. This is because 404 Not found causes the underlying TCP connection to be terminated, and in the case of reverse opens, that is the only lifeline we have into the NAT of the remote. If we terminate that connection, we do not get another.
func (r *CustomRespWriter) WriteHeader(statusCode int) {
	if statusCode == 404 {
		r.WriteHeader(204)
	} else {
		r.ResponseWriter.WriteHeader(statusCode)
	}
}
