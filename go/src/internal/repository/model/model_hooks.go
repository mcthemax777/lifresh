package model

import (
	"fmt"
	"gorm.io/gorm"
	"lifresh/define"
	"lifresh/internal/core"
)

// =====================================================
// Utility: Global Version Counter
// =====================================================

func nextGlobalVersion(tx *gorm.DB) (int64, error) {
	var last ChangeLog
	if err := tx.Order("global_version desc").Limit(1).Find(&last).Error; err != nil {
		return 0, err
	}
	return last.GlobalVersion + 1, nil
}

// =====================================================
// Base Hooks: for all versioned entities
// =====================================================

func bumpEntityVersions(tx *gorm.DB, table string, id any, actorID any, operation string, payload JSONMap) error {
	nextGV, err := nextGlobalVersion(tx)
	if err != nil {
		return err
	}

	// 1️⃣ EntityVersion +1, GlobalVersion = nextGV
	if err := tx.Model(table).
		Where("id = ?", id).
		Updates(map[string]any{
			"entity_version": gorm.Expr("entity_version + 1"),
			"global_version": nextGV,
			"updated_at":     core.Now(),
		}).Error; err != nil {
		return err
	}

	// 2️⃣ ChangeLog 기록
	log := ChangeLog{
		GlobalVersion: nextGV,
		EntityType:    table,
		EntityID:      id.(define.SnowflakeID),
		GroupType:     "",
		GroupID:       0,
		Operation:     operation,
		ActorID:       actorID.(define.SnowflakeID),
		Payload:       payload,
		CreatedAt:     core.Now(),
	}
	if err := tx.Create(&log).Error; err != nil {
		return err
	}

	// 3️⃣ Outbox로 Redis Streams에 발행 예정
	outbox := Outbox{
		EventID: fmt.Sprintf("evt-%d", nextGV),
		Stream:  "change_log",
		Payload: JSONMap{
			"entityType": table,
			"entityID":   id,
			"operation":  operation,
			"version":    nextGV,
		},
		CreatedAt: core.Now(),
	}
	return tx.Create(&outbox).Error
}

// =====================================================
// Record Hooks (Plan GroupVersion 연동)
// =====================================================

func (r *Record) AfterCreate(tx *gorm.DB) error {
	return updatePlanGroupVersion(tx, r.PlanID, "create", r)
}
func (r *Record) AfterUpdate(tx *gorm.DB) error {
	return updatePlanGroupVersion(tx, r.PlanID, "update", r)
}
func (r *Record) AfterDelete(tx *gorm.DB) error {
	return updatePlanGroupVersion(tx, r.PlanID, "delete", r)
}

// =====================================================
// Helper: Plan.GroupVersion & DataVersion 연동
// =====================================================

func updatePlanGroupVersion(tx *gorm.DB, planID any, op string, record *Record) error {
	// 1️⃣ Plan의 GroupVersion +1
	if err := tx.Model(&Plan{}).
		Where("id = ?", planID).
		Update("group_version", gorm.Expr("group_version + 1")).Error; err != nil {
		return err
	}

	// 2️⃣ DataVersion 테이블 업데이트
	var dv DataVersion
	if err := tx.First(&dv, "group_type = ? AND group_id = ?", "plan_records", planID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			dv = DataVersion{
				GroupType: "plan_records",
				GroupID:   planID.(define.SnowflakeID),
				Version:   1,
				UpdatedAt: core.Now(),
			}
			if err := tx.Create(&dv).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err := tx.Model(&DataVersion{}).
			Where("group_type = ? AND group_id = ?", "plan_records", planID).
			Updates(map[string]any{
				"version":    gorm.Expr("version + 1"),
				"updated_at": core.Now(),
			}).Error; err != nil {
			return err
		}
	}

	// 3️⃣ ChangeLog 추가
	nextGV, err := nextGlobalVersion(tx)
	if err != nil {
		return err
	}
	log := ChangeLog{
		GlobalVersion: nextGV,
		EntityType:    "record",
		EntityID:      record.ID,
		GroupType:     "plan_records",
		GroupID:       planID.(define.SnowflakeID),
		Operation:     op,
		ActorID:       0,
		Payload: JSONMap{
			"planId": planID,
		},
		CreatedAt: core.Now(),
	}
	if err := tx.Create(&log).Error; err != nil {
		return err
	}

	// 4️⃣ Outbox 발행
	outbox := Outbox{
		EventID: fmt.Sprintf("evt-%d", nextGV),
		Stream:  "plan_records",
		Payload: JSONMap{
			"planId":    planID,
			"operation": op,
			"recordId":  record.ID,
			"version":   nextGV,
			"groupVer":  gorm.Expr("group_version + 1"),
		},
		CreatedAt: core.Now(),
	}
	return tx.Create(&outbox).Error
}
