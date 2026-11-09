package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRoundTrip(t *testing.T) {
	Init("0123456789abcdef0123456789abcdef")
	tok, err := GenerateToken(42, "admin")
	if err != nil {
		t.Fatal(err)
	}
	id, role, err := VerifyToken(tok)
	if err != nil || id != 42 || role != "admin" {
		t.Fatalf("got %d %q %v", id, role, err)
	}
}

func TestRejectsBadTokens(t *testing.T) {
	Init("0123456789abcdef0123456789abcdef")

	other := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1, "exp": time.Now().Add(time.Hour).Unix()})
	forged, _ := other.SignedString([]byte("some-other-secret-some-other-secret"))
	expired, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1, "exp": time.Now().Add(-time.Hour).Unix()}).
		SignedString([]byte("0123456789abcdef0123456789abcdef"))
	noExp, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1}).
		SignedString([]byte("0123456789abcdef0123456789abcdef"))
	none, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"user_id": 1, "exp": time.Now().Add(time.Hour).Unix()}).
		SignedString(jwt.UnsafeAllowNoneSignatureType)

	for name, tok := range map[string]string{"forged": forged, "expired": expired, "noExp": noExp, "none": none, "junk": "abc"} {
		if _, _, err := VerifyToken(tok); err == nil {
			t.Errorf("%s token accepted", name)
		}
	}
}
