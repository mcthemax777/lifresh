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
			ID:            m.ID,
			LocalID:       m.LocalID,
			ParentID:      m.ParentID,
			LocalParentID: m.LocalParentID,
			Name:          m.Name,
			UserID:        m.UserID,
			Color:         m.Color,
			Type:          define.FileTypeFolder,
			Order:         m.Order,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		},
	}
}

func (r *FolderRepo) FromDomain(d *domain.Folder) *model.Folder {
	return &model.Folder{
		ID:            define.IfZero(d.ID, db.NextID()),
		LocalID:       d.LocalID,
		ParentID:      d.ParentID,
		LocalParentID: d.LocalParentID,
		Name:          d.Name,
		UserID:        d.UserID,
		Color:         d.Color,
		Order:         d.Order,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
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

func (r *FolderRepo) FindByUser(userID define.SnowflakeID) ([]*domain.Folder, error) {
	return r.Find("user_id = ?", userID)
}

func (r *FolderRepo) FindById(id define.SnowflakeID) (*domain.Folder, error) {
	return r.BaseRepo.FindById(r, id)
}

func (r *FolderRepo) FindByLocalId(id string) (*domain.Folder, error) {
	return r.BaseRepo.FindByLocalId(r, id)
}

func (r *FolderRepo) FindByParentID(id define.SnowflakeID) ([]*domain.Folder, error) {
	return r.Find("parent_id = ?", id)
}

func (r *BaseRepo[M, D]) DeleteById(id define.SnowflakeID) error {
	return r.db.Delete(new(M), "id = ?", id).Error
}
