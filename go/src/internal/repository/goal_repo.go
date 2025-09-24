package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type GoalRepo struct {
	*BaseRepo[model.Goal, domain.Goal]
}

func NewGoalRepo(db *gorm.DB) *GoalRepo {
	return &GoalRepo{NewBaseRepo[model.Goal, domain.Goal](db)}
}

func (r *GoalRepo) ToDomain(m *model.Goal) *domain.Goal {
	return &domain.Goal{
		ID:         m.ID,
		PlanID:     m.PlanID,
		CycleType:  define.GoalRepeatCycleType(m.CycleType),
		Interval:   m.Interval,
		RecordID:   m.RecordID,
		Count:      m.Count,
		StartDate:  m.StartDate,
		FinishDate: m.FinishDate,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (r *GoalRepo) FromDomain(d *domain.Goal) *model.Goal {
	return &model.Goal{
		ID:         define.IfZero(d.ID, db.NextID()),
		PlanID:     d.PlanID,
		CycleType:  model.GoalRepeatCycleType(d.CycleType),
		Interval:   d.Interval,
		RecordID:   d.RecordID,
		Count:      d.Count,
		StartDate:  d.StartDate,
		FinishDate: d.FinishDate,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

func (r *GoalRepo) First(conds ...any) (*domain.Goal, error)    { return r.BaseRepo.First(r, conds...) }
func (r *GoalRepo) Find(conds ...any) ([]*domain.Goal, error)   { return r.BaseRepo.Find(r, conds...) }
func (r *GoalRepo) Save(d *domain.Goal) (*domain.Goal, error)   { return r.BaseRepo.Save(r, d) }
func (r *GoalRepo) Update(d *domain.Goal) (*domain.Goal, error) { return r.BaseRepo.Update(r, d) }
func (r *GoalRepo) Delete(conds ...any) error                   { return r.BaseRepo.Delete(conds...) }
