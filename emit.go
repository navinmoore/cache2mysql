package main

import (
	"cache2mysql/service"
)

func main() {
	// db := mysql.MysqlInstance()
	// db.AutoMigrate(&models.User{}, &models.Goods{}, &models.Orders{})

	// service.CreateUser()
	// service.GetUserForMysql()

	// service.CreateGoods()
	// service.GetGoodsForMysql()
	// service.DecrGoodsForRedis(-1)

	// r := rabbitMQ.NewRabbitMQSimple("hello", "amqp://yairs:yairs@127.0.0.1:5672/")
	// r.PublishSimple("hello")
	// r := rabbitMQ.NewRabbitMQWorkQueue("hello", "amqp://yairs:yairs@127.0.0.1:5672/")
	// r.Task()
	// r := rabbitMQ.NewRoutingRabbitMQ("", "amqp://yairs:yairs@127.0.0.1:5672/")
	// r.RoutingPub()
	service.TestScriptDo()
}
