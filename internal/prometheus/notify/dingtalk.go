package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DingTalkNotifier struct{}

func NewDingTalkNotifier() *DingTalkNotifier {
	return &DingTalkNotifier{}
}

func (n *DingTalkNotifier) Send(ctx context.Context, subject, content string, channels []NotifyChannel) error {
	client := &http.Client{Timeout: 10 * time.Second}
	for _, ch := range channels {
		if ch.Type != "dingtalk" {
			continue
		}
		body := map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": subject,
				"text":  fmt.Sprintf("## %s\n\n%s", subject, content),
			},
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
