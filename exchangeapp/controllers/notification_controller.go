package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Getnotifications 获取用户通知列表
func GetNotifications(ctx *gin.Context) {
	username := ctx.GetString("username")
	var notifications []models.InAppNotification

	if err := global.Db.Where("username=?", username).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}

// GetUnreadCount 获取用户未读通知数量
func GetUnreadCount(ctx *gin.Context) {
	username := ctx.GetString("username")
	var count int64

	if err := global.Db.Model(&models.InAppNotification{}).Where("username=? AND is_read=?", username, false).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Unread_count": count,
	})
}

// MarkNotificationRead 标记通知为已读
func MarkNotificationRead(ctx *gin.Context) {
	id := ctx.Param("id")
	username := ctx.GetString("username")
	result := global.Db.Model(&models.InAppNotification{}).Where("id=? AND username=?", id, username).Update("is_read", true)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "notification not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "notification marked as read",
	})
}
