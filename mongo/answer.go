package mongo


// SaveAnswers 输入答卷
type SaveAnswers struct {
	Id      	string        `json:"id"`
	Answers 	[]interface{} `json:"answers"`
}

// LoadAnswers 输出答卷
type LoadAnswers struct {
	Id      string `json:"id"`
	Answers []Base `json:"answers"`
}

// Base 基础答卷
type Base struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

