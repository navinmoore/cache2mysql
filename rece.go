package main

import "cache2mysql/rabbitMQ"

func main() {
	// db := mysql.MysqlInstance()
	// db.AutoMigrate(&models.User{}, &models.Goods{}, &models.Orders{})

	// service.CreateUser()
	// service.GetUserForMysql()

	// service.CreateGoods()
	// service.GetGoodsForMysql()
	// service.DecrGoodsForRedis(-1)

	// r := rabbitMQ.NewRabbitMQSimple("hello", "amqp://yairs:yairs@127.0.0.1:5672/")
	// r.ConsumeSimple()
	// r := rabbitMQ.NewRabbitMQWorkQueue("hello", "amqp://yairs:yairs@127.0.0.1:5672/")
	// r.Work()
	r := rabbitMQ.NewRoutingRabbitMQ("", "amqp://yairs:yairs@127.0.0.1:5672/")
	r.RoutingConsume()
}
