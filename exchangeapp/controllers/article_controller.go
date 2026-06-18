package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var cachekey = "articles" //定义一个全局变量cachekey,用于存储文章列表的缓存key

// CreateAritcle函数用于创建文章
func CreateAritcle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	article.Username = ctx.GetString("username")

	//自动迁移数据库表 (确保数据库中有这张表，并且表的结构跟 article 结构体一致,则更新,如果表不存在,则创建)	(DDL)
	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//创建文章
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//创建文章后,刷新Redis缓存中的文章列表
	if err := global.RedisDB.Del(ctx, cachekey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, article)
}

// GetArticles函数用于获取所有文章
func GetArticles(ctx *gin.Context) {

	cacheData, err := global.RedisDB.Get(ctx, cachekey).Result() //尝试从Redis缓存中获取文章列表
	if err == redis.Nil {                                        //如果缓存中没有数据,则从数据库中查询文章列表
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		articleJSON, err := json.Marshal(articles) //将文章列表序列化为JSON字符串
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := global.RedisDB.Set(ctx, cachekey, articleJSON, time.Minute*10).Err(); err != nil { //将序列化后的文章列表存储到Redis缓存中,设置过期时间为10分钟
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		//返回文章列表
		ctx.JSON(http.StatusOK, articles)

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		var articles []models.Article
		if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, articles)
		return
	}
}

// GetArticleById函数用于获取指定文章的文章
func GetArticleById(ctx *gin.Context) {
	id := ctx.Param("id")
	var article models.Article
	if err := global.Db.Where("id=?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, article)
}
