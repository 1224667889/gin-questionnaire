package models

type File struct {
	Id         string `json:"id" type:"VARCHAR(64)" constraint:"PRIMARY KEY"`
	Title      string `json:"title" type:"VARCHAR(64)" constraint:"NOT NULL"`
	QuestionId string `json:"question_id" type:"VARCHAR(64)" constraint:"NOT NULL"`
	extra      string `constraint:"FOREIGN KEY(question_id) REFERENCES Question(id)"`
}
