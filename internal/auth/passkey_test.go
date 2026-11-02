package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol/webauthncbor"
)

const (
	testRPID    = "localhost"
	testOrigin  = "http://localhost"
	testAddress = "10.0.0.7"
)

func testPasskeys(t *testing.T) *Passkeys {
	t.Helper()
	passkeys, err := NewPasskeys(testRPID, testOrigin)
	if err != nil {
		t.Fatalf("NewPasskeys: %v", err)
	}
	return passkeys
}

type fakeAuthenticator struct {
	key     *ecdsa.PrivateKey
	credID  []byte
	counter uint32
	rpID    string
	origin  string
}

func newFakeAuthenticator(t *testing.T) *fakeAuthenticator {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	credID := make([]byte, 32)
	if _, err := rand.Read(credID); err != nil {
		t.Fatalf("credential id: %v", err)
	}
	return &fakeAuthenticator{key: key, credID: credID, rpID: testRPID, origin: testOrigin}
}

func b64(data []byte) string { return base64.RawURLEncoding.EncodeToString(data) }

func (a *fakeAuthenticator) clientData(kind, challenge string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":        kind,
		"challenge":   challenge,
		"origin":      a.origin,
		"crossOrigin": false,
	})
}

func (a *fakeAuthenticator) authData(flags byte, counter uint32, rest []byte) []byte {
	rpIDHash := sha256.Sum256([]byte(a.rpID))
	data := make([]byte, 0, 37+len(rest))
	data = append(data, rpIDHash[:]...)
	data = append(data, flags)
	data = binary.BigEndian.AppendUint32(data, counter)
	return append(data, rest...)
}

func (a *fakeAuthenticator) register(challenge string) ([]byte, error) {
	clientData, err := a.clientData("webauthn.create", challenge)
	if err != nil {
		return nil, err
	}

	coseKey, err := webauthncbor.Marshal(map[int]any{
		1:  int64(2),  // kty: EC2
		3:  int64(-7), // alg: ES256
		-1: int64(1),  // crv: P-256
		-2: a.key.X.FillBytes(make([]byte, 32)),
		-3: a.key.Y.FillBytes(make([]byte, 32)),
	})
	if err != nil {
		return nil, err
	}

	attested := make([]byte, 0, 16+2+len(a.credID)+len(coseKey))
	attested = append(attested, make([]byte, 16)...) // AAGUID: an unknown model
	attested = binary.BigEndian.AppendUint16(attested, uint16(len(a.credID)))
	attested = append(attested, a.credID...)
	attested = append(attested, coseKey...)

	const up, uv, at = 0x01, 0x04, 0x40
	authData := a.authData(up|uv|at, a.counter, attested)

	attestation, err := webauthncbor.Marshal(map[string]any{
		"fmt":      "none",
		"attStmt":  map[string]any{},
		"authData": authData,
	})
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"id":    b64(a.credID),
		"rawId": b64(a.credID),
		"type":  "public-key",
		"response": map[string]any{
			"clientDataJSON":    b64(clientData),
			"attestationObject": b64(attestation),
			"transports":        []string{"internal"},
		},
	})
}

func (a *fakeAuthenticator) assert(challenge string, userHandle []byte) ([]byte, error) {
	clientData, err := a.clientData("webauthn.get", challenge)
	if err != nil {
		return nil, err
	}
	clientHash := sha256.Sum256(clientData)

	const up, uv = 0x01, 0x04
	a.counter++
	authData := a.authData(up|uv, a.counter, nil)

	signed := make([]byte, 0, len(authData)+len(clientHash))
	signed = append(signed, authData...)
	signed = append(signed, clientHash[:]...)
	digest := sha256.Sum256(signed)
	signature, err := ecdsa.SignASN1(rand.Reader, a.key, digest[:])
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"id":    b64(a.credID),
		"rawId": b64(a.credID),
		"type":  "public-key",
		"response": map[string]any{
			"clientDataJSON":    b64(clientData),
			"authenticatorData": b64(authData),
			"signature":         b64(signature),
			"userHandle":        b64(userHandle),
		},
	})
}

