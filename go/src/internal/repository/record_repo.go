package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type RecordRepo struct {
	*BaseRepo[model.Record, domain.Record]
}

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{NewBaseRepo[model.Record, domain.Record](db)}
}

func (r *RecordRepo) ToDomain(m *model.Record) *domain.Record {
	return &domain.Record{
		ID:          m.ID,
		LocalID:     m.LocalID,
		PlanID:      m.PlanID,
		LocalPlanID: m.LocalPlanID,
		Values:      m.Values,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *RecordRepo) FromDomain(d *domain.Record) *model.Record {
	return &model.Record{
		ID:          define.IfZero(d.ID, db.NextID()),
		LocalID:     d.LocalID,
		PlanID:      d.PlanID,
		LocalPlanID: d.LocalPlanID,
		Values:      mcopy(d.Values),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// Values 복사 헬퍼
func mcopy(src map[string]any) model.JSONMap {
	dst := make(model.JSONMap)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (r *RecordRepo) First(conds ...any) (*domain.Record, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *RecordRepo) Find(conds ...any) ([]*domain.Record, error) {
	return r.BaseRepo.Find(r, conds...)
}
func (r *RecordRepo) Save(d *domain.Record) (*domain.Record, error)   { return r.BaseRepo.Save(r, d) }
func (r *RecordRepo) Update(d *domain.Record) (*domain.Record, error) { return r.BaseRepo.Update(r, d) }
func (r *RecordRepo) Delete(conds ...any) error                       { return r.BaseRepo.Delete(conds...) }
func (r *RecordRepo) FindByPlanID(id define.SnowflakeID) ([]*domain.Record, error) {
	return r.BaseRepo.Find(r, "plan_id = ?", id)
}
