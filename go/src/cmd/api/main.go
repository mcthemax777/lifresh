// main.go
package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lifresh/auth"
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/repository"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/handler"
	"lifresh/internal/txmgr"
	"log"
	"os"
	"time"
)

//
//// --- ID generator -----------------------------------------------------------
//var node *snowflake.Node
//
//func mustInitSnowflake() {
//	var nodeID int64 = 1
//	if v := os.Getenv("SNOWFLAKE_NODE_ID"); v != "" {
//		if id, err := snowflake.NewNodeFromString(v); err == nil {
//			node = id
//			return
//		}
//	}
//	var err error
//	node, err = snowflake.NewNode(nodeID)
//	if err != nil {
//		log.Fatalf("snowflake init failed: %v", err)
//	}
//}
//
//func newID() models.SnowflakeID { return models.SnowflakeID(node.Generate().Int64()) }
//
//// --- DTOs -------------------------------------------------------------------
//// 대부분의 필드는 models의 JSON 태그를 그대로 사용합니다. 필요 시 전용 요청 DTO를 둡니다.
//
//// CreateUserRequest: id가 비어있으면 서버가 생성
//type CreateUserRequest struct {
//	ID         *models.SnowflakeID `json:"id"`
//	Nickname   string              `json:"nickname" binding:"required"`
//	Bio        string              `json:"bio"`
//	ProfileURL *string             `json:"profileUrl"`
//}
//
//// CreateFolderRequest
//type CreateFolderRequest struct {
//	ID       *models.SnowflakeID `json:"id"`
//	ParentID *models.SnowflakeID `json:"parentId"`
//	Order    int                 `json:"order"`
//	Color    int                 `json:"color" binding:"required"`
//	Name     string              `json:"name" binding:"required"`
//}
//
//// CreatePlanRequest
//type CreatePlanRequest struct {
//	ID          *models.SnowflakeID `json:"id"`
//	ParentID    *models.SnowflakeID `json:"parentId"`
//	Order       int                 `json:"order"`
//	Color       int                 `json:"color" binding:"required"`
//	Name        string              `json:"name" binding:"required"`
//	Description string              `json:"description"`
//	DateType    models.DateType     `json:"dateType" binding:"required"`
//	StartDate   *string             `json:"startDate"`  // ISO8601 문자열
//	FinishDate  *string             `json:"finishDate"` // ISO8601 문자열
//}
//
//// CreateRecordRequest
//type CreateRecordRequest struct {
//	ID     *models.SnowflakeID `json:"id"`
//	Values models.JSONMap      `json:"values" binding:"required"`
//}
//
//// --- Handlers ---------------------------------------------------------------
//func createUser(c *gin.Context) {
//	var req CreateUserRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	u := models.User{
//		Nickname:   req.Nickname,
//		Bio:        req.Bio,
//		ProfileURL: req.ProfileURL,
//	}
//	if req.ID != nil {
//		u.ID = *req.ID
//	} else {
//		u.ID = newID()
//	}
//	if err := gormDB().Create(&u).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, u)
//}
//
//func getUser(c *gin.Context) {
//	idParam := c.Param("id")
//	var id models.SnowflakeID
//	if err := id.UnmarshalJSON([]byte("\"" + idParam + "\"")); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
//		return
//	}
//	var u models.User
//	if err := gormDB().First(&u, "id = ?", id).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, u)
//}
//
//func createFolder(c *gin.Context) {
//	var req CreateFolderRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	f := models.Folder{ParentID: req.ParentID, Order: req.Order, Color: req.Color, Name: req.Name}
//	if req.ID != nil {
//		f.ID = *req.ID
//	} else {
//		f.ID = newID()
//	}
//	if err := gormDB().Create(&f).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, f)
//}
//
//func getFolder(c *gin.Context) {
//	idParam := c.Param("id")
//	var id models.SnowflakeID
//	if err := id.UnmarshalJSON([]byte("\"" + idParam + "\"")); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
//		return
//	}
//	var f models.Folder
//	if err := gormDB().Preload("ChildrenFolders").Preload("Plans").First(&f, "id = ?", id).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, f)
//}
//
//func createPlan(c *gin.Context) {
//	var req CreatePlanRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	p := models.Plan{ParentID: req.ParentID, Order: req.Order, Color: req.Color, Name: req.Name, Description: req.Description, DateType: req.DateType}
//	if req.ID != nil {
//		p.ID = *req.ID
//	} else {
//		p.ID = newID()
//	}
//	// Parse dates if provided
//	if req.StartDate != nil {
//		if t, err := timeParse(*req.StartDate); err == nil {
//			p.StartDate = &t
//		}
//	}
//	if req.FinishDate != nil {
//		if t, err := timeParse(*req.FinishDate); err == nil {
//			p.FinishDate = &t
//		}
//	}
//	if err := gormDB().Create(&p).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, p)
//}
//
//func getPlan(c *gin.Context) {
//	idParam := c.Param("id")
//	var id models.SnowflakeID
//	if err := id.UnmarshalJSON([]byte("\"" + idParam + "\"")); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
//		return
//	}
//	var p models.Plan
//	if err := gormDB().Preload("Goals").Preload("RecordFields").Preload("RepeatRules").Preload("Records").Preload("Statistics").First(&p, "id = ?", id).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, p)
//}
//
//func createRecord(c *gin.Context) {
//	planIDParam := c.Param("id")
//	var planID models.SnowflakeID
//	if err := planID.UnmarshalJSON([]byte("\"" + planIDParam + "\"")); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planId"})
//		return
//	}
//	var req CreateRecordRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	r := models.Record{PlanID: planID, Values: req.Values}
//	if req.ID != nil {
//		r.ID = *req.ID
//	} else {
//		r.ID = newID()
//	}
//	if err := gormDB().Create(&r).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, r)
//}
//
//func listRecords(c *gin.Context) {
//	planIDParam := c.Param("id")
//	var planID models.SnowflakeID
//	if err := planID.UnmarshalJSON([]byte("\"" + planIDParam + "\"")); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planId"})
//		return
//	}
//	var list []models.Record
//	if err := gormDB().Where("plan_id = ?", planID).Order("created_at asc").Find(&list).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, list)
//}

