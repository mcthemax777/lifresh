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

	var localDbInfo = DBInfo{"root", "1234", "host.docker.internal:3306", "mysql", "Lifresh"}

	if define.OsType == define.OsTypeWindows {
		localDbInfo = DBInfo{"root", "1234", "127.0.0.1:3306", "mysql", "Lifresh"}
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

func (dh *DBHandlerImpl) InsertAccountAndPlanner(userId string, password string) error {

	tx := dbConn.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	account := models.Account{UserId: userId, Password: password, CreateTime: custom_time.Now()}

	result := tx.Create(&account)

	//account 생성
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//planner 생성
	result = tx.Create(&models.Planner{AccountNo: account.AccountNo, Title: "Planner"})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (dh *DBHandlerImpl) GetAccountByUserId(userId string, password string) (models.Account, error) {

	var account models.Account

	if err := dbConn.Where("userId = ?", userId).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func (dh *DBHandlerImpl) Login(userId string, password string) error {

	var account models.Account

	if err := dbConn.Where("userId = ?", userId).First(&account).Error; err != nil {
		return err
	}

	var planner models.Planner

	if err := dbConn.Where("accountNo = ?", account.AccountNo).First(&planner).Error; err != nil {
		return err
	}

	fmt.Println(planner.PlannerNo)

	return nil
}

func (dh *DBHandlerImpl) GetPlannerByAccountNo(accountNo int) (models.Planner, error) {

	var p models.Planner
	dbConn.Where("accountNo = ?", accountNo).First(&p)

	if p.PlannerNo == 0 {
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

func (dh *DBHandlerImpl) UpdateMainCategoryList(mainCategory *[]models.MainCategory) (int, error) {

	result := dbConn.Save(&mainCategory)

	if result.Error != nil {
		return 0, result.Error
	}

	return 0, nil
}

func (dh *DBHandlerImpl) UpdateSubCategory(subCategory *[]models.SubCategory) error {

	result := dbConn.Save(&subCategory)

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
