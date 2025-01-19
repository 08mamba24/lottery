package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:980201@tcp(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 10)

	type Lucky struct {
		gorm.Model
		Number     string `gorm:"varchar(20); not null" json:"number" binding:"required"`
		Name       string `gorm:"varchar(20); not null" json:"name" binding:"required"`
		PrizeLevel int    `gorm:"int(1); not null" json:"prizeLevel" binding:"required"`
	}

	db.AutoMigrate(&Lucky{})
	r := gin.Default()
	port := ":8080"
	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // 允许所有域名访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 路由定义
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request received successfully!"})
	})
	r.POST("/lottery", func(c *gin.Context) {
		var lucky Lucky
		if err := c.ShouldBindJSON(&lucky); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&lucky)
		// 打印中文调试日志
		fmt.Printf("收到数据: %+v\n", lucky)
		c.JSON(200, gin.H{"data": lucky})
	})

	r.Run(port)
}
