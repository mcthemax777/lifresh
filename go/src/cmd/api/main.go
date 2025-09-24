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

	//repo 생성
	accountRepo := repository.NewAccountRepo(dbConn)
	userRepo := repository.NewUserRepo(dbConn)
	folderRepo := repository.NewFolderRepo(dbConn)
	planRepo := repository.NewPlanRepo(dbConn)

	//service 생성
	authService := service.NewAuthService()
	userService := service.NewUserService(txMgr, accountRepo, userRepo, folderRepo, planRepo)

	//handler 생성
	loginHandler := handler.NewLoginHandler(authService, userService)
	jwtLoginHandler := handler.NewJwtLoginHandler(authService, userService)
	refreshHandler := handler.NewRefreshHandler()
	pingHandler := handler.NewPingHandler()
	getUserHandler := handler.NewGetUserHandler(userService)
	createUserHandler := handler.NewCreateUserHandler(userService)
	createFolderHandler := handler.NewCreateFolderHandler(userService)
	updateFolderHandler := handler.NewUpdateFolderHandler(userService)
	deleteFolderHandler := handler.NewDeleteFolderHandler(userService)

	// Auth (미들웨어에서 자동 스킵)
	authGroup := r.Group(define.RouterAuth)
	{
		authGroup.POST(define.ApiLogin, loginHandler.ApiCall)
		authGroup.POST(define.ApiJwtLogin, jwtLoginHandler.ApiCall)
		authGroup.POST(define.ApiRefresh, refreshHandler.ApiCall)
		//authGroup.POST(define.ApiPing, pingHandler.ApiCall)

		//authGroup.POST(define.ApiRefresh, api.ApiCall)
	}

	apiGroup := r.Group(define.RouterApi)
	{
		// JWT 미들웨어
		apiGroup.Use(auth.JWTAuthSkipper())
		apiGroup.GET(define.ApiPing, pingHandler.ApiCall)
		apiGroup.POST(define.ApiGetUser, getUserHandler.ApiCall)
		apiGroup.POST(define.ApiCreateUser, createUserHandler.ApiCall)
		apiGroup.POST(define.ApiCreateFolder, createFolderHandler.ApiCall)
		apiGroup.POST(define.ApiUpdateFolder, updateFolderHandler.ApiCall)
		apiGroup.POST(define.ApiDeleteFolder, deleteFolderHandler.ApiCall)
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

	addr := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("Listening on", addr)
	_ = r.Run(addr)
}
