package api

import "net/http"

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/keys", s.handleKeys)
	s.mux.HandleFunc("/api/sign", s.handleSign)
	s.mux.HandleFunc("/api/verify", s.handleVerify)
	s.mux.HandleFunc("/api/trust", s.handleTrust)
	s.mux.HandleFunc("/api/revoke", s.handleRevoke)
	s.mux.HandleFunc("/api/digest", s.handleDigest)
	if s.opts.WebDir != "" {
		fs := http.FileServer(http.Dir(s.opts.WebDir))
		s.mux.Handle("/", fs)
	}
}
