package v1

import "net/http"

func NewServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: handler,
	}
}

func Run(srv *http.Server) error {
	return srv.ListenAndServe()
}
