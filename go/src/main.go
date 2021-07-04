package main

import (
	"api"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
}

func hello(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H {})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/hello", hello)
	router.POST("/api/:name", api.ApiCall)

	router.Run(":3000")
}

// type dbInfo struct {
// 	user     string
// 	pwd      string
// 	url      string
// 	engine   string
// 	database string
// }

// var loginQuery = "select password from Account where userId = ? limit 1"

// var info = dbInfo{"root", "1234", "localhost:3306", "mysql", "Lifresh"}

// var dbHandler db.DBHandlerImpl

// func dbQuery(db dbInfo, query string, params string) (password string) {
// 	dataSource := db.user + ":" + db.pwd + "@tcp(" + db.url + ")/" + db.database
// 	conn, err := sql.Open(db.engine, dataSource)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer conn.Close()

// 	err = conn.QueryRow(query, params).Scan(&password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(password)

// 	return password
// }

// func Sum(i int, j int) int {
// 	return (i + j)
// }


// func signUp(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("request URI : " + r.RequestURI)

// 	if r.Method != "POST" {
// 		fmt.Fprintf(w, "Sorry, only POST methods are supported~~~")
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	var requestData SignUpReq

// 	err := decoder.Decode(&requestData)
// 	if err != nil {
// 		fmt.Fprintf(w, "decoder err: %s", r.Body)
// 		return
// 	}

// 	userId := requestData.UserId
// 	password := requestData.Password

// 	var account models.Account
// 	account.UserId = userId
// 	account.Password = password

// 	accountId := dbHandler.InsertAccount(&account)

// 	// if err != nil {
// 	// 	fmt.Fprintf(w, "insertAccount err: %v", err)
// 	// 	return
// 	// }

// 	fmt.Fprintf(w, "insertAccount id: %d", accountId)
// }

// func login(w http.ResponseWriter, r *http.Request) {
// 	//fmt.Printf("request URI : " + r.RequestURI)

// 	if r.Method != "POST" {
// 		fmt.Fprintf(w, "Sorry, only POST methods are supported~~~~")
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}

// 	// len := r.ContentLength
// 	// body := make([]byte, len)
// 	// r.Body.Read(body)
// 	// fmt.Println(string(body))

// 	decoder := json.NewDecoder(r.Body)
// 	var requestData LoginReq

// 	err := decoder.Decode(&requestData)
// 	if err != nil {
// 		fmt.Fprintf(w, "decoder err: %s", r.Body)
// 		return
// 	}

// 	userId := requestData.UserId
// 	password := requestData.Password

// 	fmt.Printf("id = %s\n", userId)
// 	fmt.Printf("password = %s\n", password)

// 	//pararms := [2]string{id, "fefe"}

// 	passwordInDB := dbQuery(info, loginQuery, userId)

// 	if password != passwordInDB {
// 		fmt.Printf("password is not correct : password = %s\n", password)
// 		return
// 	}

// 	//로그인 세션 만들어서 전송
// 	sessionId := "sessionId1"

// 	var responseData LoginRes
// 	responseData.SessionId = sessionId

// 	encoder := json.NewEncoder(w)
// 	encoder.Encode(responseData)

// 	//http.ServeFile(w, r, "./html/build/index.html")
// }