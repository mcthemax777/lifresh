package db

import (
	"custom_time"
	"errors"
	"fmt"
	"models"
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

// var localDbInfo = DBInfo{"root", "1234", "host.docker.internal:3306", "mysql", "Lifresh"}
var localDbInfo = DBInfo{"root", "1234", "127.0.0.1:3306", "mysql", "Lifresh"}

// type camelNamer struct
// {

// }

// func (cm *camelNamer) ColumnName(table, column string) string {
// 	return column
// }

// func (cm *camelNamer) TableName(table string) string {
// 	return table
// }

// func (cm *camelNamer) JoinTableName(table string) string {
// 	return table
// }

// func (cm *camelNamer) RelationshipFKName(ss schema.Relationship) string {
// 	return ""
// }

// func (cm *camelNamer) CheckerName(table, column string) string {
// 	return column
// }

// func (cm *camelNamer) IndexName(table, column string) string {
// 	return column
// }

func init() {
	//DBHandlerSG := DBHandlerImpl{}
	dsn := localDbInfo.user + ":" + localDbInfo.pwd + "@tcp(" + localDbInfo.url + ")/" + localDbInfo.database + "?charset=utf8&parseTime=true"
	result, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return ;
	}

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
	result = tx.Create(&models.Planner {AccountNo: account.AccountNo, Title: "Planner"});
	
	if  result.Error != nil {
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



func (dh *DBHandlerImpl) GetCFListByPlannerNo(plannerNo int) ([]models.CF, error) {

	var cfList []models.CF
	dbConn.Where("plannerNo = ?", plannerNo).Find(&cfList)

	return cfList, nil
}

func (dh *DBHandlerImpl) GetMCFListByPlannerNo(plannerNo int) ([]models.MCF, error) {

	var mcfList []models.MCF
	dbConn.Where("plannerNo = ?", plannerNo).Find(&mcfList)

	return mcfList, nil
}

func (dh *DBHandlerImpl) GetDCFListByPlannerNo(plannerNo int) ([]models.DCF, error) {

	var dcfList []models.DCF
	dbConn.Where("plannerNo = ?", plannerNo).Find(&dcfList)

	return dcfList, nil
}

func (dh *DBHandlerImpl) GetTaskPlanListByToadyNoAndPlannerNo(todayNo int, plannerNo int) ([]models.TaskPlan, error) {

	var taskPlanList []models.TaskPlan
	dbConn.Where("plannerNo = ? and todayNo = ?", plannerNo, todayNo).Scan(&taskPlanList)

	return taskPlanList, nil
}

func (dh *DBHandlerImpl) GetTaskListByToadyNoAndPlannerNo(todayNo int, plannerNo int) ([]models.Task, error) {

	var taskList []models.Task
	dbConn.Where("plannerNo = ? and todayNo = ?", plannerNo, todayNo).Scan(&taskList)

	return taskList, nil
}

func (dh *DBHandlerImpl) InsertCF(name string, plannerNo int) (int, error) {

	cf := models.CF{Name: name, PlannerNo: plannerNo}

	result := dbConn.Create(&cf)

	//cf 생성
	if result.Error != nil {
		return 0, result.Error
	}

	return cf.CfNo, nil
}

func (dh *DBHandlerImpl) InsertMCF(name string, cfNo int, plannerNo int) (int, error) {

	mcf := models.MCF{Name: name, CfNo: cfNo, PlannerNo: plannerNo}

	result := dbConn.Create(&mcf)

	//cf 생성
	if result.Error != nil {
		return 0, result.Error
	}

	return mcf.McfNo, nil
}

func (dh *DBHandlerImpl) InsertDCF(name string, mcfNo int, dcfType int, priority int, plannerNo int) (int, error) {

	dcf := models.DCF{Name: name, McfNo: mcfNo, DcfType: dcfType, Priority: priority, PlannerNo: plannerNo}

	result := dbConn.Create(&dcf)

	//cf 생성
	if result.Error != nil {
		return 0, result.Error
	}

	return dcf.DcfNo, nil
}

func (dh *DBHandlerImpl) InsertTaskPlan(taskPlan *models.TaskPlan) error {

	result := dbConn.Create(&taskPlan)

	//tp 생성
	if result.Error != nil {
		return result.Error
	}

	return nil
}


func (dh *DBHandlerImpl) InsertTask(task *models.Task) error {

	result := dbConn.Create(&task)

	//tp 생성
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) RemoveTaskPlan(taskPlan *models.TaskPlan) error {

	result := dbConn.Delete(&taskPlan)

	//tp 생성
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dh *DBHandlerImpl) RemoveTask(task *models.Task) error {

	result := dbConn.Delete(&task)

	//tp 생성
	if result.Error != nil {
		return result.Error
	}

	return nil
}