/* =========================
   서버 구동
   ========================= */

func main() {
	dbConn := db.InitDB()
	txMgr := &txmgr.Manager{
		DB:     dbConn,
		Flavor: txmgr.MySQLFlavor{}, // 또는 PostgresFlavor
		OnBegin: func(opts *sql.TxOptions) {
			log.Printf("[tx] begin iso=%v ro=%v", opts.Isolation, opts.ReadOnly)
		},
		OnCommit: func(dur time.Duration, err error) {
			//status := "ok"
			//if err != nil {
			//	status = "err"
			//}
			//metrics.ObserveTx(dur, status) // 예시
		},
		OnRollback: func(err error) {
			log.Printf("[tx] rollback err=%v", err)
		},
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS (프론트 Authorization 헤더/OPTIONS 허용)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:   []string{"Content-Length"},
	}))

	// JWT 미들웨어
	//r.Use(auth.JWTAuthSkipper())

	//repo 생성
	userRepo := repository.NewUserRepo(dbConn)
	accountRepo := repository.NewAccountRepo(dbConn)

	//service 생성
	authService := service.NewAuthService()
	userService := service.NewUserService(txMgr, accountRepo, userRepo)

	//handler 생성
	loginHandler := handler.NewLoginHandler(authService, userService)
	refreshHandler := handler.NewRefreshHandler()
	pingHandler := handler.NewPingHandler()

	// Auth (미들웨어에서 자동 스킵)
	authGroup := r.Group(define.RouterAuth)
	{
		authGroup.POST(define.ApiLogin, loginHandler.ApiCall)
		authGroup.POST(define.ApiRefresh, refreshHandler.ApiCall)
		//authGroup.POST(define.ApiPing, pingHandler.ApiCall)

		//authGroup.POST(define.ApiRefresh, api.ApiCall)
	}

	apiGroup := r.Group(define.RouterApi)
	{
		// JWT 미들웨어
		apiGroup.Use(auth.JWTAuthSkipper())
		apiGroup.GET(define.ApiPing, pingHandler.ApiCall)
		// Root / Items
		//apiGroup.GET("/v1/root", getRoot)
		//apiGroup.PUT("/v1/root", putRoot)
		//
		//apiGroup.GET("/v1/items/:id", getItem)
		//apiGroup.POST("/v1/items", createItem)
		//apiGroup.PUT("/v1/items/:id", putItem)
		//apiGroup.DELETE("/v1/items/:id", deleteItem)
		//
		//// Records
		//apiGroup.POST("/v1/plans/:planId/records", addPlanRecord)
		//apiGroup.PUT("/v1/plans/:planId/records/:recordId", updatePlanRecord)
		//apiGroup.DELETE("/v1/plans/:planId/records/:recordId", deletePlanRecord)
		//
		//// Statistics
		//apiGroup.POST("/v1/plans/:planId/statistics", addStatistics)
		//apiGroup.PUT("/v1/plans/:planId/statistics/:configId", updateStatistics)
		//apiGroup.DELETE("/v1/plans/:planId/statistics/:configId", deleteStatistics)

	}

	//// Root / Items
	//r.GET("/v1/root", getRoot)
	//r.PUT("/v1/root", putRoot)
	//
	//r.GET("/v1/items/:id", getItem)
	//r.POST("/v1/items", createItem)
	//r.PUT("/v1/items/:id", putItem)
	//r.DELETE("/v1/items/:id", deleteItem)
	//
	//// Records
	//r.POST("/v1/plans/:planId/records", addPlanRecord)
	//r.PUT("/v1/plans/:planId/records/:recordId", updatePlanRecord)
	//r.DELETE("/v1/plans/:planId/records/:recordId", deletePlanRecord)
	//
	//// Statistics
	//r.POST("/v1/plans/:planId/statistics", addStatistics)
	//r.PUT("/v1/plans/:planId/statistics/:configId", updateStatistics)
	//r.DELETE("/v1/plans/:planId/statistics/:configId", deleteStatistics)
	//
	//r.GET("/v1/api/:name", api.ApiCall)
	//r.POST("/v1/api/:name", api.ApiCall)
	//r.PUT("/v1/api/:name", api.ApiCall)
	//r.DELETE("/v1/api/:name", api.ApiCall)

	addr := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("Listening on", addr)
	_ = r.Run(addr)
}
