package models

import (
	"cache2mysql/mysql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name" gorm:"type:varchar(50);comment:名字"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex; comment:邮件"`
}

func CreateUser(name, email string) (bool, error) {
	db := mysql.MysqlInstance()
	user := &User{}
	user.Name = name
	user.Email = email
	err := db.Create(user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetUser(ID uint) *User {
	db := mysql.MysqlInstance()
	user := &User{}
	db.Find(user, 1)
	return user
}
