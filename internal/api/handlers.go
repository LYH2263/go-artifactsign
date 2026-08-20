package api

import (
	"encoding/json"
	"io"
	"net/http"

	"example.com/artifactsign"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{"ok": true, "closed": s.svc.Closed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.svc.Stats())
}

func (s *Server) handleKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID       string `json:"id"`
		Material []byte `json:"material"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.svc.RegisterKey(req.ID, req.Material); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]string{"id": req.ID})
}

func (s *Server) handleSign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Payload []byte `json:"payload"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	view, err := s.svc.Sign(req.Payload, artifactsign.Meta{Name: req.Name})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, view)
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		KeyID   string `json:"key_id"`
		Payload []byte `json:"payload"`
		Sig     []byte `json:"sig"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err := s.svc.Verify(req.KeyID, req.Payload, req.Sig)
	writeJSON(w, map[string]any{"ok": err == nil, "error": errString(err)})
}

func (s *Server) handleTrust(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, s.svc.ListTrust())
	case http.MethodPost:
		var req struct {
			ID       string `json:"id"`
			Material []byte `json:"material"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := s.svc.RegisterTrust(req.ID, req.Material); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"id": req.ID})
	default:
		http.Error(w, "method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, s.svc.ListRevoked())
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Fingerprint string `json:"fingerprint"`
		Reason      string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.svc.RevokeFingerprint(req.Fingerprint, req.Reason); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]string{"ok": "1"})
}

func (s *Server) handleDigest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	hex, err := s.svc.Digest(b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]string{"digest": hex})
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
