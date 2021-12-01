package service

import (
	"log"
	"main/models"
)

func InitData() {
	if err := models.Insert(
		&models.Role{
			Id:           "1",
			Name:         "SuperAdmin",
			Introduction: "超级管理员",
		}); err != nil {log.Fatal(err)}
	if err := models.Insert(
		&models.Role{
			Id:           "2",
			Name:         "Admin",
			Introduction: "管理员",
		}); err != nil {log.Fatal(err)}
	if err := models.Insert(
		&models.Role{
			Id:           "3",
			Name:         "User",
			Introduction: "用户",
		}); err != nil {log.Fatal(err)}
	if err := models.Insert(
		&models.Account{
			Name:     "cj",
			Account:  "cj666",
			Email:    "1224667889@qq.com",
			Password: "aaaaaa",
			RoleId:   "1",
		}); err != nil {log.Fatal(err)}
}