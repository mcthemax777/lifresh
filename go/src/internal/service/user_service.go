package service

import (
	"database/sql"
	"gorm.io/gorm"
	"lifresh/define"
	"lifresh/internal/core"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/repository"
	"lifresh/internal/txmgr"
	"sort"
)

type UserService struct {
	txm             *txmgr.Manager
	accountRepo     *repository.AccountRepo
	userRepo        *repository.UserRepo
	folderRepo      *repository.FolderRepo
	planRepo        *repository.PlanRepo
	recordFieldRepo *repository.RecordFieldRepo
	goalRepo        *repository.GoalRepo
	repeatRuleRepo  *repository.RepeatRuleRepo
	optionItemRepo  *repository.OptionItemRepo
	recordRepo      *repository.RecordRepo
	statisticsRepo  *repository.StatisticsRepo
}

func NewUserService(txm *txmgr.Manager,
	accountRepo *repository.AccountRepo,
	userRepo *repository.UserRepo,
	folderRepo *repository.FolderRepo,
	planRepo *repository.PlanRepo,
	recordFieldRepo *repository.RecordFieldRepo,
	goalRepo *repository.GoalRepo,
	repeatRuleRepo *repository.RepeatRuleRepo,
	optionItemRepo *repository.OptionItemRepo,
	recordRepo *repository.RecordRepo,
	statisticsRepo *repository.StatisticsRepo,
) *UserService {
	return &UserService{txm: txm, accountRepo: accountRepo, userRepo: userRepo, folderRepo: folderRepo, planRepo: planRepo,
		recordFieldRepo: recordFieldRepo, goalRepo: goalRepo, repeatRuleRepo: repeatRuleRepo, statisticsRepo: statisticsRepo, optionItemRepo: optionItemRepo, recordRepo: recordRepo}
}

