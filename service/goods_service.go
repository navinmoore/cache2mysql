package service

import (
	"cache2mysql/models"
	"cache2mysql/redis"
	"fmt"
)

// Name   string `json:"name" gorm:"type:varchar(50);comment:名称"`
// 	Price  int    `json:"price" gorm:"comment:价格 单位:分"`
// 	Stock  int    `json:"stock" gorm:"comment:库存"`
// 	Code   string `json:"code" gorm:"index:idx_code,unique;comment:编号"`
// 	Status bool   `json:"status" gorm:"default:true;comment:状态 false未上架 true已上架"`

func CreateGoods() {
	goods := models.Goods{}
	goods.Name = "铅笔"
	goods.Stock = 100
	goods.Price = 50
	goods.Code = "pencil1"
	models.CreateGoods(&goods)
}

func GetGoodsForMysql(ID uint) {
	goods, err := models.GetGoodsByID(ID)
	if err != nil {
		fmt.Println("goods_service.GetGoods error=", err)
	}
	fmt.Printf("%#v", goods)

}

func GetGoodsForRedis(ID uint) {
	goods, err := redis.GetGoodsForRedis(ID)
	if err != nil {
		fmt.Println("goods_service.GetGoods error=", err)
	}
	fmt.Printf("%#v", goods)
}

func DecrGoodsForRedis(ID uint, num int) int {
	count, err := redis.DecrGoodsStock(ID, num)
	if err != nil {
		fmt.Println("DecrGoodsForRedis: error=", err)
	}
	fmt.Println("count:", count)
	return int(count)
}

func TestScriptDo() {
	redis.EvalDecrScript("goods_info_1", 2)
}
