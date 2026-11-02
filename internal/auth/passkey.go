package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/andersmmg/bytewing/internal/db/sqlcgen"
)

// Passkeys (WebAuthn) are offered alongside the password, rather than as
// a replacement.
//
// Registration and revocation both re-ask for the password.

// bounds the gap between the two requests
const DefaultPasskeyCeremonyTTL = 5 * time.Minute

const MaxPasskeyNameLength = 60

// Errors returned by the passkey ceremonies.
var (
	ErrPasskeysUnavailable    = errors.New("auth: passkeys are not configured")
	ErrInvalidPasskey         = errors.New("auth: that passkey could not be verified")
	ErrPasskeyCeremonyExpired = errors.New("auth: that passkey prompt has expired, try again")
	ErrPasskeyNotFound        = errors.New("auth: no such passkey")
)

type Passkeys struct {
	wa *webauthn.WebAuthn
}

func NewPasskeys(rpID, origin string) (*Passkeys, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "ByteWing",
		RPID:          rpID,
		RPOrigins:     []string{origin},
	})
	if err != nil {
		return nil, fmt.Errorf("auth: webauthn configuration: %w", err)
	}
	return &Passkeys{wa: wa}, nil
}

// attaches WebAuthn support
func WithPasskeys(p *Passkeys) Option {
	return func(s *Service) { s.passkeys = p }
}

type Passkey struct {
	ID         int64
	Name       string
	Transports []string
	CreatedAt  time.Time
	LastUsedAt time.Time
}

// registration challenge vs login challenge
type ceremonyKind uint8

const (
	ceremonyRegistration ceremonyKind = iota + 1
	ceremonyLogin
)

type ceremony struct {
	kind      ceremonyKind
	userID    int64
	data      webauthn.SessionData
	expiresAt time.Time
}

// keeps in-flight challenges in memory.
type ceremonyStore struct {
	mu      sync.Mutex
	entries map[string]ceremony
}

func newCeremonyStore() *ceremonyStore {
	return &ceremonyStore{entries: make(map[string]ceremony)}
}

// records a fresh challenge, sweeping expired ones as it goes.
func (c *ceremonyStore) put(id string, entry ceremony, now time.Time, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry.expiresAt = now.Add(ttl)
	for key, existing := range c.entries {
		if now.After(existing.expiresAt) {
			delete(c.entries, key)
		}
	}
	c.entries[id] = entry
}

// removes a challenge and reports it
func (c *ceremonyStore) take(id string, now time.Time) (ceremony, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[id]
	if !ok {
		return ceremony{}, false
	}
	delete(c.entries, id)
	if now.After(entry.expiresAt) {
		return ceremony{}, false
	}
	return entry, true
}

func newCeremonyID() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("auth: ceremony id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

type storedPasskey struct {
	cred  webauthn.Credential
	rowID int64
}

type passkeyUser struct {
	user  User
	creds []storedPasskey
}

func (u passkeyUser) WebAuthnID() []byte {
	return []byte(strconv.FormatInt(u.user.ID, 10))
}

func (u passkeyUser) WebAuthnName() string { return u.user.Username }

func (u passkeyUser) WebAuthnDisplayName() string { return u.user.Username }

func (u passkeyUser) WebAuthnCredentials() []webauthn.Credential {
	out := make([]webauthn.Credential, 0, len(u.creds))
	for _, stored := range u.creds {
		out = append(out, stored.cred)
	}
	return out
}

// returns the instance's WebAuthn configuration, or an error if it is not
// available.
func (s *Service) engine() (*Passkeys, error) {
	if s.passkeys == nil {
		return nil, ErrPasskeysUnavailable
	}
	return s.passkeys, nil
}

// resolves an account together with every passkey it holds
func (s *Service) loadPasskeys(ctx context.Context, userID int64) (passkeyUser, []storedPasskey, error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return passkeyUser{}, nil, err
	}
	if user.Disabled {
		return passkeyUser{}, nil, ErrInvalidCredentials
	}

	rows, err := s.reads.ListWebAuthnCredentialsForUser(ctx, userID)
	if err != nil {
		return passkeyUser{}, nil, fmt.Errorf("auth: list passkeys: %w", err)
	}

	stored := make([]storedPasskey, 0, len(rows))
	for _, row := range rows {
		var cred webauthn.Credential
		if err := json.Unmarshal([]byte(row.Credential), &cred); err != nil {
			continue
		}
		stored = append(stored, storedPasskey{cred: cred, rowID: row.ID})
	}
	return passkeyUser{user: user, creds: stored}, stored, nil
}

