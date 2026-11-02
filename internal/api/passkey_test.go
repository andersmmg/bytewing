package api_test

import (
	"net/http"
	"testing"
)

// The full ceremony as tests

type passkeyStatus struct {
	Available bool `json:"available"`
}

type passkeyOptions struct {
	CeremonyID string         `json:"ceremony_id"`
	PublicKey  map[string]any `json:"public_key"`
}

func problemCode(t *testing.T, raw []byte) string {
	t.Helper()
	return decode[struct {
		Code string `json:"code"`
	}](t, raw).Code
}

func TestPasskeyStatusIsOpenAndSaysWhenNoneExist(t *testing.T) {
	t.Parallel()

	c := newClient(t)

	resp, raw := c.do(http.MethodGet, "/api/v1/auth/passkeys/status", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d: %s", resp.StatusCode, raw)
	}
	if decode[passkeyStatus](t, raw).Available {
		t.Error("an empty instance offered passkeys")
	}

	c.setUp()
	resp, raw = c.do(http.MethodGet, "/api/v1/auth/passkeys/status", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d: %s", resp.StatusCode, raw)
	}
	if decode[passkeyStatus](t, raw).Available {
		t.Error("an account with no passkey reported one available")
	}
}

func TestPasskeyRegisterOptionsDemandSessionAndPassword(t *testing.T) {
	t.Parallel()

	c := newClient(t)

	resp, raw := c.do(http.MethodPost, "/api/v1/auth/passkeys/register/options", map[string]string{"password": "x"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("signed out = %d, want 401: %s", resp.StatusCode, raw)
	}

	c.setUp()

	resp, raw = c.do(http.MethodPost, "/api/v1/auth/passkeys/register/options", map[string]string{"password": "not it"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password = %d, want 401: %s", resp.StatusCode, raw)
	}
	if code := problemCode(t, raw); code != "invalid_credentials" {
		t.Errorf("code = %q, want invalid_credentials", code)
	}

	resp, raw = c.do(http.MethodPost, "/api/v1/auth/passkeys/register/options",
		map[string]string{"password": "correct horse battery"})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("right password = %d: %s", resp.StatusCode, raw)
	}
	options := decode[passkeyOptions](t, raw)
	if options.CeremonyID == "" {
		t.Error("no ceremony id to answer with")
	}
	if options.PublicKey == nil {
		t.Fatal("no creation options for the browser")
	}
	rp, _ := options.PublicKey["rp"].(map[string]any)
	if id, _ := rp["id"].(string); id != "localhost" {
		t.Errorf("rp.id = %q, want the instance domain", id)
	}
	challenge, _ := options.PublicKey["challenge"].(string)
	if challenge == "" {
		t.Error("the options carry no challenge")
	}
}

func TestPasskeyListNeedsASession(t *testing.T) {
	t.Parallel()

	c := newClient(t)
	resp, raw := c.do(http.MethodGet, "/api/v1/auth/passkeys", nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("list = %d, want 401: %s", resp.StatusCode, raw)
	}

	c.setUp()
	resp, raw = c.do(http.MethodGet, "/api/v1/auth/passkeys", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list = %d: %s", resp.StatusCode, raw)
	}
}

func TestPasskeyVerifyRefusesUnknownChallenges(t *testing.T) {
	t.Parallel()

	open := newClient(t)
	resp, raw := open.do(http.MethodPost, "/api/v1/auth/passkeys/login/verify",
		map[string]any{"ceremony_id": "never-issued", "credential": map[string]any{}})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("login verify = %d, want 400: %s", resp.StatusCode, raw)
	}
	if code := problemCode(t, raw); code != "passkey_expired" {
		t.Errorf("login code = %q, want passkey_expired", code)
	}

	c := newClient(t)
	c.setUp()
	resp, raw = c.do(http.MethodPost, "/api/v1/auth/passkeys/register/verify",
		map[string]any{"ceremony_id": "never-issued", "name": "x", "credential": map[string]any{}})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("register verify = %d, want 400: %s", resp.StatusCode, raw)
	}
	if code := problemCode(t, raw); code != "passkey_expired" {
		t.Errorf("register code = %q, want passkey_expired", code)
	}
}

func TestPasskeyLoginOptionsAreOpen(t *testing.T) {
	t.Parallel()

	c := newClient(t)
	resp, raw := c.do(http.MethodPost, "/api/v1/auth/passkeys/login/options", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login options = %d: %s", resp.StatusCode, raw)
	}
	options := decode[passkeyOptions](t, raw)
	if options.CeremonyID == "" {
		t.Error("no ceremony id")
	}
	if challenge, _ := options.PublicKey["challenge"].(string); challenge == "" {
		t.Error("no challenge to answer")
	}
	if rpID, _ := options.PublicKey["rpId"].(string); rpID != "localhost" {
		t.Errorf("rpId = %q, want the instance domain", rpID)
	}
}
