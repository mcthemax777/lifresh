package service

import (
	"database/sql"
	"gorm.io/gorm"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/repository"
	"lifresh/internal/txmgr"
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

		var err error
		findAccount, err = s.accountRepo.FindByUID(account.ProviderUID)
		if err != nil {
			var err error
			findAccount, err = s.accountRepo.Save(account)
			if err != nil {
				return apperr.New(202, "insert error", err)
			}

			user := &domain.User{ID: 0, AccountID: findAccount.ID, Nickname: findAccount.ProviderUID, Bio: "", UpdatedAt: time.Now()}
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
