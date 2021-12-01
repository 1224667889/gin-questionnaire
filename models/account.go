package models

type Account struct {
	Id       string `json:"id" type:"serial" constraint:"PRIMARY KEY"`
	Name     string `json:"name" type:"VARCHAR(64)" constraint:"NOT NULL"`
	Account  string `json:"account" type:"VARCHAR(64)" constraint:"NOT NULL UNIQUE"`
	Email    string `json:"email" type:"VARCHAR(64)" constraint:"NOT NULL UNIQUE"`
	Password string `json:"password" type:"VARCHAR(64)" constraint:"NOT NULL"`
	RoleId   string `json:"role_id" type:"INT" constraint:"NOT NULL"`
	extra    string `constraint:"FOREIGN KEY(role_id) REFERENCES Role(id) ON DELETE cascade"`
}



