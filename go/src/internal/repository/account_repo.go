package repository

import (
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"
)

type AccountRepo struct{ dbConn *gorm.DB }

func NewAccountRepo(dbConn *gorm.DB) *AccountRepo {
	return &AccountRepo{dbConn: dbConn}
}

func (r *AccountRepo) FindByUID(uid string) (*domain.Account, error) {
	var a domain.Account

	r.dbConn.First(&a, "prov_ider_uid = ?", uid)

	if a.ID == 0 {
		return nil, apperr.New(202, "not find account", nil)
	}
	return &a, nil
}
func (r *AccountRepo) Save(account *domain.Account) (*domain.Account, error) {
	if err := r.dbConn.Error; err != nil {
		return account, err
	}

	a := &model.Account{
		ID:          db.NextID(),
		Name:        account.Name,
		Email:       account.Email,
		PhotoURL:    account.PhotoURL,
		SocialType:  int8(account.Social),
		ProviderUID: account.ProviderUID,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}

	result := r.dbConn.Create(&a)

	account.ID = a.ID
	if result.Error != nil {
		return account, result.Error
	}

	account.ID = a.ID

	return account, nil
}
