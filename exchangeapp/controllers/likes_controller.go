package controllers

import (
	"exchangeapp/global"
	"net/http"

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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully liked the article",
	})
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
