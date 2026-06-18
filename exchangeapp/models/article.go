package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title    string `binding:"required"` //标题,必填
	Content  string `binding:"required"` //内容,必填
	Preview  string `binding:"required"` //预览,必填
	Likes    int    `default:"0"`        //点赞数,默认值为0
	Username string `default:""`         //作者用户名,默认值为空字符串
}