func registerPasskey(t *testing.T, f *fixture, user User, authenticator *fakeAuthenticator, name string) Passkey {
	t.Helper()
	ctx := t.Context()

	ceremonyID, creation, err := f.svc.BeginPasskeyRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginPasskeyRegistration: %v", err)
	}
	response, err := authenticator.register(creation.Response.Challenge.String())
	if err != nil {
		t.Fatalf("authenticator register: %v", err)
	}
	passkey, err := f.svc.FinishPasskeyRegistration(ctx, user.ID, ceremonyID, name, response)
	if err != nil {
		t.Fatalf("FinishPasskeyRegistration: %v", err)
	}
	return passkey
}

func userHandle(user User) []byte { return []byte(strconv.FormatInt(user.ID, 10)) }

func TestPasskeyRegisterThenLogin(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)

	passkey := registerPasskey(t, f, user, authenticator, "Laptop")
	if passkey.ID == 0 {
		t.Fatal("the stored passkey has no id")
	}
	if phone := registerPasskey(t, f, user, newFakeAuthenticator(t), "  Phone  "); phone.Name != "Phone" {
		t.Errorf("name = %q, want the label trimmed", phone.Name)
	}

	list, err := f.svc.ListPasskeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListPasskeys: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("ListPasskeys = %d, want 2", len(list))
	}
	count, err := f.svc.CountPasskeys(ctx)
	if err != nil {
		t.Fatalf("CountPasskeys: %v", err)
	}
	if count != 2 {
		t.Errorf("CountPasskeys = %d, want 2", count)
	}

	ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
	if err != nil {
		t.Fatalf("BeginPasskeyLogin: %v", err)
	}
	response, err := authenticator.assert(assertion.Response.Challenge.String(), userHandle(user))
	if err != nil {
		t.Fatalf("authenticator assert: %v", err)
	}
	signedIn, err := f.svc.FinishPasskeyLogin(ctx, ceremonyID, response, testAddress)
	if err != nil {
		t.Fatalf("FinishPasskeyLogin: %v", err)
	}
	if signedIn.ID != user.ID {
		t.Errorf("signed in as %q, want %q", signedIn.Username, user.Username)
	}

	list, err = f.svc.ListPasskeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListPasskeys: %v", err)
	}
	for _, p := range list {
		if p.ID != passkey.ID {
			continue
		}
		if p.LastUsedAt.IsZero() {
			t.Error("the passkey that signed in has no last-used time")
		}
	}
}

func TestPasskeyCeremonyIsSingleUse(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)

	ceremonyID, creation, err := f.svc.BeginPasskeyRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginPasskeyRegistration: %v", err)
	}
	response, err := authenticator.register(creation.Response.Challenge.String())
	if err != nil {
		t.Fatalf("authenticator register: %v", err)
	}
	if _, err := f.svc.FinishPasskeyRegistration(ctx, user.ID, ceremonyID, "First", response); err != nil {
		t.Fatalf("first FinishPasskeyRegistration: %v", err)
	}

	if _, err := f.svc.FinishPasskeyRegistration(ctx, user.ID, ceremonyID, "Second", response); !errors.Is(err, ErrPasskeyCeremonyExpired) {
		t.Errorf("replayed response = %v, want ErrPasskeyCeremonyExpired", err)
	}
	list, err := f.svc.ListPasskeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListPasskeys: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("ListPasskeys = %d, want the replay to add nothing", len(list))
	}
}

func TestPasskeyCeremonyExpires(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)

	ceremonyID, creation, err := f.svc.BeginPasskeyRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginPasskeyRegistration: %v", err)
	}
	response, err := authenticator.register(creation.Response.Challenge.String())
	if err != nil {
		t.Fatalf("authenticator register: %v", err)
	}

	f.advance(DefaultPasskeyCeremonyTTL + time.Second)

	if _, err := f.svc.FinishPasskeyRegistration(ctx, user.ID, ceremonyID, "Late", response); !errors.Is(err, ErrPasskeyCeremonyExpired) {
		t.Errorf("a challenge past its TTL = %v, want ErrPasskeyCeremonyExpired", err)
	}
}

