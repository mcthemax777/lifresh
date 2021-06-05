package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type dbInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

type LoginReq struct {
	UserId   string
	Password string
}

type LoginRes struct {
	SessionId string
}

var loginQuery = "select password from Account where userId = ? limit 1"

var db = dbInfo{"root", "1234", "localhost:3306", "mysql", "Lifresh"}

func dbQuery(db dbInfo, query string, params string) (password string) {
	dataSource := db.user + ":" + db.pwd + "@tcp(" + db.url + ")/" + db.database
	conn, err := sql.Open(db.engine, dataSource)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	err = conn.QueryRow(query, params).Scan(&password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(password)

	return password
}

func Sum(i int, j int) int {
	return (i + j)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request URI : " + r.RequestURI)

	if r.Method != "GET" {
		fmt.Fprintf(w, "Sorry, only GET methods are supported!!.")
		return
	}

	http.ServeFile(w, r, "./html/login.html")
}

func login(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("request URI : " + r.RequestURI)

	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// len := r.ContentLength
	// body := make([]byte, len)
	// r.Body.Read(body)
	// fmt.Println(string(body))

	decoder := json.NewDecoder(r.Body)
	var requestData LoginReq

	err := decoder.Decode(&requestData)
	if err != nil {
		fmt.Fprintf(w, "decoder err: %s", r.Body)
		return
	}

	userId := requestData.UserId
	password := requestData.Password

	fmt.Printf("id = %s\n", userId)
	fmt.Printf("password = %s\n", password)

	//pararms := [2]string{id, "fefe"}

	passwordInDB := dbQuery(db, loginQuery, userId)

	if password != passwordInDB {
		fmt.Printf("password is not correct : password = %s\n", password)
		return
	}

	//로그인 세션 만들어서 전송
	sessionId := "sessionId1"

	var responseData LoginRes
	responseData.SessionId = sessionId

	encoder := json.NewEncoder(w)
	encoder.Encode(responseData)

	//http.ServeFile(w, r, "./html/build/index.html")
}

func main() {
	//fs := http.FileServer(http.Dir("./html"))
	http.FileServer(http.Dir("./html"))
	http.HandleFunc("/", hello)
	http.HandleFunc("/login", login)

	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatal(err)
	}
}
