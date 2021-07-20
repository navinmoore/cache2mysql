package rabbitMQ

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

var rabbitLock sync.Mutex
var RabbitMQServer *RabbitMQ

func RabbitMQInstance() *RabbitMQ {
	// if RabbitMQServer != nil {
	// 	return RabbitMQServer
	// }
	// rabbitLock.Lock()
	// defer rabbitLock.Unlock()
	// if RabbitMQServer != nil {
	// 	return RabbitMQServer
	// }
	RabbitMQServer = NewPubSubRabbitMQ("pubsub", "amqp://yairs:yairs@127.0.0.1:5672/")
	return RabbitMQServer
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	//交换机
	Exchange string
	// Key
	Key string
	// 链接信息
	Mqurl string
}

// func (r *Rabbit MQ)Channel()

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Destory 断开channel和connection
func (r *RabbitMQ) Destory() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func NewRabbitMQ(queueName, exchange, key, mqurl string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		// 链接信息
		Mqurl: mqurl,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "创建channel错误!")
	return rabbitmq
}

func BodyForm(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