func TestPasskeyCeremoniesAreNotInterchangeable(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()

	loginID, _, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
	if err != nil {
		t.Fatalf("BeginPasskeyLogin: %v", err)
	}
	if _, err := f.svc.FinishPasskeyRegistration(ctx, user.ID, loginID, "Mixed", []byte("{}")); !errors.Is(err, ErrPasskeyCeremonyExpired) {
		t.Errorf("a login challenge on the registration endpoint = %v, want ErrPasskeyCeremonyExpired", err)
	}
}

func TestPasskeyLoginRejectsBadAssertions(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)
	registerPasskey(t, f, user, authenticator, "Laptop")

	assert := func(t *testing.T, authenticator *fakeAuthenticator, user User) error {
		t.Helper()
		ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
		if err != nil {
			t.Fatalf("BeginPasskeyLogin: %v", err)
		}
		response, err := authenticator.assert(assertion.Response.Challenge.String(), userHandle(user))
		if err != nil {
			t.Fatalf("authenticator assert: %v", err)
		}
		_, err = f.svc.FinishPasskeyLogin(ctx, ceremonyID, response, testAddress)
		return err
	}

	t.Run("another origin", func(t *testing.T) {
		liar := *authenticator
		liar.origin = "https://attacker.example"
		if err := assert(t, &liar, user); !errors.Is(err, ErrInvalidPasskey) {
			t.Errorf("origin mismatch = %v, want ErrInvalidPasskey", err)
		}
	})

	t.Run("tampered signature", func(t *testing.T) {
		ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
		if err != nil {
			t.Fatalf("BeginPasskeyLogin: %v", err)
		}
		response, err := authenticator.assert(assertion.Response.Challenge.String(), userHandle(user))
		if err != nil {
			t.Fatalf("authenticator assert: %v", err)
		}
		var body struct {
			Response struct {
				Signature string `json:"signature"`
			} `json:"response"`
		}
		if err := json.Unmarshal(response, &body); err != nil {
			t.Fatalf("decode assertion: %v", err)
		}
		raw, err := base64.RawURLEncoding.DecodeString(body.Response.Signature)
		if err != nil {
			t.Fatalf("decode signature: %v", err)
		}
		raw[len(raw)-1] ^= 0xff
		var mutated map[string]any
		if err := json.Unmarshal(response, &mutated); err != nil {
			t.Fatalf("decode assertion: %v", err)
		}
		mutated["response"].(map[string]any)["signature"] = b64(raw)
		tampered, err := json.Marshal(mutated)
		if err != nil {
			t.Fatalf("encode assertion: %v", err)
		}
		if _, err := f.svc.FinishPasskeyLogin(ctx, ceremonyID, tampered, testAddress); !errors.Is(err, ErrInvalidPasskey) {
			t.Errorf("tampered signature = %v, want ErrInvalidPasskey", err)
		}
	})

	t.Run("unknown account handle", func(t *testing.T) {
		ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
		if err != nil {
			t.Fatalf("BeginPasskeyLogin: %v", err)
		}
		response, err := authenticator.assert(assertion.Response.Challenge.String(), []byte("no-such-account"))
		if err != nil {
			t.Fatalf("authenticator assert: %v", err)
		}
		if _, err := f.svc.FinishPasskeyLogin(ctx, ceremonyID, response, testAddress); !errors.Is(err, ErrInvalidPasskey) {
			t.Errorf("unknown handle = %v, want ErrInvalidPasskey", err)
		}
	})
}

func TestPasskeyLoginRefusesDisabledAccount(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)
	registerPasskey(t, f, user, authenticator, "Laptop")

	if err := f.svc.SetUserDisabled(ctx, user.ID, true); err != nil {
		t.Fatalf("SetUserDisabled: %v", err)
	}

	ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
	if err != nil {
		t.Fatalf("BeginPasskeyLogin: %v", err)
	}
	response, err := authenticator.assert(assertion.Response.Challenge.String(), userHandle(user))
	if err != nil {
		t.Fatalf("authenticator assert: %v", err)
	}
	if _, err := f.svc.FinishPasskeyLogin(ctx, ceremonyID, response, testAddress); !errors.Is(err, ErrInvalidPasskey) {
		t.Errorf("a disabled account signed in: %v", err)
	}
}

