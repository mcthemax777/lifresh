package repository

import (
	"gorm.io/gorm"
	"lifresh/internal/domain"
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

	result := r.dbConn.Create(&f)

	if result.Error != nil {
		return f, result.Error
	}

	return f, nil
}
