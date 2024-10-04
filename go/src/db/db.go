package db

import (
	"errors"
	"fmt"
	"lifresh/custom_time"
	"lifresh/define"
	"lifresh/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBHandlerSG DBHandlerImpl

var dbConn *gorm.DB

type DBInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

func init() {

	var localDbInfo = DBInfo{"root", "lifresh", "host.docker.internal:3306", "mysql", "lifresh"}

	if define.OsType == define.OsTypeWindows {
		localDbInfo = DBInfo{"root", "lifresh", "127.0.0.1:3306", "mysql", "lifresh"}
	}

	dsn := localDbInfo.user + ":" + localDbInfo.pwd + "@tcp(" + localDbInfo.url + ")/" + localDbInfo.database + "?charset=utf8&parseTime=true"

	result, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return
	}

	fmt.Println("db init all")

	dbConn = result
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
	InsertAccount(userId string, password string) error
	Login(userId string, password string) error
}

type DBHandlerImpl struct {
	//dbConn *gorm.DB
}

func (dh *DBHandlerImpl) Begin() *gorm.DB {
	return dbConn.Begin()
}

func (dh *DBHandlerImpl) Commit(d *gorm.DB) {
	d.Commit()
}

func (dh *DBHandlerImpl) Rollback(d *gorm.DB) {
	d.Rollback()
}

func (dh *DBHandlerImpl) InsertAccount(socialType int, socialToken string) (models.Account, error) {
	account := models.Account{SocialType: socialType, SocialToken: socialToken, UpdateDate: custom_time.Now(), CreateDate: custom_time.Now()}
	tx := dbConn.Begin()

	if err := tx.Error; err != nil {
		return account, err
	}

	result := tx.Create(&account)

	//account 생성
	if result.Error != nil {
		tx.Rollback()
		return account, result.Error
	}

	//planner 생성
	result = tx.Create(&models.Planner{AccountId: account.AccountId, UpdateDate: custom_time.Now()})

	if result.Error != nil {
		tx.Rollback()
		return account, result.Error
	}

	//money 생성
	result = tx.Create(&models.Money{AccountId: account.AccountId, UpdateDate: custom_time.Now()})

	if result.Error != nil {
		tx.Rollback()
		return account, result.Error
	}

	//diary 생성
	result = tx.Create(&models.Diary{AccountId: account.AccountId, UpdateDate: custom_time.Now()})

	if result.Error != nil {
		tx.Rollback()
		return account, result.Error
	}

	return account, tx.Commit().Error
}

