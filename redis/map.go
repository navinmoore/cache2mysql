package redis

import (
	"cache2mysql/models"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cast"
)

// 将goods读入redis 判断是否存在在redis,不存在读取mysql并缓存到redis
var (
	GoodsKey = "goods_info_%d"
)

func GetGoodsForRedis(ID uint) (models.Goods, error) {
	rdb := RDBInstance()
	// fmt.Println("GetGoodsForRedis", ID)
	key := fmt.Sprintf(GoodsKey, ID)
	ctx := context.Background()
	val, _ := rdb.HMGet(ctx, key, "id", "name", "price", "stock", "code").Result()
	// fmt.Println("##", val)
	if val[0] != nil {
		goods := models.Goods{}
		goods.ID = cast.ToUint(val[0])
		goods.Name = cast.ToString(val[1])
		goods.Price = cast.ToInt(val[2])
		goods.Stock = cast.ToInt(val[3])
		goods.Code = cast.ToString(val[4])

		// fmt.Println("GetGoodsForRedis val!=nil")
		return goods, nil

	} else {
		goods, err := SetGoodsToRedis(ID)
		// fmt.Println("GetGoodsForRedis val==nil", err)
		return goods, err
	}

}

func SetGoodsToRedis(ID uint) (models.Goods, error) {
	goods, err := models.GetGoodsByID(ID)
	if err != nil {
		return models.Goods{}, err
	}

	rdb := RDBInstance()
	// m := structs.Map(goods)
	m := make(map[string]interface{})
	m["id"] = goods.ID
	m["name"] = goods.Name
	m["price"] = goods.Price
	m["stock"] = goods.Stock
	m["code"] = goods.Code

	// fmt.Println("##########", m)
	key := fmt.Sprintf(GoodsKey, ID)
	ctx := context.Background()
	err = rdb.HMSet(ctx, key, m).Err()
	if err != nil {
		return models.Goods{}, err
	}
	return goods, nil
}

func DecrGoodsStock(ID uint, Num int) (int64, error) {
	numInt64 := int64(Num)
	rdb := RDBInstance()
	key := fmt.Sprintf(GoodsKey, ID)
	ctx := context.Background()
	preStock, _ := rdb.HGet(ctx, key, "stock").Result()
	// fmt.Println("#######", preStock)
	if cast.ToInt(preStock) >= Num {
		stock, err := rdb.HIncrBy(ctx, key, "stock", numInt64).Result()
		return stock, err
	} else {
		return 0, errors.New("stock can't decr")
	}

}
