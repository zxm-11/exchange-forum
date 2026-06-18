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

// CreateComment函数用于创建评论
func CreateComment(ctx *gin.Context) {

	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64) //从URL参数中获取文章ID
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}
	comment.ArticleID = uint(id)                 //将文章ID赋值给评论的ArticleID字段
	comment.Username = ctx.GetString("username") //从上下文中获取用户名，并赋值给评论的Username字段
	if err := global.Db.AutoMigrate(&comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := global.Db.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	cachekey := "article:" + ctx.Param("id") + ":comments" //定义一个缓存key,用于存储文章详情的缓存key
	global.RedisDB.Del(ctx, cachekey)                      //删除文章详情的缓存,让下一次访问文章详情时重新查询数据库并更新缓存

	//发布消息给RabbitMQ,通知其他服务有新的评论创建
	publishCommentNotification(comment)

	ctx.JSON(http.StatusOK, comment)
}

func publishCommentNotification(comment models.Comment) {
	var article models.Article
	if err := global.Db.First(&article, comment.ArticleID).Error; err != nil {
		log.Printf("[CommentNotify] Failed to find ArticleID:%d:%v", comment.ArticleID, err)
		return
	}

	notifMsg := models.CommentNotification{
		ArticleID:      comment.ArticleID,
		ArticleTitle:   article.Title,
		CommentID:      comment.ID,
		CommentAuthor:  comment.Username,
		CommentContent: comment.Content,
		Timestamp:      time.Now().Unix(),
	}
	body, err := json.Marshal(notifMsg)
	if err != nil {
		log.Printf("[CommentNotify] Failed to marshal CommentNotification:%v", err)
		return
	}

	// 发布消息给RabbitMQ,通知其他服务有新的评论创建
	err = config.PublishMessage("comment_notification", body, false)
	if err != nil {
		log.Printf("[CommentNotify] Failed to publish CommentNotification:%v", err)
		return
	}
}

// GetCommentsByArticleID函数用于获取文章下的评论列表
func GetCommentsByArticleID(ctx *gin.Context) {
	cachekey := "article:" + ctx.Param("id") + ":comments"

	cachedata, err := global.RedisDB.Get(ctx, cachekey).Result()
	if err == redis.Nil {
		var comments []models.Comment
		if err := global.Db.Where("article_id=?", ctx.Param("id")).Order("created_at DESC").Find(&comments).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		cacheJSON, err := json.Marshal(comments) //将评论列表转换为JSON格式的字节切片
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := global.RedisDB.Set(ctx, cachekey, cacheJSON, time.Minute*10).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, comments)

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		var comments []models.Comment
		if err := json.Unmarshal([]byte(cachedata), &comments); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, comments) //如果缓存中有数据,直接返回缓存中的数据
	}
}

// DeleteComment函数用于删除评论
func DeleteComment(ctx *gin.Context) {
	comment_id := ctx.Param("comment_id")
	var comment models.Comment
	if err := global.Db.Where("id=?", comment_id).First(&comment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if comment.Username == ctx.GetString("username") { //检查评论的用户名是否与当前用户的用户名匹配,如果不匹配,返回403 Forbidden错误}
		if err := global.Db.Delete(&comment).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		cachekey := "article:" + strconv.FormatUint(uint64(comment.ArticleID), 10) + ":comments"
		if err := global.RedisDB.Del(ctx, cachekey).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "comment deleted successfully",
		})
	} else {
		ctx.JSON(
			http.StatusForbidden, gin.H{"error": "you are not authorized to delete this comment"})
		return
	}
}
