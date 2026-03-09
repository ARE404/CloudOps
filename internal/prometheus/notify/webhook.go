package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type WebhookNotifier struct{}

func NewWebhookNotifier() *WebhookNotifier {
	return &WebhookNotifier{}
}

func (n *WebhookNotifier) Send(ctx context.Context, subject, content string, channels []NotifyChannel) error {
	client := &http.Client{Timeout: 10 * time.Second}
	for _, ch := range channels {
		if ch.Type != "webhook" {
			continue
		}
		body := map[string]string{
			"subject": subject,
			"content": content,
		}
		data, _ := json.Marshal(body)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, ch.Target, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}
