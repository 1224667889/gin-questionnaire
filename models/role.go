package models

type Role struct {
	Id           string `json:"id" type:"serial" constraint:"PRIMARY KEY"`
	Name         string `json:"name" type:"VARCHAR(64)" constraint:"NOT NULL UNIQUE"`
	Introduction string `json:"introduction" type:"VARCHAR(64)" constraint:""`
	extra        string `constraint:""`
}
