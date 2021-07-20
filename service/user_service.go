package service

import (
	"cache2mysql/models"
	"fmt"
)

func CreateUser() error {
	_, err := models.CreateUser("user1", "user@qq.com")
	return err
}

func GetUserForMysql() models.User {
	user := models.GetUser(1)
	fmt.Printf("%#v", *user)
	return *user
}
