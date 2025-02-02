package receivers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Awtrix struct {
	host   string
	client Client
}

func NewAwtrix(host string, client Client) (*Awtrix, error) {
	if client == nil {
		client = http.DefaultClient
	}

	return &Awtrix{
		host:   host,
		client: client,
	}, nil
}

func (a *Awtrix) SendReport(ctx context.Context, subject, _ string) error {
	subject = replaceSeverityLevels(subject)
	request, err := a.getRequest(ctx, subject)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("awtrix sent status code %d", resp.StatusCode)
	}

	return nil
}

func (a *Awtrix) getRequest(ctx context.Context, text string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/notify", a.host)
	data := map[string]any{
		"text":     text,
		"duration": 15,
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(marshalled))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req, nil
}

var (
	regex *regexp.Regexp
	once  sync.Once
)

func replaceSeverityLevels(text string) string {
	replacements := map[string]string{
		"warning":   "warn",
		"critical":  "crit",
		"error":     "err",
		"info":      "info",
		"debug":     "dbg",
		"notice":    "note",
		"emergency": "emerg",
	}

	once.Do(func() {
		var keys []string
		for key := range replacements {
			keys = append(keys, key)
		}
		pattern := `(?i)\b(` + strings.Join(keys, "|") + `)\b`
		regex = regexp.MustCompile(pattern)
	})

	result := regex.ReplaceAllStringFunc(text, func(match string) string {
		return replacements[strings.ToLower(match)]
	})

	return result
}
