package db

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

var localDbInfo = DBInfo{"root", "1234", "localhost:3306", "mysql", "Lifresh"}

type DB interface {
	connect(host string, port int, dbName string, id string, password string) *gorm.DB
	disconnect()
}

type Mysql struct {
	dbConn *gorm.DB
}

func (db Mysql) connect(host string, port int, dbName string, id string, password string) *gorm.DB {
	dsn := localDbInfo.user + ":" + localDbInfo.pwd + "@tcp(" + localDbInfo.url + ")/" + localDbInfo.database + "?charset=utf8"
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil
	}

	db.dbConn = dbConn

	return dbConn
}

func (db Mysql) disconnect() {
	db.dbConn = nil
}

type DBHandler interface {
	connect() bool
	disconnect() bool
	//insertAccount(account *Account) int
}

type DBHandlerImpl struct {
}

func (db DBHandlerImpl) connect() bool {
	return true
}

func (db DBHandlerImpl) disconnect() bool {
	return true
}

// func (db DBHandlerImpl) insertAccount(account *Account) int {
// 	result := db.dbConn.Create(account)

// 	return account.ID
// }