// ------------------------------------------------------------- registration

func (s *Service) BeginPasskeyRegistration(ctx context.Context, userID int64) (string, *protocol.CredentialCreation, error) {
	engine, err := s.engine()
	if err != nil {
		return "", nil, err
	}

	user, stored, err := s.loadPasskeys(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	// for catching duplicates
	exclusions := make([]webauthn.Credential, 0, len(stored))
	for _, held := range stored {
		exclusions = append(exclusions, held.cred)
	}

	creation, session, err := engine.wa.BeginRegistration(
		user,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationRequired,
		}),
		webauthn.WithExclusions(webauthn.Credentials(exclusions).CredentialDescriptors()),
	)
	if err != nil {
		return "", nil, fmt.Errorf("auth: begin passkey registration: %w", err)
	}

	id, err := newCeremonyID()
	if err != nil {
		return "", nil, err
	}
	s.ceremonies.put(id, ceremony{
		kind:   ceremonyRegistration,
		userID: userID,
		data:   *session,
	}, s.now(), DefaultPasskeyCeremonyTTL)

	return id, creation, nil
}

func (s *Service) FinishPasskeyRegistration(ctx context.Context, userID int64, ceremonyID, name string, response []byte) (Passkey, error) {
	engine, err := s.engine()
	if err != nil {
		return Passkey{}, err
	}

	entry, ok := s.ceremonies.take(ceremonyID, s.now())
	if !ok || entry.kind != ceremonyRegistration || entry.userID != userID {
		return Passkey{}, ErrPasskeyCeremonyExpired
	}

	user, _, err := s.loadPasskeys(ctx, userID)
	if err != nil {
		return Passkey{}, err
	}

	parsed, err := protocol.ParseCredentialCreationResponseBytes(response)
	if err != nil {
		return Passkey{}, ErrInvalidPasskey
	}
	cred, err := engine.wa.CreateCredential(user, entry.data, parsed)
	if err != nil {
		return Passkey{}, ErrInvalidPasskey
	}

	encoded, err := json.Marshal(cred)
	if err != nil {
		return Passkey{}, fmt.Errorf("auth: encode passkey: %w", err)
	}
	row, err := s.writes.CreateWebAuthnCredential(ctx, sqlcgen.CreateWebAuthnCredentialParams{
		UserID:       userID,
		CredentialID: cred.ID,
		Credential:   string(encoded),
		Name:         normalisePasskeyName(name),
		CreatedAt:    s.now().Unix(),
	})
	if err != nil {
		return Passkey{}, fmt.Errorf("auth: store passkey: %w", err)
	}
	return toPasskey(row), nil
}

// -------------------------------------------------------------------- login

func (s *Service) BeginPasskeyLogin(ctx context.Context, address string) (string, *protocol.CredentialAssertion, error) {
	engine, err := s.engine()
	if err != nil {
		return "", nil, err
	}
	if address != "" && !s.byAddress.Allow(address) {
		return "", nil, ErrRateLimited
	}

	assertion, session, err := engine.wa.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationRequired),
	)
	if err != nil {
		return "", nil, fmt.Errorf("auth: begin passkey login: %w", err)
	}

	id, err := newCeremonyID()
	if err != nil {
		return "", nil, err
	}
	s.ceremonies.put(id, ceremony{kind: ceremonyLogin, data: *session}, s.now(), DefaultPasskeyCeremonyTTL)

	return id, assertion, nil
}

