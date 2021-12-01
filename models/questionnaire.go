package models

type Questionnaire struct {
	Id          string `json:"id" type:"VARCHAR(64)" constraint:"PRIMARY KEY"`
	Title       string `json:"title" type:"VARCHAR(64)" constraint:"NOT NULL DEFAULT '未命名'"`
	Subtitle    string `json:"subtitle" type:"VARCHAR(64)" constraint:""`
	Description string `json:"description" type:"VARCHAR(255)" constraint:""`
	Deadline    string `json:"deadline" type:"DATE" constraint:"NOT NULL"`
	NeedLogin   string `json:"need_login" type:"BOOLEAN" constraint:"DEFAULT false"`
	IsOpen      string `json:"is_open" type:"BOOLEAN" constraint:"DEFAULT false"`
	HasReleased string `json:"has_released" type:"BOOLEAN" constraint:"DEFAULT false"`
	AccountId   string `json:"account_id" type:"VARCHAR(64)" constraint:"NOT NULL"`
	extra       string `constraint:"FOREIGN KEY(account_id) REFERENCES Account(id)"`
}

