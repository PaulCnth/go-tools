package main

import (
	"fmt"
	"net/http"
	"strconv"

	"time"

	"dev.com/short-url/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

const (
	JSON_SUCCESS int = 1
	JSON_ERROR   int = 0
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// 用户访问给出重定向
func visit(c *gin.Context) {
	hash := c.Param("hash")
	url, err := rdb.Get(hash).Result()
	if err != nil || len(url) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "Not found",
		})
	}

	// 如果存在，则重定向
	c.Redirect(http.StatusMovedPermanently, url)
}

// 返回当前计数器的值，自动+1计数
func getCounter() int {
	rdb.Incr("counter")
	id, _ := rdb.Get("counter").Result()
	_id, _ := strconv.Atoi(id)
	return _id
}

func add(c *gin.Context) {

	target := c.PostForm("target")
	expire := c.PostForm("expire")
	_expire, err := strconv.Atoi(expire)

	// todo:验证上述字段有效性
	id := getCounter()

	hash := helper.DecToAny(id)

	err = rdb.Set(hash, target, time.Duration(_expire)*time.Second).Err()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    JSON_SUCCESS,
		"message":   "ok",
		"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, hash),
	})

}

func main() {
	r := gin.Default()

	r.POST("/", add)       // 添加条目
	r.GET("/:hash", visit) // 用户使用，访问此链接，重定向到目标网址

	r.Run(":9090")
}
