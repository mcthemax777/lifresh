package db

import (
	"database/sql"
	"fmt"
	"github.com/sony/sonyflake"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"lifresh/define"
	"lifresh/internal/repository/model"
	"regexp"
	"strings"
)

var sf *sonyflake.Sonyflake

func initIDGen() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		// MachineID 함수에서 리전/노드 ID 제공 가능
	})
}

func NextID() define.SnowflakeID {
	id, _ := sf.NextID()
	return define.SnowflakeID(id)
}

var DBHandlerSG DBHandlerImpl

var dbConn *gorm.DB

type DBInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

// 커스텀 Namer: 기본 NamingStrategy를 임베드하고 ColumnName만 오버라이드
type CustomNamingOption struct {
	schema.NamingStrategy
	acronyms []string
	reAcr    *regexp.Regexp // ([a-z0-9])(URL|ID|API|...) 경계
	reCamel  *regexp.Regexp // camelCase -> snake_case 경계
}

func NewNamingOption(prefix string, singular bool, acronyms []string) CustomNamingOption {
	// 소문자/숫자 뒤에 약어가 오면 언더스코어를 넣기 위한 정규식
	// 예: photoURL -> photo_URL (이후 snake 처리)
	pat := fmt.Sprintf(`([a-z0-9])(%s)\b`, strings.Join(acronyms, "|"))
	return CustomNamingOption{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		acronyms: acronyms,
		reAcr:    regexp.MustCompile(pat),
		reCamel:  regexp.MustCompile(`([a-z0-9])([A-Z])`),
	}
}

// ColumnName: 약어 경계에 언더스코어 삽입 → 카멜 경계 underscoring → 전부 소문자
func (n CustomNamingOption) ColumnName(_ string, column string) string {
	// 1) 약어를 모두 대문자로 정규화 (Url→URL, Id→ID 등)
	for _, ac := range n.acronyms {
		title := strings.Title(strings.ToLower(ac)) // Url, Id, Api ...
		column = strings.ReplaceAll(column, title, ac)
		column = strings.ReplaceAll(column, strings.ToLower(ac), ac)
	}
	// 2) 컬럼 전체가 약어 하나인 경우: 바로 소문자 반환 (URL -> url)
	up := strings.ToUpper(column)
	for _, ac := range n.acronyms {
		if up == ac {
			return strings.ToLower(ac)
		}
	}
	// 3) 약어 앞에 언더스코어 삽입 (예: photoURL -> photo_URL)
	column = n.reAcr.ReplaceAllString(column, "${1}_$2")
	// 4) 일반 카멜 경계도 언더스코어 삽입 (예: apiKey -> api_Key)
	column = n.reCamel.ReplaceAllString(column, "${1}_${2}")
	// 5) 전부 소문자
	return strings.ToLower(column)
}

func InitDB() *gorm.DB {
	initIDGen()

	var localDbInfo = DBInfo{"root", "lifresh", "host.docker.internal:3306", "mysql", "lifresh"}

	if define.OsType == define.OsTypeWindows || define.OsType == define.OsTypeMac {
		localDbInfo = DBInfo{"root", "lifresh", "127.0.0.1:3306", "mysql", "lifresh"}
	}

	dsn := localDbInfo.user + ":" + localDbInfo.pwd + "@tcp(" + localDbInfo.url + ")/" + localDbInfo.database + "?charset=utf8&parseTime=true"

	namingOption := NewNamingOption(
		"",   // TablePrefix (예: "cm2_")
		true, // SingularTable
		[]string{"URL", "ID", "API", "HTML", "JSON", "IP"}, // 필요한 약어 추가
	)
	result, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// AutoMigrate 시 외래키 제약 생성 비활성화
		DisableForeignKeyConstraintWhenMigrating: true,

		NamingStrategy: namingOption,
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("db init all")

	dbConn = result

	err = model.AutoMigrate(dbConn)

	if err != nil {
		panic("failed to migrate database")
	}

	return dbConn
}

// type DB interface {
// 	connect()
// 	disconnect()
// 	create() int
// }

// type ctMysql struct {
// 	dbConn *gorm.DB
// }

// func (db ctMysql) connect() {
// 	dsn := localDbInfo.user + ":" + localDbInfo.pwd + "@tcp(" + localDbInfo.url + ")/" + localDbInfo.database + "?charset=utf8"
// 	dbconn1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		//fmt.Println("mysql connect error" + err.Error)
// 	}

