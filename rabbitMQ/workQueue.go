package rabbitMQ

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// 常见简单模式下的RabbitMQ实例
func NewRabbitMQWorkQueue(queueName, mqurl string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "", mqurl)
}

func (r *RabbitMQ) Task() {
	// QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table)
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名
		false,       // 持久化
		false,       // 是否自动删除
		false,       // 是否独立
		false,       //
		nil,
	)

	failOnError(err, "Failed to open a queue")
	body := BodyForm(os.Args)
	// Publish(exchange, key string, mandatory, immediate bool, msg Publishing)
	err = r.channel.Publish(
		"",          //exchange
		r.QueueName, // routing key
		false,       // mandatory 强制
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
	r.Destory()
}

// 帮助函数检测每一个amqp调用
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) Work() {
	// QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table)
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to open a queue")

	// Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table)
	// 消费队列
	msgs, err := r.channel.Consume(
		r.QueueName, //queue
		"",          // consumer
		false,       // 自动确认
		false,       //
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// 具体的事情
			time.Sleep(time.Second * 2)
			log.Printf("Done")
			//
			d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for messages.To exit press CTRL+C")
	<-forever
	r.Destory()
}
