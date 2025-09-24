package request

import (
	"lifresh/internal/domain"
)

type Request interface {
}

type LoginReq struct {
	SocialType  int    `json:"socialType"`
	SocialToken string `json:"socialToken"`
}

type CreateUserReq struct {
	User domain.User `json:"user"`
}

type RefreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

type PingReq struct {
}

type GetUserReq struct {
	UserID string `json:"userId"`
}

type CreateFolderReq struct {
	Folder *domain.Folder `json:"folder"`
}

type UpdateFolderReq struct {
	Folder *domain.Folder `json:"folder"`
}

type DeleteFolderReq struct {
	Folder *domain.Folder `json:"folder"`
}

type GetAccountAllDataReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}
