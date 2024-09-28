package storage

import "errors"

var (
	ErrNoFaqEntry          = errors.New("no faq entry")
	ErrNoMessage           = errors.New("no message found")
	ErrFeedbackAlreadySent = errors.New("feedback already sent")
)
