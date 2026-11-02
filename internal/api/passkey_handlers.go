package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-webauthn/webauthn/protocol"

	"github.com/andersmmg/bytewing/internal/audit"
	"github.com/andersmmg/bytewing/internal/auth"
)

// Passkeys (WebAuthn) are offered alongside the password, rather than as
// a replacement.
//
// Registration and revocation both re-ask for the password.

type passkeyResponse struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Transports []string   `json:"transports,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

type passkeyListResponse struct {
	Passkeys []passkeyResponse `json:"passkeys"`
}

type passkeyStatusResponse struct {
	Available bool `json:"available"`
}

type passwordRequest struct {
	Password string `json:"password"`
}

type passkeyOptionsResponse struct {
	CeremonyID string                                      `json:"ceremony_id"`
	PublicKey  protocol.PublicKeyCredentialCreationOptions `json:"public_key"`
}

type passkeyLoginOptionsResponse struct {
	CeremonyID string                                     `json:"ceremony_id"`
	PublicKey  protocol.PublicKeyCredentialRequestOptions `json:"public_key"`
}

type passkeyRegisterVerifyRequest struct {
	CeremonyID string          `json:"ceremony_id"`
	Name       string          `json:"name"`
	Credential json.RawMessage `json:"credential"`
}

type passkeyLoginVerifyRequest struct {
	CeremonyID string          `json:"ceremony_id"`
	Credential json.RawMessage `json:"credential"`
}

func (s *Server) handlePasskeyStatus(w http.ResponseWriter, r *http.Request) {
	count, err := s.auth.CountPasskeys(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, passkeyStatusResponse{Available: count > 0})
}

func (s *Server) handleListPasskeys(w http.ResponseWriter, r *http.Request) {
	c, ok := callerFrom(r.Context())
	if !ok {
		writeProblem(w, r, http.StatusUnauthorized, CodeUnauthenticated, "Not signed in", "")
		return
	}

	passkeys, err := s.auth.ListPasskeys(r.Context(), c.user.ID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	out := make([]passkeyResponse, 0, len(passkeys))
	for _, p := range passkeys {
		out = append(out, toPasskeyResponse(p))
	}
	writeJSON(w, r, http.StatusOK, passkeyListResponse{Passkeys: out})
}

// Start a registration ceremony
func (s *Server) handlePasskeyRegisterOptions(w http.ResponseWriter, r *http.Request) {
	c, ok := callerFrom(r.Context())
	if !ok {
		writeProblem(w, r, http.StatusUnauthorized, CodeUnauthenticated, "Not signed in", "")
		return
	}

	var body passwordRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if _, err := s.auth.Authenticate(r.Context(), c.user.Username, body.Password, clientIP(r)); err != nil {
		writeError(w, r, err)
		return
	}

	ceremonyID, creation, err := s.auth.BeginPasskeyRegistration(r.Context(), c.user.ID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, passkeyOptionsResponse{
		CeremonyID: ceremonyID,
		PublicKey:  creation.Response,
	})
}

func (s *Server) handlePasskeyRegisterVerify(w http.ResponseWriter, r *http.Request) {
	c, ok := callerFrom(r.Context())
	if !ok {
		writeProblem(w, r, http.StatusUnauthorized, CodeUnauthenticated, "Not signed in", "")
		return
	}

	var req passkeyRegisterVerifyRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	passkey, err := s.auth.FinishPasskeyRegistration(
		r.Context(), c.user.ID, req.CeremonyID, req.Name, req.Credential,
	)
	if err != nil {
		writeError(w, r, err)
		return
	}
	s.audit(r, audit.ActionPasskeyRegistered, "", map[string]any{"name": passkey.Name})

	writeJSON(w, r, http.StatusCreated, toPasskeyResponse(passkey))
}

func (s *Server) handlePasskeyRemove(w http.ResponseWriter, r *http.Request) {
	c, ok := callerFrom(r.Context())
	if !ok {
		writeProblem(w, r, http.StatusUnauthorized, CodeUnauthenticated, "Not signed in", "")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeProblem(w, r, http.StatusNotFound, CodeNotFound, "No such passkey", "")
		return
	}

	var body passwordRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if _, err := s.auth.Authenticate(r.Context(), c.user.Username, body.Password, clientIP(r)); err != nil {
		writeError(w, r, err)
		return
	}

	if err := s.auth.RemovePasskey(r.Context(), c.user.ID, id); err != nil {
		writeError(w, r, err)
		return
	}
	s.audit(r, audit.ActionPasskeyRevoked, strconv.FormatInt(id, 10), nil)

	writeJSON(w, r, http.StatusNoContent, nil)
}

func (s *Server) handlePasskeyLoginOptions(w http.ResponseWriter, r *http.Request) {
	ceremonyID, assertion, err := s.auth.BeginPasskeyLogin(r.Context(), clientIP(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, passkeyLoginOptionsResponse{
		CeremonyID: ceremonyID,
		PublicKey:  assertion.Response,
	})
}

func (s *Server) handlePasskeyLoginVerify(w http.ResponseWriter, r *http.Request) {
	var body passkeyLoginVerifyRequest
	if !decodeJSON(w, r, &body) {
		return
	}

	user, err := s.auth.FinishPasskeyLogin(r.Context(), body.CeremonyID, body.Credential, clientIP(r))
	if err != nil {
		writeError(w, r, err)
		return
	}

	token, session, err := s.auth.CreateSession(r.Context(), user.ID, r.UserAgent(), clientIP(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	if err := s.ensureHome(r.Context(), user); err != nil {
		writeError(w, r, err)
		return
	}

	http.SetCookie(w, auth.SessionCookie(token, session.ExpiresAt, s.secureCookies))
	s.auditAs(r, user.ID, audit.ActionLogin, "", map[string]any{"via": "passkey"})

	writeJSON(w, r, http.StatusOK, sessionResponse{
		User:      toUserResponse(user),
		ExpiresAt: session.ExpiresAt,
		Token:     token,
	})
}

func toPasskeyResponse(p auth.Passkey) passkeyResponse {
	out := passkeyResponse{
		ID:         p.ID,
		Name:       p.Name,
		Transports: p.Transports,
		CreatedAt:  p.CreatedAt,
	}
	if !p.LastUsedAt.IsZero() {
		last := p.LastUsedAt
		out.LastUsedAt = &last
	}
	return out
}
