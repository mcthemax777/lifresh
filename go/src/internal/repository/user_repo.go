package repository

import (
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"
)

type UserRepo struct{ dbConn *gorm.DB }

func NewUserRepo(dbConn *gorm.DB) *UserRepo {
	return &UserRepo{dbConn: dbConn}
}

func (r *UserRepo) FindByID(id int64) (domain.User, error) {
	var u domain.User
	return u, nil
}
func (r *UserRepo) Save(u *domain.User) (*domain.User, error) {
	if err := r.dbConn.Error; err != nil {
		return u, err
	}

	a := &model.User{
		ID:         db.NextID(),
		AccountID:  u.AccountID,
		Nickname:   u.Nickname,
		Bio:        u.Bio,
		ProfileURL: u.ProfileURL,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}

	result := r.dbConn.Create(&a)

	if result.Error != nil {
		return u, result.Error
	}

	u.ID = a.ID
	return u, nil
}