func (dh *DBHandlerImpl) GetAccountBySocialToken(socialType int, socialToken string) (models.Account, error) {

	var account models.Account

	tx := dbConn.Where("social_type = ? AND social_token = ?", socialType, socialToken)

	if tx.Error != nil {
		return account, tx.Error
	}

	account.SocialType = socialType
	account.SocialToken = socialToken

	if err := tx.Take(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

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

func (dh *DBHandlerImpl) GetProfileByAccountId(accountId int) (models.Profile, error) {

	var p models.Profile
	dbConn.Where("account_id = ?", accountId).First(&p)

	if p.ProfileId == 0 {
		fmt.Println("fuck")
		return p, errors.New("not exist profile")
	}

	return p, nil
}

func (dh *DBHandlerImpl) GetPlannerByAccountId(accountId int) (models.Planner, error) {

	var p models.Planner
	dbConn.Where("account_id = ?", accountId).First(&p)

	if p.PlannerId == 0 {
		fmt.Println("fuck")
		return p, errors.New("not exist planner")
	}

	return p, nil
}

func (dh *DBHandlerImpl) GetMoneyByAccountId(accountId int) (models.Money, error) {

	var p models.Money
	dbConn.Where("account_id = ?", accountId).First(&p)

	if p.MoneyId == 0 {
		fmt.Println("fuck")
		return p, errors.New("not exist money")
	}

	return p, nil
}

func (dh *DBHandlerImpl) GetDiaryByAccountId(accountId int) (models.Diary, error) {

	var p models.Diary
	dbConn.Where("account_id = ?", accountId).First(&p)

	if p.DiaryId == 0 {
		fmt.Println("fuck")
		return p, errors.New("not exist Diary")
	}

	return p, nil
}

func (dh *DBHandlerImpl) GetPlanCategoryListByPlannerId(plannerId int) ([]models.PlanCategory, error) {

	var list []models.PlanCategory
	dbConn.Where("planner_id = ?", plannerId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetPlanListByPlannerId(plannerId int) ([]models.Plan, error) {

	var list []models.Plan
	dbConn.Where("planner_id = ?", plannerId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetPlanHistoryListByPlannerId(plannerId int) ([]models.PlanHistory, error) {

	var list []models.PlanHistory
	dbConn.Where("planner_id = ?", plannerId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyCategoryListByMoneyId(moneyId int) ([]models.MoneyCategory, error) {

	var list []models.MoneyCategory
	dbConn.Where("money_id = ?", moneyId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyHistoryListByMoneyId(moneyId int) ([]models.MoneyHistory, error) {

	var list []models.MoneyHistory
	dbConn.Where("money_id = ?", moneyId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetDiaryCategoryListByMoneyId(diaryId int) ([]models.DiaryCategory, error) {

	var list []models.DiaryCategory
	dbConn.Where("diary_id = ?", diaryId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetDiaryHistoryListByMoneyId(diaryId int) ([]models.DiaryHistory, error) {

	var list []models.DiaryHistory
	dbConn.Where("diary_id = ?", diaryId).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) InsertDiaryCategoryList(diaryCategoryList *[]models.DiaryCategory) error {

	result := dbConn.Create(&diaryCategoryList)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateDiaryCategoryList(diaryCategoryList *[]models.DiaryCategory) error {

	result := dbConn.Save(&diaryCategoryList)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertDiaryHistoryList(diaryHistoryList *[]models.DiaryHistory) error {

	result := dbConn.Create(&diaryHistoryList)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateDiaryHistoryList(diaryHistoryList *[]models.DiaryHistory) error {

	result := dbConn.Save(&diaryHistoryList)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) DeleteDiaryCategoryList(diaryId int, diaryCategoryIdList []int) error {

	result := dbConn.Where("diary_id = ? AND diary_category_id IN ?", diaryId, diaryCategoryIdList).Delete(&models.DiaryCategory{})

	return result.Error
}

func (dh *DBHandlerImpl) DeleteDiaryHistoryList(diaryId int, diaryHistoryIdList []int) error {

	result := dbConn.Where("diary_id = ? AND diary_history_id IN ?", diaryId, diaryHistoryIdList).Delete(&models.DiaryHistory{})

	return result.Error
}

func (dh *DBHandlerImpl) GetPlannerByAccountNo(accountNo int) (models.Planner, error) {

	var p models.Planner
	dbConn.Where("accountNo = ?", accountNo).First(&p)

	if p.PlannerId == 0 {
		fmt.Println("fuck")
		return p, errors.New("not exist planner")
	}

	return p, nil
}

func (dh *DBHandlerImpl) GetTodayByPlannerNoAndTime(plannerNo int, myTime time.Time) (models.Today, error) {

	//받은 시간의 이전 00시, 다음 00시 구하기
	minToday := time.Date(myTime.Year(), myTime.Month(), myTime.Day(), 0, 0, 0, 0, time.UTC)

	var today models.Today
	dbConn.Where("plannerNo = ? and todayDate >= ?", plannerNo, minToday).First(&today)

	//없다면 만들어주기
	if today.TodayNo == 0 {
		today.PlannerNo = plannerNo
		//today.TodayDate = CustomTime(minToday)
		today.Diary = "diary"

		result := dbConn.Create(&today)

		if result.Error != nil {
			return today, result.Error
		}
	}

	return today, nil
}

func (dh *DBHandlerImpl) GetTodayByTodayNo(todayNo int, plannerNo int) (models.Today, error) {

	var today models.Today
	dbConn.Where("plannerNo = ? and todayNo = ?", plannerNo, todayNo).First(&today)

	return today, nil
}

func (dh *DBHandlerImpl) GetTodayListByPlannerNo(plannerNo int) ([]models.Today, error) {

	var list []models.Today
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyManagerListByPlannerNoAndMoneyManagerNoList(plannerNo int, moneyManagerNoList []int) ([]models.MoneyManager, error) {

	var list []models.MoneyManager
	dbConn.Where("plannerNo = ? and moneyManagerNo in (?)", plannerNo, moneyManagerNoList).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyManagerListByPlannerNo(plannerNo int) ([]models.MoneyManager, error) {

	var list []models.MoneyManager
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMainCategoryListByPlannerNo(plannerNo int) ([]models.MainCategory, error) {

	var list []models.MainCategory
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMainCategoryListByPlannerNoAndCategoryTypeList(plannerNo int, categoryTypeList []int) ([]models.MainCategory, error) {

	var list []models.MainCategory
	dbConn.Where("plannerNo = ? and categoryType in (?)", plannerNo, categoryTypeList).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetSubCategoryListByPlannerNoAndMainCategoryNoList(plannerNo int, mainCategoryNoList []int) ([]models.SubCategory, error) {

	var list []models.SubCategory
	dbConn.Where("plannerNo = ? and mainCategoryNo in (?)", plannerNo, mainCategoryNoList).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetSubCategoryListByPlannerNo(plannerNo int) ([]models.SubCategory, error) {

	var list []models.SubCategory
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetScheduleTaskListByPlannerNo(plannerNo int) ([]models.ScheduleTask, error) {

	var list []models.ScheduleTask
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetToDoTaskListByPlannerNo(plannerNo int) ([]models.ToDoTask, error) {

	var list []models.ToDoTask
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyTaskListByPlannerNo(plannerNo int) ([]models.MoneyTask, error) {

	var list []models.MoneyTask
	dbConn.Where("plannerNo = ?", plannerNo).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) GetMoneyTaskListByPlannerNoAndMoneyTaskNoList(plannerNo int, moneyTaskNoList []int) ([]models.MoneyTask, error) {

	var list []models.MoneyTask
	dbConn.Where("plannerNo = ? and moneyTaskNo in (?)", plannerNo, moneyTaskNoList).Find(&list)

	return list, nil
}

func (dh *DBHandlerImpl) InsertMainCategory(mainCategory *[]models.MainCategory) error {

	result := dbConn.Create(&mainCategory)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertSubCategory(subCategory *[]models.SubCategory) error {

	result := dbConn.Create(&subCategory)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertMoneyManager(moneyManager *[]models.MoneyManager) error {

	result := dbConn.Create(&moneyManager)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertScheduleTask(scheduleTask *[]models.ScheduleTask) error {

	result := dbConn.Create(&scheduleTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertToDoTask(toDoTask *[]models.ToDoTask) error {

	result := dbConn.Create(&toDoTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) InsertMoneyTask(moneyTask *[]models.MoneyTask) error {

	result := dbConn.Create(&moneyTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateMainCategory(mainCategory *[]models.MainCategory) error {

	result := dbConn.Save(&mainCategory)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateSubCategory(subCategory *[]models.SubCategory) error {

	result := dbConn.Save(&subCategory)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateMoneyManager(moneyManager *[]models.MoneyManager) error {

	result := dbConn.Save(&moneyManager)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateScheduleTask(scheduleTask *[]models.ScheduleTask) error {

	result := dbConn.Save(&scheduleTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateToDoTask(toDoTask *[]models.ToDoTask) error {

	result := dbConn.Save(&toDoTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) UpdateMoneyTask(moneyTask *[]models.MoneyTask) error {

	result := dbConn.Save(&moneyTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) DeleteMainCategoryList(plannerNo int, mainCategoryNoList []int) (error, error) {

	result := dbConn.Where("plannerNo = ? AND mainCategoryNo IN ?", plannerNo, mainCategoryNoList).Delete(&models.MainCategory{})

	//tp 생성
	if result.Error != nil {
		return result.Error, nil
	}

	return nil, nil
}

func (dh *DBHandlerImpl) DeleteSubCategoryList(plannerNo int, subCategoryNoList []int) (error, error) {

	result := dbConn.Where("plannerNo = ? AND subCategoryNo IN ?", plannerNo, subCategoryNoList).Delete(&models.SubCategory{})

	//tp 생성
	if result.Error != nil {
		return result.Error, nil
	}

	return nil, nil
}

func (dh *DBHandlerImpl) DeleteScheduleTaskList(plannerNo int, scheduleTaskNoList []int) (error, error) {

	result := dbConn.Where("plannerNo = ? AND scheduleTaskNo IN ?", plannerNo, scheduleTaskNoList).Delete(&models.ScheduleTask{})

	//tp 생성
	if result.Error != nil {
		return result.Error, nil
	}

	return nil, nil
}

func (dh *DBHandlerImpl) DeleteToDoTaskList(plannerNo int, toDoTaskNoList []int) (error, error) {

	result := dbConn.Where("plannerNo = ? AND toDoTaskNo IN ?", plannerNo, toDoTaskNoList).Delete(&models.ToDoTask{})

	//tp 생성
	if result.Error != nil {
		return result.Error, nil
	}

	return nil, nil
}

func (dh *DBHandlerImpl) DeleteMoneyTaskList(plannerNo int, moneyTaskNo []int) (error, error) {

	result := dbConn.Where("plannerNo = ? AND moneyTaskNo IN ?", plannerNo, moneyTaskNo).Delete(&models.MoneyTask{})

	//tp 생성
	if result.Error != nil {
		return result.Error, nil
	}

	return nil, nil
}
