package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := gin.Default()

	router.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, err := rdb.Get(c, key).Result()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.String(http.StatusOK, value)
	})

	router.PUT("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Query("value")
		if err := rdb.Set(c, key, value, 0).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set key"})
			return
		}
		c.Status(http.StatusOK)
	})

	err := router.Run(":8082")
	if err != nil {
		log.Fatal("Failed to start redis cache: ", err)
	}
}
