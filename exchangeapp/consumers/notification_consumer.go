package consumers

import (
	"encoding/json"
	"exchangeapp/config"
	"exchangeapp/models"
	"exchangeapp/services"
	"time"

	"log"

	"github.com/streadway/amqp"
)

func StartNotificationConsumer(prefetch int) {
	log.Printf("[NotifyConsumer] Starting with prefetch=%d", prefetch)

	conn := config.GetRabbitMQConn()
	if conn == nil {
		log.Fatalf("[NotifyConsumer] RabbitMQ connectiion is nil")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[NotifyConsumer] Faile to create channel:%v", err)
	}
	defer ch.Close()

	err = ch.Qos(
		prefetch,
		0,
		false,
	)
	if err != nil {
		log.Fatalf("[NotifyConsumer] Failed to set Qos:%v", err)
	}

	msgs, err := ch.Consume(
		config.Notificationqueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[NotifyConsumer] Failed to start consumer:%v", err)
	}

	for msg := range msgs {
		processMessage(msg) // 处理消息
	}

	log.Println("[NotifyConsumer] Consumer closed successfully")
}

// 处理单条消息
func processMessage(msg amqp.Delivery) {
	var notifyMsg models.CommentNotification
	if err := json.Unmarshal(msg.Body, &notifyMsg); err != nil {
		log.Printf("[NotifyConsumer] Failed to unmarshal:%v,body:%s", err, string(msg.Body))
		//消息格式错误,直接进入死信队列
		msg.Nack(false, false)
		return
	}
	log.Printf("[NotifyConsumer] Processing message:articleID=%d,commentID=%d,author=%s",
		notifyMsg.ArticleID,
		notifyMsg.CommentID,
		notifyMsg.CommentAuthor)

	siteErr := services.SendSiteNotification(&notifyMsg)
	if siteErr != nil {
		log.Printf("[NotifyConsumer] Failed to send site notification:%v", siteErr)

		//判断错误类型,如果是网络错误,则重新发送消息
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			if reTryErr := services.SendSiteNotification(&notifyMsg); reTryErr == nil {
				log.Printf("[NotifyConsumer] Retried successed")
				msg.Ack(false)
				return
			}
		}

		//3次重试失败,进入死信队列
		log.Printf("[NotifyConsumer] All retried failed for artilceID=%d,sending to DLQ", notifyMsg.ArticleID)

		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
	log.Printf("[NotifyConsumer] Notification processed successfully for article=%d", notifyMsg.ArticleID)
}
