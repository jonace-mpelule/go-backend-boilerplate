package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/username/project-name/internal/config"
)

const lokiPushPath = "/loki/api/v1/push"

type lokiWriteSyncer struct {
	client   *http.Client
	endpoint string
	tenantID string
	username string
	password string
	labels   map[string]string

	mu      sync.Mutex
	entries []lokiEntry

	ticker *time.Ticker
	done   chan struct{}
	once   sync.Once
}

type lokiEntry struct {
	Timestamp time.Time
	Line      string
}

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func newLokiWriteSyncer(appCfg config.AppConfig, lokiCfg config.LokiConfig) (*lokiWriteSyncer, error) {
	if !lokiCfg.Enabled {
		return nil, nil
	}

	endpoint := strings.TrimRight(lokiCfg.URL, "/")
	if !strings.HasSuffix(endpoint, lokiPushPath) {
		endpoint += lokiPushPath
	}

	hostname, _ := os.Hostname()
	labels := map[string]string{
		"app":      appCfg.Name,
		"env":      appCfg.Env,
		"service":  appCfg.Name,
		"instance": hostname,
	}

	for key, value := range lokiCfg.Labels {
		labels[key] = value
	}

	syncer := &lokiWriteSyncer{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoint: endpoint,
		tenantID: lokiCfg.TenantID,
		username: lokiCfg.BasicAuthUsername,
		password: lokiCfg.BasicAuthPassword,
		labels:   labels,
		ticker:   time.NewTicker(lokiCfg.BatchWait),
		done:     make(chan struct{}),
	}

	go syncer.run()

	return syncer, nil
}

func (l *lokiWriteSyncer) Write(p []byte) (int, error) {
	line := strings.TrimSpace(string(p))
	if line == "" {
		return len(p), nil
	}

	l.mu.Lock()
	l.entries = append(l.entries, lokiEntry{
		Timestamp: time.Now().UTC(),
		Line:      line,
	})
	l.mu.Unlock()

	return len(p), nil
}

func (l *lokiWriteSyncer) Sync() error {
	return l.flush()
}

func (l *lokiWriteSyncer) Close() error {
	var closeErr error
	l.once.Do(func() {
		close(l.done)
		l.ticker.Stop()
		closeErr = l.flush()
	})
	return closeErr
}

func (l *lokiWriteSyncer) run() {
	for {
		select {
		case <-l.done:
			return
		case <-l.ticker.C:
			if err := l.flush(); err != nil {
				fmt.Fprintf(os.Stderr, "loki flush error: %v\n", err)
			}
		}
	}
}

func (l *lokiWriteSyncer) flush() error {
	l.mu.Lock()
	if len(l.entries) == 0 {
		l.mu.Unlock()
		return nil
	}

	entries := make([]lokiEntry, len(l.entries))
	copy(entries, l.entries)
	l.entries = nil
	l.mu.Unlock()

	values := make([][]string, 0, len(entries))
	for _, entry := range entries {
		values = append(values, []string{
			fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
			entry.Line,
		})
	}

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: l.labels,
				Values: values,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, l.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if l.tenantID != "" {
		req.Header.Set("X-Scope-OrgID", l.tenantID)
	}
	if l.username != "" || l.password != "" {
		req.SetBasicAuth(l.username, l.password)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki push returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}
