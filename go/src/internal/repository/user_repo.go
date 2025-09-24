package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo[model.User, domain.User]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{NewBaseRepo[model.User, domain.User](db)}
}

// Mapper 구현
func (r *UserRepo) ToDomain(m *model.User) *domain.User {
	return &domain.User{
		ID:         m.ID,
		AccountID:  m.AccountID,
		Nickname:   m.Nickname,
		Bio:        m.Bio,
		ProfileURL: m.ProfileURL,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (r *UserRepo) FromDomain(d *domain.User) *model.User {
	return &model.User{
		ID:         define.IfZero(d.ID, db.NextID()),
		AccountID:  d.AccountID,
		Nickname:   d.Nickname,
		Bio:        d.Bio,
		ProfileURL: d.ProfileURL,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

func (r *UserRepo) FindByID(id define.SnowflakeID) (*domain.User, error) {
	return r.BaseRepo.FindById(r, id)
}

func (r *UserRepo) FindByAccountID(id define.SnowflakeID) (*domain.User, error) {
	return r.First(r, "account_id = ?", id)
}

func (r *UserRepo) Save(d *domain.User) (*domain.User, error) {
	return r.BaseRepo.Save(r, d)
}

func (r *UserRepo) Update(d *domain.User) (*domain.User, error) {
	return r.BaseRepo.Update(r, d)
}

func (r *UserRepo) DeleteByID(id define.SnowflakeID) error {
	return r.Delete("id = ?", id)
}
