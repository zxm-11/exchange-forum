package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ArticleID uint   `json:"article_id"` //文章ID
	Username  string `json:"username"`   //用户名
	Content   string `json:"content" binding:"required"`    //评论内容，必填
}
