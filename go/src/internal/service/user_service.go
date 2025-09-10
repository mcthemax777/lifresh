package service

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/repository"
	"lifresh/internal/txmgr"
	"strconv"
	"time"
)

type UserService struct {
	txm         *txmgr.Manager
	accountRepo *repository.AccountRepo
	userRepo    *repository.UserRepo
}

func NewUserService(txm *txmgr.Manager, accountRepo *repository.AccountRepo, userRepo *repository.UserRepo) *UserService {
	return &UserService{txm: txm, accountRepo: accountRepo, userRepo: userRepo}
}

func (s *UserService) Login(account *domain.Account) (*domain.Account, *apperr.AppError) {

	var findAccount *domain.Account

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		isNew := false
		//게스트 로그인인데 토큰 보냈으면 존재하는 유저라고 판단해야됨
		if account.Social == define.SocialTypeGuest && account.ProviderUID == "" {
			account.ProviderUID = "guest-" + uuid.New().String()
			account.Email = strconv.FormatUint(uint64(db.NextID()), 10)
			isNew = true
		} else {
			var err error
			findAccount, err = s.accountRepo.FindByUID(account.ProviderUID)
			if err != nil {
				if account.Social == define.SocialTypeGuest && account.ProviderUID != "" {
					return apperr.New(202, "guest find error", err)
				}
				isNew = true
			}
		}

		if isNew {
			findAccount, err := s.accountRepo.Save(account)
			if err != nil {
				return apperr.New(202, "insert error", err)
			}

			user := &domain.User{ID: findAccount.ID, AccountID: findAccount.ID, Nickname: findAccount.ProviderUID, Bio: "", UpdatedAt: time.Now()}
			user, err = s.userRepo.Save(user)
			if err != nil {
				return apperr.New(202, "insert error", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return findAccount, nil
}
