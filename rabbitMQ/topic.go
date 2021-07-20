package rabbitMQ

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func NewTopicMQ(mqurl string) *RabbitMQ {
	r := NewRabbitMQ("", "topicExchange", "", mqurl)
	return r
}

func severityForm(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}

// ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table)
func (r *RabbitMQ) TopicPub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "TopicPub Exchange")
	body := BodyForm(os.Args)
	r.channel.Publish(
		r.Exchange,
		severityForm(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	r.failOnErr(err, "TopicPub Publish")
	defer r.Destory()
}

func (r *RabbitMQ) TopicSub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "TopicSub exchange")
	// body := BodyForm(os.Args)
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "TopicSub QueueDeclare")

	for _, s := range os.Args[1:] {
		err = r.channel.QueueBind(
			r.QueueName,
			s,
			r.Exchange,
			false,
			nil,
		)
		r.failOnErr(err, "TopicSub QueueBind")
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever

}
