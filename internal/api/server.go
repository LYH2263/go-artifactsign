package api

import (
	"net/http"

	"example.com/artifactsign"
)

// Options HTTP 选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server HTTP 门面。
type Server struct {
	svc  *artifactsign.Service
	opts Options
	mux  *http.ServeMux
}

// New 构造。
func New(svc *artifactsign.Service, opts Options) *Server {
	s := &Server{svc: svc, opts: opts, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}
