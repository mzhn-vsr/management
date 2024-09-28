package entity

type ClassifierResponse struct {
	Output struct {
		Class1 string `json:"c1"`
		Class2 string `json:"c2"`
	} `json:"output"`
	Metadata struct {
		RunId          string   `json:"run_id"`
		FeedbackTokens []string `json:"feedback_tokens"`
	} `json:"metadata"`
}
