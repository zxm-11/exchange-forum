package models

import "gorm.io/gorm"

type CommentNotification struct { //评论通知消息
	ArticleID      uint   `json:"article_id"`
	ArticleTitle   string `json:"article_title"`
	CommentID      uint   `json:"comment_id"`
	CommentAuthor  string `json:"comment_author"`
	CommentContent string `json:"comment_content"` //评论摘要
	Timestamp      int64  `json:"timestamp"`
}

type InAppNotification struct { //站内通知(持久化到数据库)
	gorm.Model
	Username string `json:"username"` //用户名
	Title    string `json:"title"`
	Content  string `json:"content"`                      //通知内容
	Link     string `json:"link"`                         //跳转链接
	IsRead   bool   `json:"is_read" gorm:"default:false"` //是否已读
}
