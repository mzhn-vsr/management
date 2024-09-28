package dto

type Feedback struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	IsUseful bool   `json:"isUseful"`
}

type FeedbackStats struct {
	Positive int `json:"positive"`
	Negative int `json:"negative"`
	Total    int `json:"total"`
}
