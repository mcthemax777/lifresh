package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type AccountRepo struct {
	*BaseRepo[model.Account, domain.Account]
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{NewBaseRepo[model.Account, domain.Account](db)}
}

func (r *AccountRepo) ToDomain(m *model.Account) *domain.Account {
	return &domain.Account{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		PhotoURL:    m.PhotoURL,
		SocialType:  define.SocialType(m.SocialType),
		ProviderUID: m.ProviderUID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *AccountRepo) FromDomain(d *domain.Account) *model.Account {
	id := define.IfZero(d.ID, db.NextID())
	return &model.Account{
		ID:          id,
		Name:        d.Name,
		Email:       d.Email,
		PhotoURL:    d.PhotoURL,
		SocialType:  int8(d.SocialType),
		ProviderUID: d.ProviderUID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// CRUD methods
func (r *AccountRepo) FindByID(id define.SnowflakeID) (*domain.Account, error) {
	return r.BaseRepo.FindById(r, id)
}

func (r *AccountRepo) FindByUID(uid string) (*domain.Account, error) {
	return r.First(r, "provider_uid = ?", uid)
}

func (r *AccountRepo) FindAll() ([]*domain.Account, error) {
	return r.Find(r)
}

func (r *AccountRepo) Save(d *domain.Account) (*domain.Account, error) {
	return r.BaseRepo.Save(r, d)
}

func (r *AccountRepo) Update(d *domain.Account) (*domain.Account, error) {
	return r.BaseRepo.Update(r, d)
}

func (r *AccountRepo) DeleteByUID(uid string) error {
	return r.Delete("provider_uid = ?", uid)
}
