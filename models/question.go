package models

type Question struct {
	Id              string `json:"id" type:"VARCHAR(64)" constraint:"PRIMARY KEY"`
	QuestionType    string `json:"question_type" type:"VARCHAR(64)" constraint:"NOT NULL"`
	Title           string `json:"title" type:"VARCHAR(64)" constraint:"NOT NULL"`
	Description     string `json:"description" type:"VARCHAR(64)" constraint:"NOT NULL"`
	IsRequired      string `json:"is_required" type:"VARCHAR(64)" constraint:"NOT NULL"`
	QuestionnaireId string `json:"questionnaire_id" type:"VARCHAR(64)" constraint:"NOT NULL"`
	extra           string `constraint:"FOREIGN KEY(questionnaire_id) REFERENCES Questionnaire(id) ON DELETE cascade"`
}
