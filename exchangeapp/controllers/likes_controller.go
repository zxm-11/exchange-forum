package controllers

import (
	"encoding/json"
	"exchangeapp/config"
	"exchangeapp/global"
	"exchangeapp/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// 点赞文章
func LikeArticles(ctx *gin.Context) {
	articleId := ctx.Param("id")
	LikeKey := "article:" + articleId + ":likes" //点赞数的key格式为 article:文章ID:likes
	//增加点赞数
	if err := global.RedisDB.Incr(ctx, LikeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 发送点赞消息到RabbitMQ
	publishLikeMessage(articleId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully liked the article",
	})
}

func publishLikeMessage(articleId string) {
	id, err := strconv.ParseUint(articleId, 10, 64)
	if err != nil {
		log.Printf("[LikeMessagge] Invalid articleID:%s,err:%v", articleId, err)
		return
	}
	LikeMsg := models.LikeMessage{
		ArticleID: uint(id),
		Timestamp: time.Now().Unix(),
	}
	body, err := json.Marshal(LikeMsg)
	if err != nil {
		log.Printf("[LikeMessage] Failed to marshal:%v", err)
		return
	}
	err = config.PublishMessage("article.like", body, true)
	if err != nil {
		log.Printf("[LikeMessage] Failed to publish:%v", err)
	}
}

// 获取文章点赞数
func GetArticleLikes(ctx *gin.Context) {
	articleId := ctx.Param("id")
	LikeKey := "article:" + articleId + ":likes"

	likes, err := global.RedisDB.Get(ctx, LikeKey).Result()

	if err == redis.Nil { //如果点赞数不存在,则默认为0
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})

}
