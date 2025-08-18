// auth/apple_oidc.go
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type AppleClaims struct {
	Subject        string `json:"sub"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"` // "true"/"false"
	IsPrivateEmail string `json:"is_private_email"`
	Aud            string `json:"aud"`
	Nonce          string `json:"nonce"`
}

// 여러 client_id 허용 + (선택) rawNonce 검증
func VerifyAppleIDTokenAllowList(ctx context.Context, idToken string, clientIDs []string, rawNonce string) (*AppleClaims, error) {
	provider, err := oidc.NewProvider(ctx, "https://appleid.apple.com")
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}
	// aud는 직접 검사할 것이므로 SkipClientIDCheck
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	tok, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}
	var c AppleClaims
	if err := tok.Claims(&c); err != nil {
		return nil, fmt.Errorf("claims: %w", err)
	}
	if c.Subject == "" {
		return nil, errors.New("missing sub")
	}

	// aud 허용 리스트 검사
	ok := false
	for _, cid := range clientIDs {
		if strings.TrimSpace(c.Aud) == strings.TrimSpace(cid) {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("audience mismatch: %s", c.Aud)
	}

	// (선택) nonce 검증: 토큰의 nonce는 rawNonce의 SHA-256 해시
	if rawNonce != "" && c.Nonce != "" {
		sum := sha256.Sum256([]byte(rawNonce))
		// Apple은 보통 base64url(무패딩) 사용, 환경에 따라 hex로 오는 경우도 대비
		base64url := base64.RawURLEncoding.EncodeToString(sum[:])
		hexsum := hex.EncodeToString(sum[:])
		if c.Nonce != base64url && !strings.EqualFold(c.Nonce, hexsum) {
			return nil, errors.New("nonce mismatch")
		}
	}

	return &c, nil
}
