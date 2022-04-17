package main

import (
	"fmt"
	// "api"
	// "github.com/gin-gonic/contrib/static"
	// "github.com/gin-gonic/gin"
	// "github.com/itsjamie/gin-cors"
	"sort"
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func init() {
}

type myDataType struct {
    num 	int
    count  	int
	cha		int
}

func main() {

	f, err := excelize.OpenFile("C:/workspace/lifresh/go/lotto.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
	mySlice := make([]myDataType, 0)

	var count[46]int

	for i := 5; i < 255; i++ {
		for j := 66; j < 72; j++ {
			s := string(j)
			num1, err := f.GetCellValue("Sheet1", s + strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println(num1)

			result, _ := strconv.Atoi(num1);

			count[result] += 1;
		}
	}

	for i := 1; i < 46; i++ {
		mySlice = append(mySlice, myDataType{i, count[i], 1})
	}

	sort.Slice(mySlice, func(i, j int) bool {
        return mySlice[i].count < mySlice[j].count
    })
    fmt.Println(mySlice)



	// router := gin.Default()

	// router.Use(cors.Middleware(cors.Config{
	// 	Origins:        "*",
	// 	Methods:        "GET, PUT, POST, DELETE",
	// 	RequestHeaders: "Origin, Authorization, Content-Type",
	// 	ExposedHeaders: "",
	// 	MaxAge: 50 * 1000,
	// 	Credentials: true,
	// 	ValidateHeaders: false,
	// }))

	// router.Use(static.Serve("/", static.LocalFile("../html", true)))
	// // router.LoadHTMLGlob("templates/*")
	// //router.GET("/", hello)
	// router.POST("/api/:name", api.ApiCall)

	// router.Run(":8000")
}