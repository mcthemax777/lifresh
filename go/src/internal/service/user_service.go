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
)

type UserService struct {
	txm         *txmgr.Manager
	accountRepo *repository.AccountRepo
	userRepo    *repository.UserRepo
	folderRepo  *repository.FolderRepo
	planRepo    *repository.PlanRepo
}

func NewUserService(txm *txmgr.Manager,
	accountRepo *repository.AccountRepo,
	userRepo *repository.UserRepo,
	folderRepo *repository.FolderRepo,
	planRepo *repository.PlanRepo,
) *UserService {
	return &UserService{txm: txm, accountRepo: accountRepo, userRepo: userRepo, folderRepo: folderRepo, planRepo: planRepo}
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
	var folders []*domain.Folder
	var plans []*domain.Plan

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		var err error
		user, err = s.userRepo.FindByAccountID(accountID)
		if err != nil {
			return apperr.New(202, "find user error", err)
		}

		folders, err = s.folderRepo.FindByUser(user.ID)
		if err != nil {
			return apperr.New(202, "find folders error", err)
		}

		plans, err = s.planRepo.FindByUser(user.ID)
		if err != nil {
			return apperr.New(202, "find plans error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	//user 에 folder, plan 구조 세팅
	buildUserTree(user, folders, plans)

	return user, nil
}

func (s *UserService) GetUser(userID define.SnowflakeID) (*domain.User, *apperr.AppError) {

	var user *domain.User
	var folders []*domain.Folder
	var plans []*domain.Plan

	err := s.txm.WithinTx(txmgr.Opts{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false, // 가입도 되어야 해서 false
	}, func(tx *gorm.DB) error {

		var err error
		user, err = s.userRepo.FindByID(userID)
		if err != nil {
			return apperr.New(202, "find user error", err)
		}

		folders, err = s.folderRepo.FindByUser(userID)
		if err != nil {
			return apperr.New(202, "find folders error", err)
		}

		plans, err = s.planRepo.FindByUser(userID)
		if err != nil {
			return apperr.New(202, "find plans error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	//user 에 folder, plan 구조 세팅
	buildUserTree(user, folders, plans)

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

		parentFolder, err := s.folderRepo.FindById(folder.ParentID)
		if err != nil {
			return apperr.New(202, "save folder error", err)
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
		var err error
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
		//TODO. 유저 정보 일치 확인, 부모 폴더가 존재하는지 확인
		var err error
		err = s.folderRepo.Delete(folder)
		if err != nil {
			return apperr.New(202, "delete folder error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err.(*apperr.AppError)
	}

	return folder, nil
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
}
