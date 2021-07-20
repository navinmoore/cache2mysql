package rabbitMQ

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

// NewRabbitMQ(queueName, exchange, key, mqurl string)
func NewPubSubRabbitMQ(queueName, mqurl string) *RabbitMQ {
	r := NewRabbitMQ("", "pubsubexchange", "", mqurl)
	return r
}

// ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table)

func (r *RabbitMQ) Pub() {
	body := BodyForm(os.Args)
	// Publish(exchange, key string, mandatory, immediate bool, msg Publishing)

	err := r.channel.ExchangeDeclare(
		r.Exchange, //exchage name
		"fanout",   // 种类 fanout, direct, topic, headers
		true,       // 持久的
		false,      //自动删除
		false,      //
		false,
		nil,
	)
	r.failOnErr(err, "PubSub exchange err=")
	// for i := 0; i < 10; i++ {
	// 	err = r.channel.Publish(
	// 		r.Exchange,
	// 		"",
	// 		false,
	// 		false,
	// 		amqp.Publishing{
	// 			DeliveryMode: amqp.Persistent,
	// 			ContentType:  "text/plain",
	// 			Body:         []byte(body),
	// 		},
	// 	)

	// }

	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	r.failOnErr(err, "Pub err")
	log.Printf(" [x] send %s", body)
	defer r.Destory()

}

func (r *RabbitMQ) Sub() {
	// 定义exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange, //exchage name
		"fanout",   // 种类 fanout, direct, topic, headers
		true,       // 持久的
		false,      //自动删除
		false,      //
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a exchange")

	_, err = r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		true,
		false,
		nil,
	)
	// r.QueueName = q.Name
	r.failOnErr(err, "Failed to declare a queue")

	// QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table)
	//绑定交换器
	err = r.channel.QueueBind(
		r.QueueName,
		"",
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	// 消费
	// Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table)
	msgs, err := r.channel.Consume(
		r.QueueName,
		"", //customer
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Sub Consume")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