func (s *Service) FinishPasskeyLogin(ctx context.Context, ceremonyID string, response []byte, address string) (User, error) {
	engine, err := s.engine()
	if err != nil {
		return User{}, err
	}
	if address != "" && !s.byAddress.Allow(address) {
		return User{}, ErrRateLimited
	}

	entry, ok := s.ceremonies.take(ceremonyID, s.now())
	if !ok || entry.kind != ceremonyLogin {
		return User{}, ErrPasskeyCeremonyExpired
	}

	fail := func(username string) (User, error) {
		if username != "" {
			s.byAccount.Fail(username)
		}
		if address != "" {
			s.byAddress.Fail(address)
		}
		return User{}, ErrInvalidPasskey
	}

	parsed, err := protocol.ParseCredentialRequestResponseBytes(response)
	if err != nil {
		return fail("")
	}

	var resolved *User
	var usedRowID int64
	handler := func(rawID []byte, userHandle []byte) (webauthn.User, error) {
		id, err := strconv.ParseInt(string(userHandle), 10, 64)
		if err != nil {
			return nil, ErrUserNotFound
		}
		user, stored, err := s.loadPasskeys(ctx, id)
		if err != nil {
			return nil, err
		}
		resolved = &user.user
		for _, held := range stored {
			if bytes.Equal(held.cred.ID, rawID) {
				usedRowID = held.rowID
				break
			}
		}
		return user, nil
	}

	user, cred, err := engine.wa.ValidatePasskeyLogin(handler, entry.data, parsed)
	if err != nil || resolved == nil || user == nil {
		username := ""
		if resolved != nil {
			username = resolved.Username
		}
		return fail(username)
	}

	s.recordPasskeyUse(ctx, usedRowID, cred)

	s.byAccount.Reset(resolved.Username)
	if address != "" {
		s.byAddress.Reset(address)
	}
	return *resolved, nil
}

func (s *Service) recordPasskeyUse(ctx context.Context, rowID int64, cred *webauthn.Credential) {
	if rowID == 0 {
		return
	}
	encoded, err := json.Marshal(cred)
	if err != nil {
		return
	}
	_ = s.writes.UpdateWebAuthnCredentialUse(ctx, sqlcgen.UpdateWebAuthnCredentialUseParams{
		Credential: string(encoded),
		LastUsedAt: sql.NullInt64{Int64: s.now().Unix(), Valid: true},
		ID:         rowID,
	})
}

// --------------------------------------------------------------- management

func (s *Service) ListPasskeys(ctx context.Context, userID int64) ([]Passkey, error) {
	rows, err := s.reads.ListWebAuthnCredentialsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("auth: list passkeys: %w", err)
	}
	passkeys := make([]Passkey, 0, len(rows))
	for _, row := range rows {
		passkeys = append(passkeys, toPasskey(row))
	}
	return passkeys, nil
}

func (s *Service) RemovePasskey(ctx context.Context, userID, passkeyID int64) error {
	rows, err := s.writes.DeleteWebAuthnCredentialForUser(ctx, sqlcgen.DeleteWebAuthnCredentialForUserParams{
		ID:     passkeyID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("auth: remove passkey: %w", err)
	}
	if rows == 0 {
		return ErrPasskeyNotFound
	}
	return nil
}

func (s *Service) CountPasskeys(ctx context.Context) (int64, error) {
	count, err := s.reads.CountWebAuthnCredentials(ctx)
	if err != nil {
		return 0, fmt.Errorf("auth: count passkeys: %w", err)
	}
	return count, nil
}

// ------------------------------------------------------------------ helpers

func normalisePasskeyName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Passkey"
	}
	if runes := []rune(name); len(runes) > MaxPasskeyNameLength {
		name = string(runes[:MaxPasskeyNameLength])
	}
	return name
}

func toPasskey(row sqlcgen.WebauthnCredential) Passkey {
	passkey := Passkey{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: time.Unix(row.CreatedAt, 0),
	}
	if row.LastUsedAt.Valid {
		passkey.LastUsedAt = time.Unix(row.LastUsedAt.Int64, 0)
	}
	var cred webauthn.Credential
	if err := json.Unmarshal([]byte(row.Credential), &cred); err == nil {
		for _, transport := range cred.Transport {
			passkey.Transports = append(passkey.Transports, string(transport))
		}
	}
	return passkey
}
