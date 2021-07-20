package rabbitMQ

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// 常见简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName, mqurl string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "", mqurl)
}

// 生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	// 申请队列，如果队列不存在自动创建，如果存在则跳过创建 保证队列存在，消息能发送到队列
	// name string, durable, autoDelete, exclusive, noWait bool, args Table
	fmt.Println("PublishSimple", r.QueueName)
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 自动删除
		false,
		//排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println("PublishSimple", err)
	}

	//发送消息到队列
	// func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error {
	r.channel.Publish(
		r.Exchange,  //exchange
		r.QueueName, // routing key??
		//mandatory 如果为true, 会根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		// immediate 如果为true, 当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.Destory()
}

// 消费
func (r *RabbitMQ) ConsumeSimple() {
	// QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table)
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//持久化
		false,
		// 自动删除
		false,
		// 排他性
		false,
		//阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		//自动应答
		true,
		//排他
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message:%s", d.Body)
			fmt.Println(d.Body)
		}
	}()

	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
	r.Destory()
}
