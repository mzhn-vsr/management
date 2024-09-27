package entity

import "time"

type FaqEntry struct {
	Id          string     `db:"id" json:"id"`
	Question    string     `db:"question" json:"question"`
	Answer      string     `db:"answer" json:"answer"`
	Classifier1 *string    `db:"classifier1" json:"classifier1"`
	Classifier2 *string    `db:"classifier2" json:"classifier2"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
}
