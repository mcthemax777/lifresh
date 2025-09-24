package repository

import (
	"gorm.io/gorm"
	"lifresh/define"
)

// Mapper: model <-> domain 변환 인터페이스
type Mapper[M any, D any] interface {
	ToDomain(*M) *D
	FromDomain(*D) *M
}

// BaseRepo: 공통 CRUD
type BaseRepo[M any, D any] struct {
	db *gorm.DB
}

func NewBaseRepo[M any, D any](db *gorm.DB) *BaseRepo[M, D] {
	return &BaseRepo[M, D]{db: db}
}

// First: 조건에 맞는 단일 row
func (r *BaseRepo[M, D]) First(mapper Mapper[M, D], conds ...any) (*D, error) {
	var m M
	result := r.db.First(&m, conds...)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapper.ToDomain(&m), nil
}

// Find: 조건에 맞는 여러 row
func (r *BaseRepo[M, D]) Find(mapper Mapper[M, D], conds ...any) ([]*D, error) {
	var models []M
	result := r.db.Find(&models, conds...)
	if result.Error != nil {
		return nil, result.Error
	}
	domains := make([]*D, 0, len(models))
	for i := range models {
		domains = append(domains, mapper.ToDomain(&models[i]))
	}
	return domains, nil
}

func (r *BaseRepo[M, D]) FindById(mapper Mapper[M, D], id define.SnowflakeID) (*D, error) {
	return r.First(mapper, "id = ?", id)
}

// Save: 새 row 생성
func (r *BaseRepo[M, D]) Save(mapper Mapper[M, D], d *D) (*D, error) {
	m := mapper.FromDomain(d)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return mapper.ToDomain(m), nil
}

// Update: 전체 row 갱신
func (r *BaseRepo[M, D]) Update(mapper Mapper[M, D], d *D) (*D, error) {
	m := mapper.FromDomain(d)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return mapper.ToDomain(m), nil
}

// Delete: 조건에 맞는 row 삭제
func (r *BaseRepo[M, D]) Delete(conds ...any) error {
	return r.db.Delete(new(M), conds...).Error
}
