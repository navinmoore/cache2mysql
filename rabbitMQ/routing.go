package rabbitMQ

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func severityForm(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}

func NewRoutingRabbitMQ(queueName, mqurl string) *RabbitMQ {
	r := NewRabbitMQ("", "routingExchange", "", mqurl)
	return r
}

func (r *RabbitMQ) RoutingPub() {
	//ExchangeDeclare 定义exchane
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "RoutingPub exchangeDeclare")
	body := BodyForm(os.Args)
	// Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error
	err = r.channel.Publish(
		r.Exchange,
		severityForm(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	r.failOnErr(err, "RoutingPub publish")
	defer r.Destory()
}

func (r *RabbitMQ) RoutingConsume() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "RoutingConsume Exchange")
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		true, //exclusive
		false,
		nil,
	)
	r.failOnErr(err, "RoutingConsume QueueDeclare")

	if len(os.Args) < 2 {
		r.failOnErr(nil, "RoutingConsume args")
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		err = r.channel.QueueBind(
			r.QueueName,
			s,
			r.Exchange,
			false,
			nil,
		)
		r.failOnErr(err, "RoutingConsume QueueBind")
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

	r.failOnErr(err, "RoutingConsume Consume")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever

}
