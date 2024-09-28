package entity

import "time"

type Feedback struct {
	Id        string    `db:"id"`
	Question  string    `db:"question"`
	Answer    string    `db:"answer"`
	IsUseful  *bool     `db:"is_useful"`
	CreatedAt time.Time `db:"created_at"`
}
