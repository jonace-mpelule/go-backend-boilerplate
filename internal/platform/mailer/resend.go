package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ResendMailer struct {
	apiKey  string
	baseURL string
	from    string
	client  *http.Client
}

func NewResend(apiKey, baseURL, from string) Mailer {
	return &ResendMailer{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		from:    from,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *ResendMailer) Send(ctx context.Context, to, subject, html string) error {
	payload := map[string]any{
		"from":    r.from,
		"to":      []string{to},
		"subject": subject,
		"html":    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend send failed with status %d", resp.StatusCode)
	}

	return nil
}
