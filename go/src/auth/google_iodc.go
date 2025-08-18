package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Google ID Token에서 가져올 표준/확장 클레임
type GoogleClaims struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// auth/google_oidc.go 에 추가
func VerifyGoogleIDTokenAllowList(ctx context.Context, idToken string, clientIDs []string) (*GoogleClaims, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}

	// 우선 ClientID 체크를 건너뛰고 직접 aud 검사
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	tok, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verify idToken: %w", err)
	}

	var std struct {
		Aud string `json:"aud"`
		GoogleClaims
	}
	if err := tok.Claims(&std); err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}

	ok := false
	for _, cid := range clientIDs {
		if std.Aud == cid {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("audience mismatch: %s", std.Aud)
	}

	if std.Subject == "" {
		return nil, errors.New("invalid claims: missing sub")
	}
	return &std.GoogleClaims, nil
}

// VerifyGoogleIDToken:
// - 구글 OIDC Provider 메타에서 JWKS를 자동 캐시/로테이션
// - iss, aud, exp 검증 (aud는 GOOGLE_OAUTH_CLIENT_ID와 일치해야)
// - 성공 시 클레임 반환
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleClaims, error) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("missing GOOGLE_OAUTH_CLIENT_ID")
	}

	// Provider는 내부적으로 JWKs를 주기적으로 갱신/캐시합니다.
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}

	// aud 검증용 설정
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
		// SkipIssuerCheck: false (기본: iss = https://accounts.google.com)
		// SkipExpiryCheck:  false (기본: 만료 자동 검증)
		// ClockSkew:        허용하려면 컨텍스트 deadline/tolerance로 조정
	})

	// 실제 검증 (서명/만료/iss/aud)
	idTok, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verify idToken: %w", err)
	}

	// 클레임 파싱
	var claims GoogleClaims
	if err := idTok.Claims(&claims); err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}
	// 필수 sub 확인
	if claims.Subject == "" {
		return nil, errors.New("invalid claims: missing sub")
	}
	return &claims, nil
}

// ContextWithTimeout: 네트워크 호출에 타임아웃을 주기 위한 helper
func ContextWithTimeout(seconds int) (context.Context, context.CancelFunc) {
	if seconds <= 0 {
		seconds = 5
	}
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}
