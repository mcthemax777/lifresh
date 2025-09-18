package repository

import (
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"
	"time"
)

type FolderRepo struct{ dbConn *gorm.DB }

func NewFolderRepo(dbConn *gorm.DB) *FolderRepo {
	return &FolderRepo{dbConn: dbConn}
}

func (r *FolderRepo) FindByID(id int64) (*domain.Folder, error) {
	var u *domain.Folder
	return u, nil
}
func (r *FolderRepo) Save(f *domain.Folder) (*domain.Folder, error) {

	if err := r.dbConn.Error; err != nil {
		return f, err
	}

	a := &model.Folder{
		ID:        db.NextID(),
		ParentID:  f.ParentID,
		UserID:    f.UserID,
		Order:     f.Order,
		Color:     f.Color,
		Name:      f.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	result := r.dbConn.Create(&a)

	if result.Error != nil {
		return f, result.Error
	}

	f.ID = a.ID
	return f, nil
}
