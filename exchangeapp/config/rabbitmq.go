package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// 队列名称常量
const (
	Likequeue         = "article_likes"
	Notificationqueue = "comment_notifications" // 评论通知队列
	DeadLetterqueue   = "dead_letter"           // 死信队列((处理失败消息))
)

// 交换机名称常量
const (
	DefaultExchange = "exchangeappexchange" //默认交换机)(名称)
)

var (
	rabbitConn   *amqp.Connection
	rabbituChan  *amqp.Channel
	rabbitMu     sync.Mutex
	rabbitNotify chan *amqp.Error //Connection 异常通知 channel
)

func InitRabbitMQ() {
	cfg := AppConfig.rabbitMQ
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost)

	var err error
	rabbitConn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ:%v", err)
	}
	rabbituChan, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel:%v", err)
	}

	declareExchange(rabbituChan)

	declareQueue(rabbituChan)

	// 注册连接异常回调 —— 当连接断开时，rabbitNotify 会收到通知
	rabbitNotify = rabbitConn.NotifyClose(make(chan *amqp.Error))
	log.Printf("RabbitMQ initialzed successfully")
}

func declareExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		DefaultExchange,
		"topic",
		true, //持久化交换机
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange:%v", err)
	}
}

func declareQueue(ch *amqp.Channel) {

	//1.死信队列
	_, err := ch.QueueDeclare(
		DeadLetterqueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": 86400000, // 消息在死信队列中的存活时间: 24小时 (毫秒)
		},
	)
	if err != nil {
		log.Fatalf("Failed to declare queue:%v", err)
	}
	ch.QueueBind(DeadLetterqueue, "dead.letter", DefaultExchange, false, nil)

	//2.点赞队列
	_, err = ch.QueueDeclare(
		Likequeue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    DefaultExchange, //死信回到defaultexchange交换机
			"x-dead-letter-routing-key": DeadLetterqueue, //死信路由键
		},
	)
	if err != nil {
		log.Fatalf("Failed to declare queue:%v", err)
	}
	ch.QueueBind(Likequeue, "artilce.like", DefaultExchange, false, nil)

	//3.评论通知队列
	_, err = ch.QueueDeclare(
		Notificationqueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    DefaultExchange,
			"x-dead-letter-routing-key": DeadLetterqueue,
		},
	)
	if err != nil {
		log.Fatalf("Failed to declare queue:%v", err)
	}
	ch.QueueBind(Notificationqueue, "comment.notify", DefaultExchange, false, nil)
}

// GetChannel 线程安全地获取 RabbitMQ Channel
func GetChannel() (*amqp.Channel, func()) {
	rabbitMu.Lock()
	return rabbituChan, func() {
		rabbitMu.Unlock()
	}
}

// ReconnectRabbitMQ处理RabbitMQ断线重连
func ReconnectRabbitMQ() {
	for {
		err, ok := <-rabbitNotify
		if !ok {
			// 连接已关闭，退出循环
			return
		}

		log.Fatalf("RabbitMQ connection closed:%v,attempting reconnect...", err)

		//无限重连直到连接成功
		for {
			time.Sleep(3 * time.Second)
			log.Println("Attempting to reconnect to RabbitMQ...")
			cfg := AppConfig.rabbitMQ
			dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost)
			conn, DialErr := amqp.Dial(dsn)
			if DialErr != nil {
				log.Fatalf("Failed to reconnect to RabbitMQ:%v,Retrying...", DialErr)
				continue
			}
			ch, chErr := rabbitConn.Channel()
			if chErr != nil {
				log.Fatalf("Failed to create channel:%v,Retrying...", chErr)
				continue
			}
			//原子更新全局连接
			rabbitMu.Lock()
			if rabbituChan != nil {
				rabbituChan.Close()
			}
			if rabbitConn != nil {
				rabbitConn.Close()
			}
			rabbitConn = conn
			rabbituChan = ch
			rabbitMu.Unlock()

			//重新声明交换机和队列
			declareExchange(ch)
			declareQueue(ch)

			//重新注册异常通知
			rabbitNotify = conn.NotifyClose(make(chan *amqp.Error))

			log.Printf("RabbitMQ reconnected successfully...")
			break
		}
	}
}

// 优雅的关闭RabbitMQ连接
func CloseRabbitMQ() {
	if rabbituChan != nil {
		rabbituChan.Close()
	}
	if rabbitConn != nil {
		rabbitConn.Close()
	}
	log.Println("RabbitMQ connection closed successfully")

}

// PubilshMesssage
func PublishMessage(routingkey string, body []byte, persistent bool) error {
	ch, release := GetChannel()
	defer release() //释放Channel

	deliveryMode := amqp.Transient //非持久化消息(消息仅在内存中，性能更好但可能丢失)
	if persistent {
		deliveryMode = amqp.Persistent //持久化消息(消息被写入磁盘，适合不能丢失的业务数据 (如点赞))
	}

	return ch.Publish(
		DefaultExchange,
		routingkey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: deliveryMode,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
