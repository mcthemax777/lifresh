package repository

import (
	"gorm.io/gorm"
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"
)

type RecordFieldRepo struct {
	*BaseRepo[model.RecordField, domain.RecordField]
}

func NewRecordFieldRepo(db *gorm.DB) *RecordFieldRepo {
	return &RecordFieldRepo{NewBaseRepo[model.RecordField, domain.RecordField](db)}
}

func (r *RecordFieldRepo) ToDomain(m *model.RecordField) *domain.RecordField {
	return &domain.RecordField{
		ID:            m.ID,
		PlanID:        m.PlanID,
		Name:          m.Name,
		Type:          m.Type,
		IsRepeatField: m.IsRepeatField,
		Unit:          m.Unit,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (r *RecordFieldRepo) FromDomain(d *domain.RecordField) *model.RecordField {
	return &model.RecordField{
		ID:            define.IfZero(d.ID, db.NextID()),
		PlanID:        d.PlanID,
		Name:          d.Name,
		Type:          d.Type,
		IsRepeatField: d.IsRepeatField,
		Unit:          d.Unit,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func (r *RecordFieldRepo) First(conds ...any) (*domain.RecordField, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *RecordFieldRepo) Find(conds ...any) ([]*domain.RecordField, error) {
	return r.BaseRepo.Find(r, conds...)
}
func (r *RecordFieldRepo) Save(d *domain.RecordField) (*domain.RecordField, error) {
	return r.BaseRepo.Save(r, d)
}
func (r *RecordFieldRepo) Update(d *domain.RecordField) (*domain.RecordField, error) {
	return r.BaseRepo.Update(r, d)
}
func (r *RecordFieldRepo) Delete(conds ...any) error { return r.BaseRepo.Delete(conds...) }
