package repository

import (
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/domain"
	"lifresh/internal/repository/model"

	"gorm.io/gorm"
)

type PlanRepo struct {
	*BaseRepo[model.Plan, domain.Plan]
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{NewBaseRepo[model.Plan, domain.Plan](db)}
}

func (r *PlanRepo) ToDomain(m *model.Plan) *domain.Plan {
	return &domain.Plan{
		SystemFile: domain.SystemFile{
			ID:        m.ID,
			Name:      m.Name,
			ParentID:  m.ParentID,
			UserID:    m.UserID,
			Order:     m.Order,
			Color:     m.Color,
			Type:      define.FileTypePlan,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Description: m.Description,
		StartDate:   m.StartDate,
		FinishDate:  m.FinishDate,
		DateType:    define.DateType(m.DateType),
	}
}

func (r *PlanRepo) FromDomain(d *domain.Plan) *model.Plan {
	return &model.Plan{
		ID:          define.IfZero(d.ID, db.NextID()),
		Name:        d.Name,
		ParentID:    d.ParentID,
		UserID:      d.UserID,
		Order:       d.Order,
		Color:       d.Color,
		Description: d.Description,
		StartDate:   d.StartDate,
		FinishDate:  d.FinishDate,
		DateType:    model.DateType(d.DateType),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func (r *PlanRepo) First(conds ...any) (*domain.Plan, error)    { return r.BaseRepo.First(r, conds...) }
func (r *PlanRepo) Find(conds ...any) ([]*domain.Plan, error)   { return r.BaseRepo.Find(r, conds...) }
func (r *PlanRepo) Save(d *domain.Plan) (*domain.Plan, error)   { return r.BaseRepo.Save(r, d) }
func (r *PlanRepo) Update(d *domain.Plan) (*domain.Plan, error) { return r.BaseRepo.Update(r, d) }
func (r *PlanRepo) Delete(conds ...any) error                   { return r.BaseRepo.Delete(conds...) }

func (r *PlanRepo) FindByUser(userID define.SnowflakeID) ([]*domain.Plan, error) {
	return r.Find("user_id = ?", userID)
}
