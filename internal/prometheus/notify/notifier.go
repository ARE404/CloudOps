package notify

import "context"

// NotifyChannel represents a notification channel type
type NotifyChannel struct {
	Type   string `json:"type"`   // dingtalk | email | webhook
	Target string `json:"target"` // webhook URL or email address
}

// Notifier sends alert notifications
type Notifier interface {
	Send(ctx context.Context, subject, content string, channels []NotifyChannel) error
}
