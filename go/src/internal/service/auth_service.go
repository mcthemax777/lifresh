package service

import (
	"lifresh/auth"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/txmgr"
	"os"
	"strings"
)

type AuthService struct {
	txm *txmgr.Manager
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Auth(socialToken string, socialType define.SocialType) (*domain.Account, *apperr.AppError) {
	var account domain.Account

	// 1) Google idToken이 오면 OIDC로 검증
	var uid string
	var email string
	var name string
	//// 로그인 핸들러 내부 (요지)
	switch {
	case socialType == define.SocialTypeGuest:
		uid = "guest-" + socialToken
		email = uid
		name = uid
	case socialType == define.SocialTypeGoogle:
		ids := strings.Split(os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"), ",")
		ctx, cancel := auth.ContextWithTimeout(50)
		defer cancel()
		claims, err := auth.VerifyGoogleIDTokenAllowList(ctx, socialToken, ids)
		if err != nil {
			return nil, apperr.New(201, "invalid_google_idtoken", err)
			//return ResponseToByteArray(response.CreateFailResponse(201, "invalid_google_idtoken")), err
		}
		uid = "google-" + claims.Subject
		email = claims.Email
		name = claims.Name
		//picture := claims.Picture

	case socialType == define.SocialTypeApple:
		// 여러 client_id 허용 (iOS 번들ID, 웹 Service ID 등)
		ids := strings.Split(os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"), ",")
		ctx, cancel := auth.ContextWithTimeout(5)
		defer cancel()
		ac, err := auth.VerifyAppleIDTokenAllowList(ctx, socialToken, ids, "")
		if err != nil {
			return nil, apperr.New(201, "invalid_apple_idtoken", err)
			//return ResponseToByteArray(response.CreateFailResponse(201, "invalid_apple_idtoken")), err
		}
		uid = "apple-" + ac.Subject
		email = ac.Email
		name = "apple_none"
	}

	account.ProviderUID = uid
	account.SocialType = socialType
	account.Email = email
	account.Name = name

	return &account, nil
}