func TestRemovePasskey(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	user := f.setUp(t)
	ctx := t.Context()
	authenticator := newFakeAuthenticator(t)
	passkey := registerPasskey(t, f, user, authenticator, "Laptop")

	if err := f.svc.RemovePasskey(ctx, user.ID, passkey.ID); err != nil {
		t.Fatalf("RemovePasskey: %v", err)
	}
	list, err := f.svc.ListPasskeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListPasskeys: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("ListPasskeys after removing = %d, want 0", len(list))
	}

	ceremonyID, assertion, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
	if err != nil {
		t.Fatalf("BeginPasskeyLogin: %v", err)
	}
	response, err := authenticator.assert(assertion.Response.Challenge.String(), userHandle(user))
	if err != nil {
		t.Fatalf("authenticator assert: %v", err)
	}
	if _, err := f.svc.FinishPasskeyLogin(ctx, ceremonyID, response, testAddress); !errors.Is(err, ErrInvalidPasskey) {
		t.Errorf("a removed passkey still signed in: %v", err)
	}

	if err := f.svc.RemovePasskey(ctx, user.ID, passkey.ID); !errors.Is(err, ErrPasskeyNotFound) {
		t.Errorf("removing an unknown id = %v, want ErrPasskeyNotFound", err)
	}
}

func TestPasskeysUnavailableWhenNotConfigured(t *testing.T) {
	t.Parallel()

	f := newFixture(t)
	ctx := t.Context()

	if _, _, err := f.svc.BeginPasskeyLogin(ctx, testAddress); !errors.Is(err, ErrPasskeysUnavailable) {
		t.Errorf("BeginPasskeyLogin = %v, want ErrPasskeysUnavailable", err)
	}
	user := f.setUp(t)
	if _, _, err := f.svc.BeginPasskeyRegistration(ctx, user.ID); !errors.Is(err, ErrPasskeysUnavailable) {
		t.Errorf("BeginPasskeyRegistration = %v, want ErrPasskeysUnavailable", err)
	}
}

func TestPasskeyLoginThrottlesFailures(t *testing.T) {
	t.Parallel()

	f := newFixture(t, WithPasskeys(testPasskeys(t)))
	f.setUp(t)
	ctx := t.Context()

	throttled := false
	for range defaultIPLimit + 2 {
		ceremonyID, _, err := f.svc.BeginPasskeyLogin(ctx, testAddress)
		if errors.Is(err, ErrRateLimited) {
			throttled = true
			break
		}
		if err != nil {
			t.Fatalf("BeginPasskeyLogin: %v", err)
		}
		_, err = f.svc.FinishPasskeyLogin(ctx, ceremonyID, []byte(`{}`), testAddress)
		switch {
		case errors.Is(err, ErrRateLimited):
			throttled = true
		case !errors.Is(err, ErrInvalidPasskey):
			t.Fatalf("FinishPasskeyLogin = %v, want ErrInvalidPasskey", err)
		}
		if throttled {
			break
		}
	}
	if !throttled {
		t.Errorf("the address was never throttled after %d failures", defaultIPLimit+2)
	}
}

func TestCeremonyStoreIsSingleUseAndPrunes(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_800_000_000, 0)
	store := newCeremonyStore()
	entry := ceremony{kind: ceremonyLogin}

	store.put("one", entry, now, time.Minute)
	store.put("stale", entry, now.Add(-2*time.Minute), time.Minute)

	got, ok := store.take("one", now)
	if !ok || got.kind != ceremonyLogin {
		t.Fatalf("take(one) = %v, %v, want the stored ceremony", got, ok)
	}
	if _, ok := store.take("one", now); ok {
		t.Error("the same challenge was served twice")
	}

	store.put("late", entry, now.Add(-time.Minute), time.Second)
	if _, ok := store.take("late", now); ok {
		t.Error("an expired challenge was served")
	}
	if _, ok := store.take("late", now); ok {
		t.Error("an expired challenge stayed after being refused")
	}
}
