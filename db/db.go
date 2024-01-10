package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

func main() {
	dsn := "user:password@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := gin.Default()

	store := persistence.NewRedisCache(rdb)

	router.GET("/products/:id", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var product Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	}))

	router.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	router.PUT("/products/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		}
		var product Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if err := db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	router.DELETE("/products/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var product Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err := db.Delete(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}
		c.Status(http.StatusOK)
	})

	err = router.Run(":8081")
	if err != nil {
		log.Fatal("Failed to start gin: ", err)
	}
}
