package repository

import (
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"
)

type OptionItemRepo struct {
	*BaseRepo[model.OptionItem, domain.OptionItem]
}

func NewOptionItemRepo(db *gorm.DB) *OptionItemRepo {
	return &OptionItemRepo{NewBaseRepo[model.OptionItem, domain.OptionItem](db)}
}

func (r *OptionItemRepo) ToDomain(m *model.OptionItem) *domain.OptionItem {
	return &domain.OptionItem{
		ID:                 m.ID,
		LocalID:            m.LocalID,
		PlanID:             m.PlanID,
		LocalPlanID:        m.LocalPlanID,
		ParentID:           m.ParentID,
		LocalParentID:      m.LocalParentID,
		RecordFieldID:      m.RecordFieldID,
		LocalRecordFieldID: m.LocalRecordFieldID,
		Name:               m.Name,
		Order:              m.Order,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func (r *OptionItemRepo) FromDomain(d *domain.OptionItem) *model.OptionItem {
	return &model.OptionItem{
		ID:                 define.IfZero(d.ID, db.NextID()),
		LocalID:            d.LocalID,
		PlanID:             d.PlanID,
		LocalPlanID:        d.LocalPlanID,
		ParentID:           d.ParentID,
		LocalParentID:      d.LocalParentID,
		RecordFieldID:      d.RecordFieldID,
		LocalRecordFieldID: d.LocalRecordFieldID,
		Name:               d.Name,
		Order:              d.Order,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

func (r *OptionItemRepo) First(conds ...any) (*domain.OptionItem, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *OptionItemRepo) Find(conds ...any) ([]*domain.OptionItem, error) {
	return r.BaseRepo.Find(r, conds...)
}
func (r *OptionItemRepo) Save(d *domain.OptionItem) (*domain.OptionItem, error) {
	return r.BaseRepo.Save(r, d)
}
func (r *OptionItemRepo) Update(d *domain.OptionItem) (*domain.OptionItem, error) {
	return r.BaseRepo.Update(r, d)
}
func (r *OptionItemRepo) Delete(conds ...any) error { return r.BaseRepo.Delete(conds...) }
func (r *OptionItemRepo) FindByRecordFieldID(id define.SnowflakeID) ([]*domain.OptionItem, error) {
	return r.BaseRepo.Find(r, "record_field_id = ?", id)
}