// 	db.dbConn = dbconn1
// }

// func (db ctMysql) disconnect() {
// 	db.dbConn = nil
// }

type DBHandler interface {
	InsertAccount(socialType int, socialToken string) (model.Account, error)
	Login(userId string, password string) error
}

type DBHandlerImpl struct {
	//dbConn *gorm.DB
}

func (dh *DBHandlerImpl) Begin() *gorm.DB {
	return dbConn.Begin(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
}

func (dh *DBHandlerImpl) Commit(d *gorm.DB) {
	d.Commit()
}

func (dh *DBHandlerImpl) Rollback(d *gorm.DB) {
	d.Rollback()
}

//
//func (dh *DBHandlerImpl) InsertAccount(socialType int, uid string, name string, email string) (models.Account, error) {
//	account := models.Account{SocialType: models.SocialType(socialType), ProviderUID: &uid} // SocialToken: socialToken, UpdateDate: custom_time.Now(), CreateDate: custom_time.Now()}
//	tx := dbConn.Begin()
//
//	if err := tx.Error; err != nil {
//		return account, err
//	}
//
//	result := tx.Create(&account)
//
//	//account 생성
//	if result.Error != nil {
//		tx.Rollback()
//		return account, result.Error
//	}
//
//	user := models.User{ID: account.ID, AccountID: account.ID, Nickname: uid, Bio: "", UpdatedAt: time.Now()}
//	//user 생성
//	result = tx.Create(&user)
//
//	if result.Error != nil {
//		tx.Rollback()
//		return account, result.Error
//	}
//
//	//root 생성
//	root := models.Folder{ID: account.ID, UserID: user.ID, Name: "root", UpdatedAt: time.Now()}
//	result = tx.Create(&root)
//
//	if result.Error != nil {
//		tx.Rollback()
//		return account, result.Error
//	}
//
//	//account.User = &user
//	//user.Root = &root
//
//	return account, tx.Commit().Error
//}
//
//func (dh *DBHandlerImpl) GetAccountByUID(uid string) (models.Account, error) {
//
//	var account models.Account
//	if err := dbConn.Where("provider_uid = ?", uid).First(&account).Error; err != nil {
//		return account, err
//	}
//
//	return account, nil
//
//	//var account models.Account
//	//var user models.User
//	//var root models.Folder
//	//
//	//// 같은 스냅샷에서 account -> user -> root를 읽기 위한 읽기 전용 트랜잭션
//	//tx := dbConn.Begin(&sql.TxOptions{
//	//	Isolation: sql.LevelRepeatableRead,
//	//	ReadOnly:  true,
//	//})
//	//if tx.Error != nil {
//	//	return account, tx.Error
//	//}
//	//// 롤백 안전장치
//	//defer func() { _ = tx.Rollback() }()
//	//
//	//// 1) account
//	//if err := tx.Where("provider_uid = ?", uid).Take(&account).Error; err != nil {
//	//	if errors.Is(err, gorm.ErrRecordNotFound) {
//	//		return account, err
//	//	}
//	//	return account, err
//	//}
//	//
//	//// 2) user
//	//if err := tx.First(&user, "account_id = ?", account.ID).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//	//	return account, err
//	//}
//	//
//	//// 3) root
//	//_ = tx.Where("user_id = ? AND parent_id IS NULL", user.ID).Take(&root).Error
//	//
//	//if err := tx.Commit().Error; err != nil {
//	//	return account, err
//	//}
//	//
//	//account.User = &user
//	//user.Root = &root
//	//
//	//return account, nil
//}
//
//func (dh *DBHandlerImpl) GetUserByUID(userID define.SnowflakeID) (*models.User, []models.Folder, []models.Plan, error) {
//
//	var user models.User
//
//	// 같은 스냅샷에서 account -> user -> root를 읽기 위한 읽기 전용 트랜잭션
//	tx := dbConn.Begin(&sql.TxOptions{
//		Isolation: sql.LevelRepeatableRead,
//		ReadOnly:  true,
//	})
//	if tx.Error != nil {
//		return nil, nil, nil, tx.Error
//	}
//	// 롤백 안전장치
//	defer tx.Rollback()
//
//	// get user
//	if err := tx.First(&user, "id = ?", userID).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, nil, nil, err
//	}
//
//	var folderList []models.Folder
//	var planList []models.Plan
//
//	// get folder
//	if err := tx.Where("user_id = ?", userID).Find(&folderList).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, nil, nil, err
//	}
//
//	_ = tx.Where("user_id = ?", userID).Find(&planList).Error
//
//	if err := tx.Commit().Error; err != nil {
//		return nil, nil, nil, err
//	}
//
//	return &user, folderList, planList, nil
//}
//
//func (dh *DBHandlerImpl) GetFolderByID(folderID define.SnowflakeID) (models.Folder, error) {
//
//	var item models.Folder
//
//	// 같은 스냅샷에서 account -> user -> root를 읽기 위한 읽기 전용 트랜잭션
//	tx := dbConn.Begin(&sql.TxOptions{
//		Isolation: sql.LevelRepeatableRead,
//		ReadOnly:  true,
//	})
//	if tx.Error != nil {
//		return item, tx.Error
//	}
//	// 롤백 안전장치
//	defer func() { _ = tx.Rollback() }()
//
//	// 3) list
//	_ = tx.Where("id = ?", folderID).Find(&item).Error
//
//	if err := tx.Commit().Error; err != nil {
//		return item, err
//	}
//
//	userID := item.UserID
//
//	var folderList []models.Folder
//	var planList []models.Plan
//
//	// 3) list
//	_ = tx.Where("user_id = ?", userID).Find(&folderList).Error
//
//	if err := tx.Commit().Error; err != nil {
//		return item, err
//	}
//
//	folderMap := make(map[define.SnowflakeID]*models.Folder)
//	for _, folder := range folderList {
//		folderMap[folder.ID] = &folder
//	}
//
//	for _, folder := range folderList {
//		if folder.ParentID == nil {
//
//		} else {
//			parent := folderMap[*folder.ParentID]
//			parent.ChildrenFolders = append(parent.ChildrenFolders, folder)
//		}
//	}
//
//	_ = tx.Where("parent_id IN ?", userID).Find(&planList).Error
//
//	if err := tx.Commit().Error; err != nil {
//		return item, err
//	}
//
//	return item, nil
//}

//func (dh *DBHandlerImpl) GetAccountByUserId(userId string, password string) (models.Account, error) {
//
//	var account models.Account
//
//	if err := dbConn.Where("userId = ?", userId).First(&account).Error; err != nil {
//		return account, err
//	}
//
//	return account, nil
//}

//func (dh *DBHandlerImpl) Login(userId string, password string) error {
//
//	var account models.Account
//
//	if err := dbConn.Where("userId = ?", userId).First(&account).Error; err != nil {
//		return err
//	}
//
//	var planner models.Planner
//
//	if err := dbConn.Where("accountNo = ?", account.AccountNo).First(&planner).Error; err != nil {
//		return err
//	}
//
//	fmt.Println(planner.PlannerId)
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertUserAndRootFolder(accountId int, nickname string, rootFolderName string) (models.User, error) {
//	user := models.User{AccountId: accountId, Nickname: nickname, ProfileImageUrl: "", UpdateDate: custom_time.Now()}
//	tx := dbConn.Begin()
//
//	if err := tx.Error; err != nil {
//		return user, err
//	}
//
//	result := tx.Create(&user)
//
//	//user 생성
//	if result.Error != nil {
//		tx.Rollback()
//		return user, result.Error
//	}
//
//	//root plan 생성
//	result = tx.Create(&models.Plan{
//		UserId:         user.Id,
//		Name:           rootFolderName,
//		Description:    "",
//		RGBColor:       1000000000,
//		ParentId:       0,
//		Type:           0,
//		PermissionType: 0,
//		StartDate:      custom_time.Now(),
//		EndDate:        custom_time.Now(),
//		Sort:           0,
//		UpdateDate:     custom_time.Now(),
//	})
//
//	if result.Error != nil {
//		tx.Rollback()
//		return user, result.Error
//	}
//
//	return user, tx.Commit().Error
//}
//
//func (dh *DBHandlerImpl) GetUserByAccountId(accountId int) (models.User, error) {
//
//	var p models.User
//	dbConn.Where("account_id = ?", accountId).First(&p)
//
//	if p.Id == 0 {
//		return p, errors.New("not exist user")
//	}
//
//	return p, nil
//}
//
///// select
//
//func (dh *DBHandlerImpl) GetPlanListByUserId(userId int) ([]models.Plan, error) {
//
//	var list []models.Plan
//	dbConn.Where("user_id = ?", userId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetPlanRecordListByUserId(userId int) ([]models.PlanRecord, error) {
//
//	var list []models.PlanRecord
//	dbConn.Where("user_id = ?", userId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetPlanRecordOperatorListByUserId(userId int) ([]models.PlanRecordOperator, error) {
//
//	var list []models.PlanRecordOperator
//	dbConn.Where("user_id = ?", userId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetPlanGoalListByUserId(userId int) ([]models.PlanGoal, error) {
//
//	var list []models.PlanGoal
//	dbConn.Where("user_id = ?", userId).Find(&list)
//
//	return list, nil
//}
//
///// insert
//
//func (dh *DBHandlerImpl) InsertPlanList(planList *[]models.Plan) error {
//
//	result := dbConn.Create(&planList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertPlanRecordList(planRecordList *[]models.PlanRecord) error {
//
//	result := dbConn.Create(&planRecordList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertPlanRecordOperatorList(planRecordOperatorList *[]models.PlanRecordOperator) error {
//
//	result := dbConn.Create(&planRecordOperatorList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertPlanGoalList(planGoalList *[]models.PlanGoal) error {
//
//	result := dbConn.Create(&planGoalList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
///// update
//
//func (dh *DBHandlerImpl) UpdatePlanList(planList *[]models.Plan) error {
//
//	result := dbConn.Save(&planList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdatePlanRecordList(planRecordList *[]models.PlanRecord) error {
//
//	result := dbConn.Save(&planRecordList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdatePlanRecordOperatorList(planRecordOperatorList *[]models.PlanRecordOperator) error {
//
//	result := dbConn.Save(&planRecordOperatorList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdatePlanGoalList(planGoalList *[]models.PlanGoal) error {
//
//	result := dbConn.Save(&planGoalList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
///// delete
//
//func (dh *DBHandlerImpl) DeletePlanList(userId int, planIdList []int) error {
//
//	result := dbConn.Where("user_id = ? AND id IN ?", userId, planIdList).Delete(&models.Plan{})
//
//	return result.Error
//}
//
//func (dh *DBHandlerImpl) DeletePlanRecordList(userId int, planRecordIdList []int) error {
//
//	result := dbConn.Where("user_id = ? AND id IN ?", userId, planRecordIdList).Delete(&models.PlanRecord{})
//
//	return result.Error
//}
//
//func (dh *DBHandlerImpl) DeletePlanRecordOperatorList(userId int, planRecordOperatorIdList []int) error {
//
//	result := dbConn.Where("user_id = ? AND id IN ?", userId, planRecordOperatorIdList).Delete(&models.PlanRecordOperator{})
//
//	return result.Error
//}
//
//func (dh *DBHandlerImpl) DeletePlanGoalList(userId int, planGoalIdList []int) error {
//
//	result := dbConn.Where("user_id = ? AND id IN ?", userId, planGoalIdList).Delete(&models.PlanGoal{})
//
//	return result.Error
//}
//
///// 안씀
//
//func (dh *DBHandlerImpl) GetPlanHistoryListByUserId(plannerId int) ([]models.PlanHistory, error) {
//
//	var list []models.PlanHistory
//	dbConn.Where("planner_id = ?", plannerId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyCategoryListByMoneyId(moneyId int) ([]models.MoneyCategory, error) {
//
//	var list []models.MoneyCategory
//	dbConn.Where("money_id = ?", moneyId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyHistoryListByMoneyId(moneyId int) ([]models.MoneyHistory, error) {
//
//	var list []models.MoneyHistory
//	dbConn.Where("money_id = ?", moneyId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetDiaryCategoryListByMoneyId(diaryId int) ([]models.DiaryCategory, error) {
//
//	var list []models.DiaryCategory
//	dbConn.Where("diary_id = ?", diaryId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetDiaryHistoryListByMoneyId(diaryId int) ([]models.DiaryHistory, error) {
//
//	var list []models.DiaryHistory
//	dbConn.Where("diary_id = ?", diaryId).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) InsertDiaryCategoryList(diaryCategoryList *[]models.DiaryCategory) error {
//
//	result := dbConn.Create(&diaryCategoryList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateDiaryCategoryList(diaryCategoryList *[]models.DiaryCategory) error {
//
//	result := dbConn.Save(&diaryCategoryList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertDiaryHistoryList(diaryHistoryList *[]models.DiaryHistory) error {
//
//	result := dbConn.Create(&diaryHistoryList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateDiaryHistoryList(diaryHistoryList *[]models.DiaryHistory) error {
//
//	result := dbConn.Save(&diaryHistoryList)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) DeleteDiaryCategoryList(diaryId int, diaryCategoryIdList []int) error {
//
//	result := dbConn.Where("diary_id = ? AND diary_category_id IN ?", diaryId, diaryCategoryIdList).Delete(&models.DiaryCategory{})
//
//	return result.Error
//}
//
//func (dh *DBHandlerImpl) DeleteDiaryHistoryList(diaryId int, diaryHistoryIdList []int) error {
//
//	result := dbConn.Where("diary_id = ? AND diary_history_id IN ?", diaryId, diaryHistoryIdList).Delete(&models.DiaryHistory{})
//
//	return result.Error
//}
//
//func (dh *DBHandlerImpl) GetPlannerByAccountNo(accountNo int) (models.Planner, error) {
//
//	var p models.Planner
//	dbConn.Where("accountNo = ?", accountNo).First(&p)
//
//	if p.PlannerId == 0 {
//		fmt.Println("fuck")
//		return p, errors.New("not exist planner")
//	}
//
//	return p, nil
//}
//
//func (dh *DBHandlerImpl) GetTodayByPlannerNoAndTime(plannerNo int, myTime time.Time) (models.Today, error) {
//
//	//받은 시간의 이전 00시, 다음 00시 구하기
//	minToday := time.Date(myTime.Year(), myTime.Month(), myTime.Day(), 0, 0, 0, 0, time.UTC)
//
//	var today models.Today
//	dbConn.Where("plannerNo = ? and todayDate >= ?", plannerNo, minToday).First(&today)
//
//	//없다면 만들어주기
//	if today.TodayNo == 0 {
//		today.PlannerNo = plannerNo
//		//today.TodayDate = CustomTime(minToday)
//		today.Diary = "diary"
//
//		result := dbConn.Create(&today)
//
//		if result.Error != nil {
//			return today, result.Error
//		}
//	}
//
//	return today, nil
//}
//
//func (dh *DBHandlerImpl) GetTodayByTodayNo(todayNo int, plannerNo int) (models.Today, error) {
//
//	var today models.Today
//	dbConn.Where("plannerNo = ? and todayNo = ?", plannerNo, todayNo).First(&today)
//
//	return today, nil
//}
//
//func (dh *DBHandlerImpl) GetTodayListByPlannerNo(plannerNo int) ([]models.Today, error) {
//
//	var list []models.Today
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyManagerListByPlannerNoAndMoneyManagerNoList(plannerNo int, moneyManagerNoList []int) ([]models.MoneyManager, error) {
//
//	var list []models.MoneyManager
//	dbConn.Where("plannerNo = ? and moneyManagerNo in (?)", plannerNo, moneyManagerNoList).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyManagerListByPlannerNo(plannerNo int) ([]models.MoneyManager, error) {
//
//	var list []models.MoneyManager
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMainCategoryListByPlannerNo(plannerNo int) ([]models.MainCategory, error) {
//
//	var list []models.MainCategory
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMainCategoryListByPlannerNoAndCategoryTypeList(plannerNo int, categoryTypeList []int) ([]models.MainCategory, error) {
//
//	var list []models.MainCategory
//	dbConn.Where("plannerNo = ? and categoryType in (?)", plannerNo, categoryTypeList).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetSubCategoryListByPlannerNoAndMainCategoryNoList(plannerNo int, mainCategoryNoList []int) ([]models.SubCategory, error) {
//
//	var list []models.SubCategory
//	dbConn.Where("plannerNo = ? and mainCategoryNo in (?)", plannerNo, mainCategoryNoList).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetSubCategoryListByPlannerNo(plannerNo int) ([]models.SubCategory, error) {
//
//	var list []models.SubCategory
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetScheduleTaskListByPlannerNo(plannerNo int) ([]models.ScheduleTask, error) {
//
//	var list []models.ScheduleTask
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetToDoTaskListByPlannerNo(plannerNo int) ([]models.ToDoTask, error) {
//
//	var list []models.ToDoTask
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyTaskListByPlannerNo(plannerNo int) ([]models.MoneyTask, error) {
//
//	var list []models.MoneyTask
//	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) GetMoneyTaskListByPlannerNoAndMoneyTaskNoList(plannerNo int, moneyTaskNoList []int) ([]models.MoneyTask, error) {
//
//	var list []models.MoneyTask
//	dbConn.Where("plannerNo = ? and moneyTaskNo in (?)", plannerNo, moneyTaskNoList).Find(&list)
//
//	return list, nil
//}
//
//func (dh *DBHandlerImpl) InsertMainCategory(mainCategory *[]models.MainCategory) error {
//
//	result := dbConn.Create(&mainCategory)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertSubCategory(subCategory *[]models.SubCategory) error {
//
//	result := dbConn.Create(&subCategory)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertMoneyManager(moneyManager *[]models.MoneyManager) error {
//
//	result := dbConn.Create(&moneyManager)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertScheduleTask(scheduleTask *[]models.ScheduleTask) error {
//
//	result := dbConn.Create(&scheduleTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertToDoTask(toDoTask *[]models.ToDoTask) error {
//
//	result := dbConn.Create(&toDoTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) InsertMoneyTask(moneyTask *[]models.MoneyTask) error {
//
//	result := dbConn.Create(&moneyTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateMainCategory(mainCategory *[]models.MainCategory) error {
//
//	result := dbConn.Save(&mainCategory)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateSubCategory(subCategory *[]models.SubCategory) error {
//
//	result := dbConn.Save(&subCategory)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateMoneyManager(moneyManager *[]models.MoneyManager) error {
//
//	result := dbConn.Save(&moneyManager)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateScheduleTask(scheduleTask *[]models.ScheduleTask) error {
//
//	result := dbConn.Save(&scheduleTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateToDoTask(toDoTask *[]models.ToDoTask) error {
//
//	result := dbConn.Save(&toDoTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) UpdateMoneyTask(moneyTask *[]models.MoneyTask) error {
//
//	result := dbConn.Save(&moneyTask)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	return nil
//}
//
//func (dh *DBHandlerImpl) DeleteMainCategoryList(plannerNo int, mainCategoryNoList []int) (error, error) {
//
//	result := dbConn.Where("plannerNo = ? AND mainCategoryNo IN ?", plannerNo, mainCategoryNoList).Delete(&models.MainCategory{})
//
//	//tp 생성
//	if result.Error != nil {
//		return result.Error, nil
//	}
//
//	return nil, nil
//}
//
//func (dh *DBHandlerImpl) DeleteSubCategoryList(plannerNo int, subCategoryNoList []int) (error, error) {
//
//	result := dbConn.Where("plannerNo = ? AND subCategoryNo IN ?", plannerNo, subCategoryNoList).Delete(&models.SubCategory{})
//
//	//tp 생성
//	if result.Error != nil {
//		return result.Error, nil
//	}
//
//	return nil, nil
//}
//
//func (dh *DBHandlerImpl) DeleteScheduleTaskList(plannerNo int, scheduleTaskNoList []int) (error, error) {
//
//	result := dbConn.Where("plannerNo = ? AND scheduleTaskNo IN ?", plannerNo, scheduleTaskNoList).Delete(&models.ScheduleTask{})
//
//	//tp 생성
//	if result.Error != nil {
//		return result.Error, nil
//	}
//
//	return nil, nil
//}
//
//func (dh *DBHandlerImpl) DeleteToDoTaskList(plannerNo int, toDoTaskNoList []int) (error, error) {
//
//	result := dbConn.Where("plannerNo = ? AND toDoTaskNo IN ?", plannerNo, toDoTaskNoList).Delete(&models.ToDoTask{})
//
//	//tp 생성
//	if result.Error != nil {
//		return result.Error, nil
//	}
//
//	return nil, nil
//}
//
//func (dh *DBHandlerImpl) DeleteMoneyTaskList(plannerNo int, moneyTaskNo []int) (error, error) {
//
//	result := dbConn.Where("plannerNo = ? AND moneyTaskNo IN ?", plannerNo, moneyTaskNo).Delete(&models.MoneyTask{})
//
//	//tp 생성
//	if result.Error != nil {
//		return result.Error, nil
//	}
//
//	return nil, nil
//}
