package consumers

import (
	"encoding/json"
	"exchangeapp/config"
	"exchangeapp/global"
	"exchangeapp/models"
	"log"
	"time"
)

type likesBuffer map[uint]int

// 启动点赞消费者
func LikeConsumer(prefetch int, flushInterval time.Duration) {
	log.Printf("[LikeConsumer] Starting with prefetch=%d,flushInterval=%v", prefetch, flushInterval) //prefetch:每次从队列中获取的消息数量,flushInterval:刷新时间间隔

	conn := config.GetRabbitMQConn()
	if conn == nil {
		log.Fatalf("[LikeConsumer] Connection is nil")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[LikeConsumer] Failed to open a channel:%v", err)
	}

	defer ch.Close()

	//配置Qos
	err = ch.Qos(prefetch, 0, false)
	if err != nil {
		log.Fatalf("[LikeConsumer] Failed to set Qos:%v", err)
	}

	//消费
	msgs, err := ch.Consume(
		config.Likequeue,
		"",
		false, //auto-ack关闭(手动确认模式)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[LikeConsumer] Failed to consume:%v", err)
	}

	//初始化聚合缓冲区
	buffer := make(likesBuffer)

	//启动定时刷库
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	log.Printf("[LikeConsumer] Timer started with flushInterval=%v,waiting for messages...", flushInterval)

	var likeMsg models.LikeMessage

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Println("[LikeConsumer] Message channel closed, flushing remaining data...")
				flushToDatabase(buffer)
				return
			}

			if err := json.Unmarshal(msg.Body, &likeMsg); err != nil {
				log.Printf("[LikeConsumer] Failed to unmarshal message: %v, body: %s", err, string(msg.Body))
				msg.Nack(false, false) //requeue=false → 进入 DLQ
				continue
			}
			//聚合
			buffer[likeMsg.ArticleID]++
			msg.Ack(false) //手动确认消费成功

			log.Printf("[LikeConsumer] Processed message for articleID: %d, count: %d", likeMsg.ArticleID, buffer[likeMsg.ArticleID])

		case <-ticker.C:
			//定时器触发，刷新数据库
			if len(buffer) > 0 {
				log.Printf("[LikeConsumer] Timer triggered,flushing%d articles...", len(buffer))
				flushToDatabase(buffer)
				buffer = make(likesBuffer) //重新初始化缓冲区
			}
		}
	}
 
}

// 将聚合的点赞增量批量写入 MySQL数据库
func flushToDatabase(buffer likesBuffer) {
	for articleID, delta := range buffer {
		result := global.Db.Model(&models.Article{}).Where("id=?", articleID).UpdateColumn("likes", global.Db.Raw("likes+?", delta))

		if result.Error != nil {
			log.Printf("[LikeConsumer] failed to update likes for articleID:%d,err:%v", articleID, result.Error)
			continue
		}
		if result.RowsAffected == 0 { //改变行量
			log.Printf("[LikeConsumer] Article%d not found,skipping like delta%d", articleID, delta)
			continue
		}

		log.Printf("[LikeConsumer] Updated likes for articleID:%d, delta:%d", articleID, delta)
	}

}




