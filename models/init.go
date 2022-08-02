package models

import (
	"context"
	//"getcharzp.cn/define"
	// "github.com/go-redis/redis/v8"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Redis *redis.Client
var Ctx = context.Background()

func Init() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/hykoj?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("open db error:  ", err)
		fmt.Println("open db error: ", err)
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0,
	})
	_, err = Redis.Ping(Ctx).Result()
	if err != nil {
		log.Println("redis 连接错误,  ", err)
		return
	}
}
