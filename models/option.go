package models

type Option struct {
	Id         string `json:"id" type:"VARCHAR(64)" constraint:"PRIMARY KEY"`
	Content    string `json:"content" type:"VARCHAR(255)" constraint:"NOT NULL"`
	QuestionId string `json:"question_id" type:"VARCHAR(64)" constraint:"NOT NULL"`
	extra      string `constraint:"FOREIGN KEY(question_id) REFERENCES Question(id) ON DELETE cascade"`
}
