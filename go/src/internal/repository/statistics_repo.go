package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type StatisticsRepo struct {
	*BaseRepo[model.Statistics, domain.Statistics]
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{NewBaseRepo[model.Statistics, domain.Statistics](db)}
}

func (r *StatisticsRepo) ToDomain(m *model.Statistics) *domain.Statistics {
	return &domain.Statistics{
		ID:            m.ID,
		PlanID:        m.PlanID,
		Name:          m.Name,
		ChartType:     define.StatisticsChartType(m.ChartType),
		XAxisType:     m.XAxisType,
		XAxisPath:     m.XAxisPath,
		YAxisField:    m.YAxisField,
		XAxisIsRepeat: m.XAxisIsRepeat,
		YAxisIsRepeat: m.YAxisIsRepeat,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (r *StatisticsRepo) FromDomain(d *domain.Statistics) *model.Statistics {
	return &model.Statistics{
		ID:            define.IfZero(d.ID, db.NextID()),
		PlanID:        d.PlanID,
		Name:          d.Name,
		ChartType:     model.StatisticsChartType(d.ChartType),
		XAxisType:     d.XAxisType,
		XAxisPath:     d.XAxisPath,
		YAxisField:    d.YAxisField,
		XAxisIsRepeat: d.XAxisIsRepeat,
		YAxisIsRepeat: d.YAxisIsRepeat,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func (r *StatisticsRepo) First(conds ...any) (*domain.Statistics, error) {
	return r.BaseRepo.First(r, conds...)
}
func (r *StatisticsRepo) Find(conds ...any) ([]*domain.Statistics, error) {
	return r.BaseRepo.Find(r, conds...)
}
func (r *StatisticsRepo) Save(d *domain.Statistics) (*domain.Statistics, error) {
	return r.BaseRepo.Save(r, d)
}
func (r *StatisticsRepo) Update(d *domain.Statistics) (*domain.Statistics, error) {
	return r.BaseRepo.Update(r, d)
}
func (r *StatisticsRepo) Delete(conds ...any) error { return r.BaseRepo.Delete(conds...) }
