package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	cachekey :="article:"+ctx.Param("id")+":comments" //定义一个缓存key,用于存储文章详情的缓存key
	global.RedisDB.Del(ctx,cachekey) //删除文章详情的缓存,让下一次访问文章详情时重新查询数据库并更新缓存
	ctx.JSON(http.StatusOK, comment)
}
