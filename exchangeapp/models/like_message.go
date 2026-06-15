package models

type LikeMessage struct {
	ArticleID uint `json:"article_id"`
	Timestamp int64  `json:"timestamp"` // 点赞时间戳
}
