package feedbackservice

import "errors"

var (
	ErrNoMessage           = errors.New("no message found")
	ErrFeedbackAlreadySent = errors.New("feedback already sent")
)
