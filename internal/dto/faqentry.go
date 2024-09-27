package dto

type FaqEntryCreate struct {
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
	Classifier1 *string `json:"classifier1,omitempty"`
	Classifier2 *string `json:"classifier2,omitempty"`
}

type FaqEntryUpdate struct {
	Id          string  `json:"id"`
	Question    *string `json:"question"`
	Answer      *string `json:"answer"`
	Classifier1 *string `json:"classifier1"`
	Classifier2 *string `json:"classifier2"`
}

type FaqEntryList struct {
	Pagination
}
