package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yellow78/mini-mes/backend/internal/handler"
	"github.com/yellow78/mini-mes/backend/internal/repository"
	"github.com/yellow78/mini-mes/backend/internal/service"
	ws "github.com/yellow78/mini-mes/backend/pkg/websocket"
)

func main() {
	// 載入 .env（本機開發用，Docker 環境由 compose 注入）
	_ = godotenv.Load()

	db := connectDB()
	defer db.Close()

	// 建立 Hub 並啟動廣播迴圈
	hub := ws.NewHub()
	go hub.Run()

	// 初始化各層
	equipRepo   := repository.NewEquipmentRepository(db)
	lotRepo     := repository.NewLotRepository(db)
	alarmRepo   := repository.NewAlarmRepository(db)
	spcRepo     := repository.NewSpcRepository(db)

	equipSvc    := service.NewEquipmentService(equipRepo)
	lotSvc      := service.NewLotService(lotRepo)
	dispatchSvc := service.NewDispatchService(equipRepo, lotRepo)

	equipHandler := handler.NewEquipmentHandler(equipSvc)
	lotHandler   := handler.NewLotHandler(lotSvc, dispatchSvc)
	alarmHandler := handler.NewAlarmHandler(alarmRepo, spcRepo, equipSvc)

	// 設定 Gin router
	r := gin.Default()

	// CORS（開發階段允許所有來源）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWS(c.Writer, c.Request)
	})

	// API 路由
	v1 := r.Group("/api/v1")
	{
		eq := v1.Group("/equipment")
		eq.GET("",           equipHandler.ListEquipments)
		eq.GET("/:id",       equipHandler.GetEquipment)
		eq.PUT("/:id/status", equipHandler.UpdateStatus)
		eq.POST("/:id/hold", equipHandler.HoldEquipment)

		lots := v1.Group("/lots")
		lots.GET("",              lotHandler.ListLots)
		lots.POST("",             lotHandler.CreateLot)
		lots.GET("/:id",          lotHandler.GetLot)
		lots.POST("/:id/dispatch", lotHandler.DispatchLot)

		alarms := v1.Group("/alarms")
		alarms.GET("",                     alarmHandler.ListAlarms)
		alarms.PUT("/:id/acknowledge",     alarmHandler.AcknowledgeAlarm)

		v1.GET("/spc/:equipment_id", alarmHandler.GetSpc)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[Server] 啟動於 :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[Server] 啟動失敗: %v", err)
	}
}

func connectDB() *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "mes_dev"),
		getEnv("DB_USER", "mes"),
		getEnv("DB_PASSWORD", "mes_password"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("[DB] 連線失敗: %v", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	log.Println("[DB] PostgreSQL 連線成功")
	return db
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
