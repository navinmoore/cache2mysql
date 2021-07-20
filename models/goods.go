package models

import (
	"cache2mysql/mysql"
)

// 商品
type Goods struct {
	BaseModel
	Name   string `json:"name" gorm:"type:varchar(50);comment:名称"`
	Price  int    `json:"price" gorm:"comment:价格 单位:分"`
	Stock  int    `json:"stock" gorm:"comment:库存"`
	Code   string `json:"code" gorm:"index:idx_code,unique;comment:编号"`
	Status bool   `json:"status" gorm:"default:true;comment:状态 false未上架 true已上架"`
}

type Orders struct {
	BaseModel
	UserID  uint `json:"user_id"`
	GoodsID uint `json:"goods_id"`
	Num     int  `json:"num" gorm:"comment:个数"`
	Money   int  `json:"money" gorm:"comment:总价 分"`
	Status  int  `json:"status" gorm:"type:tinyint;default:0;comment:状态 0未支付 1已支付"`
}

func CreateGoods(goods *Goods) *Goods {
	db := mysql.MysqlInstance()
	db.Create(goods)
	return goods
}

func GetGoodsByID(ID uint) (Goods, error) {
	var goods Goods
	db := mysql.MysqlInstance()
	if err := db.Find(&goods, ID).Error; err != nil {
		return goods, err
	}
	if !goods.Status {
		return Goods{}, nil
	}
	return goods, nil
}

//事务 更新库存 incr为false 是减库存  true是增库存
func IncrementGoodsStock(ID uint, num int, incr bool) (rows int, err error) {
	rows = 0
	db := mysql.MysqlInstance()
	// goods := Goods{}
	tx := db.Begin()
	// tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&goods, ID)
	if incr {
		if err = tx.Exec("update goods set stock = stock + ? where id = ?", num, ID).Error; err != nil {
			tx.Rollback()
			return
		}
	} else {
		if err = tx.Exec("update goods set stock = stock - ? where id = ? and stock >= ?", num, ID, num).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		return
	}
	rows = int(tx.RowsAffected)
	return
}

// 更新商品信息，不包括库存
func UpdateGoods(ID uint, m map[string]interface{}) (goods *Goods, err error) {
	goods.ID = ID
	db := mysql.MysqlInstance()
	tx := db.Begin()
	// 更新字段 除了 stock，id （包括零值字段的所有字段）eg:就算code更新后为0也更新
	// tx.Model(goods).Omit("id", "stock").Updates(m).Error
	if err = tx.Model(goods).Select("name", "price", "code", "status").Updates(m).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		return
	}
	return
}

// false创建失败  true 创建成功
func CreateOrders(userID uint, goodsID uint, num int) (bool, error) {
	db := mysql.MysqlInstance()
	goods := &Goods{}
	orders := &Orders{}
	tx := db.Begin()
	err := tx.Model(goods).Find(goods, goodsID).Error
	if err != nil {
		return false, err
	}
	orders.GoodsID = goodsID
	orders.Money = goods.Price * num
	orders.Num = num
	orders.UserID = userID
	err = tx.Create(orders).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if err = tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil
}

//删除订单
func DeleteOrder(ID uint) error {
	db := mysql.MysqlInstance()
	err := db.Delete(&Orders{}, ID).Error
	return err
}

// 支付订单
func PayOrder(ID uint) (bool, error) {
	order := &Orders{}
	db := mysql.MysqlInstance()
	tx := db.Begin()
	err := tx.Model(order).Where("status = ?", 0).Update("status", 1).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if err = tx.Commit().Error; err != nil {
		return false, err

	}
	return true, nil
}
