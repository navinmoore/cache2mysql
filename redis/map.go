package redis

import (
	"cache2mysql/models"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
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
	if cast.ToInt(preStock) >= Num {
		stock, err := rdb.HIncrBy(ctx, key, "stock", numInt64).Result()
		return stock, err
	} else {
		return 0, errors.New("stock can't decr")
	}

}

func CreateDecrScript() *redis.Script {
	script := redis.NewScript(`
		local num = ARGV[1]
		local Fnum = 0 - num
		local stock = redis.call("hget", KEYS[1], "stock")
		if stock < num  then
			return 0
		end
		
		redis.call("hincrby", KEYS[1], "stock", Fnum)
		return 1
	`)
	return script
}

func EvalDecrScript(goodskey string, num int) {
	rdb := RDBInstance()
	script := CreateDecrScript()
	sha, err := script.Load(rdb.Context(), rdb).Result()
	if err != nil {
		fmt.Println("data.EvalDecrScript err=", err)
	}
	keys := make([]string, 1)
	keys[0] = goodskey
	ret := rdb.EvalSha(rdb.Context(), sha, keys, num)
	if result, err := ret.Result(); err != nil {
		fmt.Printf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Println("result:", result)
	}

}

// func TestScript() *redis.Script {
// 	script := redis.NewScript(`local foo = redis.call("ping");return foo`)
// 	return script
// }

// func TestScriptDo() {
// 	rdb := RDBInstance()
// 	script := TestScript()
// 	sha, err := script.Load(rdb.Context(), rdb).Result()
// 	if err != nil {
// 		fmt.Println("TestScriptDo err=", err)
// 	}
// 	ret := rdb.EvalSha(rdb.Context(), sha, []string{})
// 	if result, err := ret.Result(); err != nil {
// 		fmt.Printf("Execute Redis fail: %v", err.Error())
// 	} else {
// 		fmt.Println("result:", result)
// 	}
// }
