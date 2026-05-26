package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ArticleID uint   `json:"article_id" binding:"required"` //文章ID,必填
	Username  string `json:"username" binding:"required"`   //用户名,必填
	Content   string `json:"content" binding:"required"`    //评论内容，必填
}
