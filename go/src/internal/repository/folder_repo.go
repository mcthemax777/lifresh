package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type FolderRepo struct {
	*BaseRepo[model.Folder, domain.Folder]
}

func NewFolderRepo(db *gorm.DB) *FolderRepo {
	return &FolderRepo{NewBaseRepo[model.Folder, domain.Folder](db)}
}

func (r *FolderRepo) ToDomain(m *model.Folder) *domain.Folder {
	return &domain.Folder{
		SystemFile: domain.SystemFile{
			ID:        m.ID,
			Name:      m.Name,
			ParentID:  m.ParentID,
			UserID:    m.UserID,
			Order:     m.Order,
			Color:     m.Color,
			Type:      define.FileTypeFolder,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	}
}

func (r *FolderRepo) FromDomain(d *domain.Folder) *model.Folder {
	return &model.Folder{
		ID:        define.IfZero(d.ID, db.NextID()),
		Name:      d.Name,
		ParentID:  d.ParentID,
		UserID:    d.UserID,
		Order:     d.Order,
		Color:     d.Color,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *FolderRepo) First(conds ...any) (*domain.Folder, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *FolderRepo) Find(conds ...any) ([]*domain.Folder, error) {

	return r.BaseRepo.Find(r, conds...)
}
func (r *FolderRepo) Save(d *domain.Folder) (*domain.Folder, error)   { return r.BaseRepo.Save(r, d) }
func (r *FolderRepo) Update(d *domain.Folder) (*domain.Folder, error) { return r.BaseRepo.Update(r, d) }

func (r *FolderRepo) Delete(d *domain.Folder) error {
	return r.BaseRepo.Delete("user_id = ? AND id = ?", d.UserID, d.ID)
}

func (r *FolderRepo) FindByUser(userID define.SnowflakeID) ([]*domain.Folder, error) {
	return r.Find("user_id = ?", userID)
}

func (r *FolderRepo) FindById(id define.SnowflakeID) (*domain.Folder, error) {
	return r.BaseRepo.FindById(r, id)
}
