package entity

type ChatInvokeResponse struct {
	Output struct {
		Content string `json:"content"`
	} `json:"output"`
}

type ChatInvokeAnswer struct {
	Id     string
	Answer string `json:"answer"`
	Class1 string `json:"class_1"`
	Class2 string `json:"class_2"`
}

type ChatInvokeOutput struct {
	Answer string `json:"answer"`
}