func (s *UserService) Login(account *domain.Account) (*domain.Account, *apperr.AppError) {

	var findAccount *domain.Account

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		var err error
		findAccount, err = s.accountRepo.FindByUID(account.ProviderUID)
		if err != nil {
			var err error
			findAccount, err = s.accountRepo.Save(account)
			if err != nil {
				return apperr.New(202, "insert error", err)
			}

		} else {
			user, _ := s.userRepo.FindByAccountID(findAccount.ID)
			findAccount.User = user
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return findAccount, nil
}

func (s *UserService) GetAccountByID(accountID define.SnowflakeID) (*domain.Account, *apperr.AppError) {

	var findAccount *domain.Account

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		var err error
		findAccount, err = s.accountRepo.FindByID(accountID)
		if err != nil {
			return apperr.New(202, "insert error", err)
		}

		user, _ := s.userRepo.FindByAccountID(findAccount.ID)
		findAccount.User = user

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return findAccount, nil
}

func (s *UserService) GetUserByAccountID(accountID define.SnowflakeID) (*domain.User, *apperr.AppError) {

	var user *domain.User

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		var err error
		if user, err = s.userRepo.FindByAccountID(accountID); err != nil {
			return apperr.New(202, "find user error", err)
		}
		if err := s.getData(user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return user, nil
}

func (s *UserService) GetUser(userID define.SnowflakeID) (*domain.User, *apperr.AppError) {

	var user *domain.User

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		var err error

		// 유저
		if user, err = s.userRepo.FindByID(userID); err != nil {
			return apperr.New(202, "find user error", err)
		}

		if err := s.getData(user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return user, nil
}

func (s *UserService) CreateUser(user *domain.User) (*domain.User, *apperr.AppError) {

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		var err error

		user.CreatedAt = core.Now()
		user.UpdatedAt = core.Now()

		user, err = s.userRepo.Save(user)
		if err != nil {
			return apperr.New(202, "insert error", err)
		}

		root := &domain.Folder{
			SystemFile: domain.SystemFile{
				ID:        0,
				Name:      "최상위 폴더",
				ParentID:  0,
				UserID:    user.ID,
				Order:     0,
				Color:     0,
				Type:      define.FileTypeFolder,
				CreatedAt: core.Now(),
				UpdatedAt: core.Now(),
			},
		}
		root, err = s.folderRepo.Save(root)
		if err != nil {
			return apperr.New(202, "insert error", err)
		}

		user.Root = root

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return user, nil
}

func (s *UserService) UpdateUser(user *domain.User) (*domain.User, *apperr.AppError) {

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		var err error

		user.UpdatedAt = core.Now()
		user, err = s.userRepo.Update(user)
		if err != nil {
			return apperr.New(202, "Update error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return user, nil
}

func (s *UserService) CreateFolder(folder *domain.Folder) (*domain.Folder, *apperr.AppError) {

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		//TODO. 유저 정보 일치 확인, 부모 폴더가 존재하는지 확인

		var parentFolder *domain.Folder
		var err error
		if folder.ParentID == 0 {
			parentFolder, err = s.folderRepo.FindByLocalId(folder.LocalParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}

			folder.ParentID = parentFolder.ID
		} else {
			parentFolder, err = s.folderRepo.FindById(folder.ParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}
		}

		folder.UserID = parentFolder.UserID
		folder, err = s.folderRepo.Save(folder)
		if err != nil {
			return apperr.New(202, "save folder error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return folder, nil
}

func (s *UserService) UpdateFolder(folder *domain.Folder) (*domain.Folder, *apperr.AppError) {

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		//TODO. 유저 정보 일치 확인, 부모 폴더가 존재하는지 확인

		// 아이디가 없으면 로컬아이디로 탐색
		if folder.ID == 0 {
			findFolder, err := s.folderRepo.FindByLocalId(folder.LocalID)
			if err != nil {
				return apperr.New(202, "not find folder error", err)
			}

			folder.ID = findFolder.ID
		}

		var parentFolder *domain.Folder
		var err error
		if folder.ParentID == 0 {
			parentFolder, err = s.folderRepo.FindByLocalId(folder.LocalParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}

			folder.ParentID = parentFolder.ID
		} else {
			parentFolder, err = s.folderRepo.FindById(folder.ParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}
		}

		folder.UserID = parentFolder.UserID

		folder, err = s.folderRepo.Update(folder)
		if err != nil {
			return apperr.New(202, "save folder error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return folder, nil
}

func (s *UserService) DeleteFolder(folder *domain.Folder) (*domain.Folder, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {

		// 아이디가 없으면 로컬아이디로 탐색
		if folder.ID == 0 {
			findFolder, err := s.folderRepo.FindByLocalId(folder.LocalID)
			if err != nil {
				return apperr.New(202, "not find folder error", err)
			}

			folder.ID = findFolder.ID
		}

		if err := s.deleteFolder(folder); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return folder, nil
}

func (s *UserService) deleteFolder(folder *domain.Folder) *apperr.AppError {

	// 1. 자식 폴더들 재귀 삭제
	childFolders, err := s.folderRepo.FindByParentID(folder.ID)
	if err != nil {
		return apperr.New(202, "find child folders error", err)
	}
	for _, cf := range childFolders {
		if err := s.deleteFolder(cf); err != nil {
			return err
		}
	}

	// 2. 자식 플랜들 삭제
	childPlans, err := s.planRepo.FindByParentID(folder.ID)
	if err != nil {
		return apperr.New(202, "find child plans error", err)
	}
	for _, cp := range childPlans {
		if err := s.deletePlan(cp); err != nil {
			return err
		}
	}

	// 3. 폴더 자체 삭제
	if err := s.folderRepo.DeleteById(folder.ID); err != nil {
		return apperr.New(202, "delete folder error", err)
	}

	return nil
}

func (s *UserService) CreatePlan(plan *domain.Plan) (*domain.Plan, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {

		if plan.ParentID == 0 {
			parentFolder, err := s.folderRepo.FindByLocalId(plan.LocalParentID)
			if err != nil {
				return apperr.New(202, "FindByLocalId error", err)
			}

			plan.ParentID = parentFolder.ID
		}

		// 1. 부모 폴더 확인 및 UserID 설정
		parentFolder, err := s.folderRepo.FindById(plan.ParentID)
		if err != nil {
			return apperr.New(202, "save plan error", err)
		}
		plan.UserID = parentFolder.UserID
		plan.CreatedAt = core.Now()
		plan.UpdatedAt = core.Now()

		// 2. Plan 저장
		newPlan, err := s.planRepo.Save(plan)
		if err != nil {
			return apperr.New(202, "save plan error", err)
		}

		plan.ID = newPlan.ID

		// 3. Goals 저장
		for i := range plan.Goals {
			plan.Goals[i].PlanID = plan.ID
			plan.Goals[i].CreatedAt = core.Now()
			plan.Goals[i].UpdatedAt = core.Now()
			if _, err := s.goalRepo.Save(plan.Goals[i]); err != nil {
				return apperr.New(202, "save goal error", err)
			}
		}

		// 4. RecordFields 저장 (+ OptionItem 트리)
		for i := range plan.RecordFields {
			rf := plan.RecordFields[i]
			rf.PlanID = plan.ID
			rf.CreatedAt = core.Now()
			rf.UpdatedAt = core.Now()

			newRF, err := s.recordFieldRepo.Save(rf)
			if err != nil {
				return apperr.New(202, "save record field error", err)
			}
			rf.ID = newRF.ID

			// OptionItem 트리 저장 (최상위 부모 ParentID=0)
			for j := range rf.Options {
				if err := s.saveOptionItemTree(rf.Options[j], rf.ID, 0); err != nil {
					return err
				}
			}
		}

		// 5. RepeatRules 저장
		for i := range plan.RepeatRules {
			plan.RepeatRules[i].PlanID = plan.ID
			plan.RepeatRules[i].CreatedAt = core.Now()
			plan.RepeatRules[i].UpdatedAt = core.Now()
			if _, err := s.repeatRuleRepo.Save(plan.RepeatRules[i]); err != nil {
				return apperr.New(202, "save repeat rule error", err)
			}
		}

		// 6. Statistics 저장
		for i := range plan.Statistics {
			plan.Statistics[i].PlanID = plan.ID
			plan.Statistics[i].CreatedAt = core.Now()
			plan.Statistics[i].UpdatedAt = core.Now()
			if _, err := s.statisticsRepo.Save(plan.Statistics[i]); err != nil {
				return apperr.New(202, "save statistics error", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return plan, nil
}

// OptionItem 트리 저장 (재귀)
func (s *UserService) saveOptionItemTree(opt *domain.OptionItem, recordFieldID define.SnowflakeID, parentID define.SnowflakeID) error {
	opt.RecordFieldID = recordFieldID
	opt.ParentID = parentID
	opt.CreatedAt = core.Now()
	opt.UpdatedAt = core.Now()

	// 1. 현재 노드 저장
	newOpt, err := s.optionItemRepo.Save(opt)
	if err != nil {
		return apperr.New(202, "save option item error", err)
	}
	opt.ID = newOpt.ID // 저장된 ID 반영

	// 2. 자식 노드들 저장
	for i := range opt.Children {
		if err := s.saveOptionItemTree(opt.Children[i], recordFieldID, opt.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) UpdatePlan(plan *domain.Plan) (*domain.Plan, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {

		// 아이디가 없으면 로컬아이디로 탐색
		if plan.ID == 0 {
			findPlan, err := s.planRepo.FindByLocalId(plan.LocalID)
			if err != nil {
				return apperr.New(202, "not find plan error", err)
			}

			plan.ID = findPlan.ID
		}

		var parentFolder *domain.Folder
		var err error
		if plan.ParentID == 0 {
			parentFolder, err = s.folderRepo.FindByLocalId(plan.LocalParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}

			plan.ParentID = parentFolder.ID
		} else {
			parentFolder, err = s.folderRepo.FindById(plan.ParentID)
			if err != nil {
				return apperr.New(202, "save folder error", err)
			}
		}

		plan.UserID = parentFolder.UserID

		// 1. Plan 업데이트
		plan.UpdatedAt = core.Now()
		_, err = s.planRepo.Update(plan)
		if err != nil {
			return apperr.New(202, "update plan error", err)
		}

		// =========================
		// Goals
		// =========================
		existingGoals, _ := s.goalRepo.FindByPlanID(plan.ID)
		goalMap := make(map[define.SnowflakeID]bool)

		for i := range plan.Goals {
			g := plan.Goals[i]
			g.PlanID = plan.ID
			g.UpdatedAt = core.Now()

			if g.ID == 0 {
				g.CreatedAt = core.Now()
				if _, err := s.goalRepo.Save(g); err != nil {
					return apperr.New(202, "save goal error", err)
				}
			} else {
				if _, err := s.goalRepo.Update(g); err != nil {
					return apperr.New(202, "update goal error", err)
				}
			}
			goalMap[g.ID] = true
		}
		for _, eg := range existingGoals {
			if _, ok := goalMap[eg.ID]; !ok {
				if err := s.goalRepo.DeleteById(eg.ID); err != nil {
					return apperr.New(202, "delete goal error", err)
				}
			}
		}

		// =========================
		// RecordFields (+OptionItems)
		// =========================
		existingFields, _ := s.recordFieldRepo.FindByPlanID(plan.ID)
		fieldMap := make(map[define.SnowflakeID]bool)

		for i := range plan.RecordFields {
			rf := plan.RecordFields[i]
			rf.PlanID = plan.ID
			rf.UpdatedAt = core.Now()

			if rf.ID == 0 {
				rf.CreatedAt = core.Now()
				newRF, err := s.recordFieldRepo.Save(rf)
				if err != nil {
					return apperr.New(202, "save record field error", err)
				}
				rf.ID = newRF.ID
			} else {
				if _, err := s.recordFieldRepo.Update(rf); err != nil {
					return apperr.New(202, "update record field error", err)
				}
			}
			fieldMap[rf.ID] = true

			// OptionItems diff: DB에 있는 애들 목록
			existingOpts, _ := s.optionItemRepo.FindByRecordFieldID(rf.ID)
			optMap := make(map[define.SnowflakeID]bool)

			// 클라에서 온 OptionItem 트리 저장
			for j := range rf.Options {
				if err := s.upsertOptionItemTree(rf.Options[j], rf.ID, 0, optMap); err != nil {
					return err
				}
				optMap[rf.Options[j].ID] = true
			}

			// 기존 DB 중 요청에 없는 건 삭제
			for _, eo := range existingOpts {
				if _, ok := optMap[eo.ID]; !ok {
					if err := s.optionItemRepo.DeleteById(eo.ID); err != nil {
						return apperr.New(202, "delete option item error", err)
					}
				}
			}
		}
		for _, ef := range existingFields {
			if _, ok := fieldMap[ef.ID]; !ok {
				if err := s.recordFieldRepo.DeleteById(ef.ID); err != nil {
					return apperr.New(202, "delete record field error", err)
				}
			}
		}

		// =========================
		// RepeatRules
		// =========================
		existingRules, _ := s.repeatRuleRepo.FindByPlanID(plan.ID)
		ruleMap := make(map[define.SnowflakeID]bool)

		for i := range plan.RepeatRules {
			rr := plan.RepeatRules[i]
			rr.PlanID = plan.ID
			rr.UpdatedAt = core.Now()

			if rr.ID == 0 {
				rr.CreatedAt = core.Now()
				if _, err := s.repeatRuleRepo.Save(rr); err != nil {
					return apperr.New(202, "save repeat rule error", err)
				}
			} else {
				if _, err := s.repeatRuleRepo.Update(rr); err != nil {
					return apperr.New(202, "update repeat rule error", err)
				}
			}
			ruleMap[rr.ID] = true
		}
		for _, er := range existingRules {
			if _, ok := ruleMap[er.ID]; !ok {
				if err := s.repeatRuleRepo.DeleteById(er.ID); err != nil {
					return apperr.New(202, "delete repeat rule error", err)
				}
			}
		}

		// =========================
		// Statistics
		// =========================
		existingStats, _ := s.statisticsRepo.FindByPlanID(plan.ID)
		statMap := make(map[define.SnowflakeID]bool)

		for i := range plan.Statistics {
			st := plan.Statistics[i]
			st.PlanID = plan.ID
			st.UpdatedAt = core.Now()

			if st.ID == 0 {
				st.CreatedAt = core.Now()
				if _, err := s.statisticsRepo.Save(st); err != nil {
					return apperr.New(202, "save statistics error", err)
				}
			} else {
				if _, err := s.statisticsRepo.Update(st); err != nil {
					return apperr.New(202, "update statistics error", err)
				}
			}
			statMap[st.ID] = true
		}
		for _, es := range existingStats {
			if _, ok := statMap[es.ID]; !ok {
				if err := s.statisticsRepo.DeleteById(es.ID); err != nil {
					return apperr.New(202, "delete statistics error", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return plan, nil
}

func (s *UserService) DeletePlan(plan *domain.Plan) (*domain.Plan, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {

		// 아이디가 없으면 로컬아이디로 탐색
		if plan.ID == 0 {
			findPlan, err := s.planRepo.FindByLocalId(plan.LocalID)
			if err != nil {
				return apperr.New(202, "not find plan error", err)
			}

			plan.ID = findPlan.ID
		}

		// 1. Goals 삭제
		goals, _ := s.goalRepo.FindByPlanID(plan.ID)
		for _, g := range goals {
			if err := s.goalRepo.DeleteById(g.ID); err != nil {
				return apperr.New(202, "delete goal error", err)
			}
		}

		// 2. RecordFields & OptionItems 삭제
		fields, _ := s.recordFieldRepo.FindByPlanID(plan.ID)
		for _, rf := range fields {
			opts, _ := s.optionItemRepo.FindByRecordFieldID(rf.ID)
			for _, o := range opts {
				if err := s.optionItemRepo.DeleteById(o.ID); err != nil {
					return apperr.New(202, "delete option item error", err)
				}
			}
			if err := s.recordFieldRepo.DeleteById(rf.ID); err != nil {
				return apperr.New(202, "delete record field error", err)
			}
		}

		// 3. RepeatRules 삭제
		rules, _ := s.repeatRuleRepo.FindByPlanID(plan.ID)
		for _, rr := range rules {
			if err := s.repeatRuleRepo.DeleteById(rr.ID); err != nil {
				return apperr.New(202, "delete repeat rule error", err)
			}
		}

		// 4. Statistics 삭제
		stats, _ := s.statisticsRepo.FindByPlanID(plan.ID)
		for _, st := range stats {
			if err := s.statisticsRepo.DeleteById(st.ID); err != nil {
				return apperr.New(202, "delete statistics error", err)
			}
		}

		// 5. Records 삭제
		records, _ := s.recordRepo.FindByPlanID(plan.ID)
		for _, r := range records {
			if err := s.recordRepo.DeleteById(r.ID); err != nil {
				return apperr.New(202, "delete record error", err)
			}
		}

		// 6. Plan 삭제
		if err := s.planRepo.DeleteById(plan.ID); err != nil {
			return apperr.New(202, "delete plan error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return plan, nil
}

func (s *UserService) deletePlan(plan *domain.Plan) *apperr.AppError {

	// 1. Goals 삭제
	goals, _ := s.goalRepo.FindByPlanID(plan.ID)
	for _, g := range goals {
		if err := s.goalRepo.DeleteById(g.ID); err != nil {
			return apperr.New(202, "delete goal error", err)
		}
	}

	// 2. RecordFields & OptionItems 삭제
	fields, _ := s.recordFieldRepo.FindByPlanID(plan.ID)
	for _, rf := range fields {
		opts, _ := s.optionItemRepo.FindByRecordFieldID(rf.ID)
		for _, o := range opts {
			if err := s.optionItemRepo.DeleteById(o.ID); err != nil {
				return apperr.New(202, "delete option item error", err)
			}
		}
		if err := s.recordFieldRepo.DeleteById(rf.ID); err != nil {
			return apperr.New(202, "delete record field error", err)
		}
	}

	// 3. RepeatRules 삭제
	rules, _ := s.repeatRuleRepo.FindByPlanID(plan.ID)
	for _, rr := range rules {
		if err := s.repeatRuleRepo.DeleteById(rr.ID); err != nil {
			return apperr.New(202, "delete repeat rule error", err)
		}
	}

	// 4. Statistics 삭제
	stats, _ := s.statisticsRepo.FindByPlanID(plan.ID)
	for _, st := range stats {
		if err := s.statisticsRepo.DeleteById(st.ID); err != nil {
			return apperr.New(202, "delete statistics error", err)
		}
	}

	// 5. Records 삭제
	records, _ := s.recordRepo.FindByPlanID(plan.ID)
	for _, r := range records {
		if err := s.recordRepo.DeleteById(r.ID); err != nil {
			return apperr.New(202, "delete record error", err)
		}
	}

	// 6. Plan 삭제
	if err := s.planRepo.DeleteById(plan.ID); err != nil {
		return apperr.New(202, "delete plan error", err)
	}

	return nil
}

func (s *UserService) CreateRecord(record *domain.Record) (*domain.Record, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {

		// Plan 존재 여부 및 UserID 확인
		plan, err := s.planRepo.FindById(record.PlanID)
		if err != nil {
			return apperr.New(202, "plan not found", err)
		}

		// 유저 ID 동기화
		record.PlanID = plan.ID
		record.CreatedAt = core.Now()
		record.UpdatedAt = core.Now()

		record, err = s.recordRepo.Save(record)
		if err != nil {
			return apperr.New(202, "save record error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return record, nil
}

func (s *UserService) UpdateRecord(record *domain.Record) (*domain.Record, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		var err error

		record.UpdatedAt = core.Now()
		record, err = s.recordRepo.Update(record)
		if err != nil {
			return apperr.New(202, "update record error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return record, nil
}

func (s *UserService) DeleteRecord(record *domain.Record) (*domain.Record, *apperr.AppError) {
	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}, func(tx *gorm.DB) error {
		var err error

		err = s.recordRepo.DeleteById(record.ID)
		if err != nil {
			return apperr.New(202, "delete record error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}
	return record, nil
}

func (s *UserService) getData(user *domain.User) *apperr.AppError {
	// 폴더
	folders, err := s.folderRepo.FindByUser(user.ID)
	if err != nil {
		return apperr.New(202, "find folders error", err)
	}

	// 플랜
	plans, err := s.planRepo.FindByUser(user.ID)
	if err != nil {
		return apperr.New(202, "find plans error", err)
	}

	// 플랜별 하위 엔티티 채우기
	for _, p := range plans {
		// Goals
		if goals, err := s.goalRepo.FindByPlanID(p.ID); err == nil {
			p.Goals = goals
		} else {
			return apperr.New(202, "find goals error", err)
		}

		// RecordFields (+ OptionItems)
		if fields, err := s.recordFieldRepo.FindByPlanID(p.ID); err == nil {
			for i := range fields {
				opts, _ := s.optionItemRepo.FindByRecordFieldID(fields[i].ID)
				fields[i].Options = buildOptionItemTree(opts)
			}
			p.RecordFields = fields
		} else {
			return apperr.New(202, "find record fields error", err)
		}

		// RepeatRules
		if rules, err := s.repeatRuleRepo.FindByPlanID(p.ID); err == nil {
			p.RepeatRules = rules
		} else {
			return apperr.New(202, "find repeat rules error", err)
		}

		// Statistics
		if stats, err := s.statisticsRepo.FindByPlanID(p.ID); err == nil {
			p.Statistics = stats
		} else {
			return apperr.New(202, "find statistics error", err)
		}

		// Records
		if records, err := s.recordRepo.FindByPlanID(p.ID); err == nil {
			p.Records = records
		} else {
			return apperr.New(202, "find records error", err)
		}

		// ✅ 빈 slice 보정 (nil → [])
		if p.Goals == nil {
			p.Goals = []*domain.Goal{}
		}
		if p.RecordFields == nil {
			p.RecordFields = []*domain.RecordField{}
		}
		if p.RepeatRules == nil {
			p.RepeatRules = []*domain.RepeatRule{}
		}
		if p.Statistics == nil {
			p.Statistics = []*domain.Statistics{}
		}
		if p.Records == nil {
			p.Records = []*domain.Record{}
		}
	}

	// user 에 folder/plan 트리 세팅
	buildUserTree(user, folders, plans)

	return nil
}

// helper: Folder/Plan을 트리 구조에 맞게 붙여주는 함수
func buildUserTree(user *domain.User, folders []*domain.Folder, plans []*domain.Plan) {
	// 빠른 접근을 위해 folder map 생성
	folderMap := make(map[define.SnowflakeID]*domain.Folder)
	for _, f := range folders {
		f.Children = []domain.ChildSystemFile{}
		folderMap[f.ID] = f
	}

	// 먼저 Root 세팅 (ParentID == 0 인 폴더를 Root로)
	for _, f := range folders {
		if f.ParentID == 0 {
			user.Root = f
			break
		}
	}

	// 폴더 트리 세팅
	for _, f := range folders {
		if f.ParentID == 0 {
			continue // Root
		}
		if parent, ok := folderMap[f.ParentID]; ok {
			parent.Children = append(parent.Children, f)
		}
	}

	// 플랜 트리 세팅
	for _, p := range plans {
		if parent, ok := folderMap[p.ParentID]; ok {
			parent.Children = append(parent.Children, p)
		}
	}

	// 모든 폴더의 Children 정렬 (order 기준)
	for _, f := range folderMap {
		sort.SliceStable(f.Children, func(i, j int) bool {
			var oi, oj int

			switch ci := f.Children[i].(type) {
			case *domain.Folder:
				oi = ci.Order
			case *domain.Plan:
				oi = ci.Order
			}

			switch cj := f.Children[j].(type) {
			case *domain.Folder:
				oj = cj.Order
			case *domain.Plan:
				oj = cj.Order
			}

			return oi > oj

		})
	}
}

func buildOptionItemTree(items []*domain.OptionItem) []*domain.OptionItem {
	itemMap := make(map[define.SnowflakeID]*domain.OptionItem)
	var roots []*domain.OptionItem

	// 초기화
	for _, item := range items {
		item.Children = []*domain.OptionItem{}
		itemMap[item.ID] = item
	}

	// 트리 구조 구성
	for _, item := range items {
		if item.ParentID == 0 {
			roots = append(roots, item)
		} else {
			if parent, ok := itemMap[item.ParentID]; ok {
				parent.Children = append(parent.Children, item)
			}
		}
	}

	// 재귀적으로 정렬
	var sortChildren func(opts []*domain.OptionItem)
	sortChildren = func(opts []*domain.OptionItem) {
		sort.SliceStable(opts, func(i, j int) bool {
			return opts[i].Order < opts[j].Order
		})
		for _, opt := range opts {
			if len(opt.Children) > 0 {
				sortChildren(opt.Children)
			}
		}
	}

	// root 기준으로 정렬 시작
	sortChildren(roots)

	return roots
}

func (s *UserService) upsertOptionItemTree(opt *domain.OptionItem, recordFieldID, parentID define.SnowflakeID, optMap map[define.SnowflakeID]bool) error {
	opt.RecordFieldID = recordFieldID
	opt.ParentID = parentID
	opt.UpdatedAt = core.Now()

	if opt.ID == 0 {
		opt.CreatedAt = core.Now()
		newOpt, err := s.optionItemRepo.Save(opt)
		if err != nil {
			return apperr.New(202, "save option item error", err)
		}
		opt.ID = newOpt.ID

	} else {
		if _, err := s.optionItemRepo.Update(opt); err != nil {
			return apperr.New(202, "update option item error", err)
		}
	}

	optMap[opt.ID] = true

	// 자식 재귀 처리
	for i := range opt.Children {
		if err := s.upsertOptionItemTree(opt.Children[i], recordFieldID, opt.ID, optMap); err != nil {
			return err
		}
	}
	return nil
}
