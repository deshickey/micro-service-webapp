package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string
	Type  string
	Options []Option
	Price uint
}

type Option struct {
	gorm.Model
	Size  string
	Quantity int
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	router := gin.Default()

	router.GET("/products", func(c *gin.Context) {
		var products []Product
		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
			return
		}
		c.JSON(http.StatusOK, products)
	})

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

	router.GET("/products/:id", func(c *gin.Context) {
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
	})

	router.PUT("/products/:id", func(c *gin.Context) {
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
		if err := db.Delete(&Product{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
	})

	err = router.Run(":8081")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
