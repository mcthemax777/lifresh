package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type RepeatRuleRepo struct {
	*BaseRepo[model.RepeatRule, domain.RepeatRule]
}

func NewRepeatRuleRepo(db *gorm.DB) *RepeatRuleRepo {
	return &RepeatRuleRepo{NewBaseRepo[model.RepeatRule, domain.RepeatRule](db)}
}

func (r *RepeatRuleRepo) ToDomain(m *model.RepeatRule) *domain.RepeatRule {
	return &domain.RepeatRule{
		ID:             m.ID,
		LocalID:        m.LocalID,
		PlanID:         m.PlanID,
		LocalPlanID:    m.LocalPlanID,
		Unit:           define.RepeatUnit(m.Unit),
		Interval:       m.Interval,
		PerRepeatCount: m.PerRepeatCount,
		Months:         m.Months,
		Days:           m.Days,
		Weeks:          m.Weeks,
		Weekdays:       m.Weekdays,
		Times:          m.Times,
		StartDate:      m.StartDate,
		EndDate:        m.EndDate,
		Order:          m.Order,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (r *RepeatRuleRepo) FromDomain(d *domain.RepeatRule) *model.RepeatRule {
	return &model.RepeatRule{
		ID:             define.IfZero(d.ID, db.NextID()),
		PlanID:         d.PlanID,
		Unit:           model.RepeatUnit(d.Unit),
		Interval:       d.Interval,
		PerRepeatCount: d.PerRepeatCount,
		Months:         d.Months,
		Days:           d.Days,
		Weeks:          d.Weeks,
		Weekdays:       d.Weekdays,
		Times:          d.Times,
		StartDate:      d.StartDate,
		EndDate:        d.EndDate,
		Order:          d.Order,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

func (r *RepeatRuleRepo) First(conds ...any) (*domain.RepeatRule, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *RepeatRuleRepo) Find(conds ...any) ([]*domain.RepeatRule, error) {
	return r.BaseRepo.Find(r, conds...)
}
func (r *RepeatRuleRepo) Save(d *domain.RepeatRule) (*domain.RepeatRule, error) {
	return r.BaseRepo.Save(r, d)
}
func (r *RepeatRuleRepo) Update(d *domain.RepeatRule) (*domain.RepeatRule, error) {
	return r.BaseRepo.Update(r, d)
}
func (r *RepeatRuleRepo) Delete(conds ...any) error { return r.BaseRepo.Delete(conds...) }
func (r *RepeatRuleRepo) FindByPlanID(id define.SnowflakeID) ([]*domain.RepeatRule, error) {
	return r.BaseRepo.Find(r, "plan_id = ?", id)
}
