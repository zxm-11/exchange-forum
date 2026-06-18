package services

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"log"
)

func SendSiteNotification(notif *models.CommentNotification) error {
	var article models.Article
	if err := global.Db.First(&article, notif.ArticleID).Error; err != nil {
		return fmt.Errorf("Failed to find article %d:%w", notif.ArticleID, err)
	}
	notifRecord := models.InAppNotification{
		Username: article.Username,
		Title:    article.Title,
		Content: fmt.Sprintf("用户%s评论了文章《%s》:%s",
			notif.CommentAuthor,
			notif.ArticleTitle,
			truncateContent(notif.CommentContent, 100),
		),
		Link:   fmt.Sprintf("/articles/%d", notif.ArticleID),
		IsRead: false,
	}

	// 写入数据库
	if err := global.Db.Create(&notifRecord).Error; err != nil {
		return fmt.Errorf("Failed to create notification%w", err)
	}

	log.Printf("[Notification] Created site notification #%d for article %d", notifRecord.ID, notif.ArticleID)
	return nil

}

func truncateContent(content string, maxLen int) string {
	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}
	return string(runes[:maxLen]) + "..."
}
