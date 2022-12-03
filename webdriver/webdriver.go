package webdriver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ikawaha/navigator/webdriver/service"
	"github.com/ikawaha/navigator/webdriver/session"
)

// WebDriver represents a web driver service/client.
type WebDriver struct {
	Timeout    time.Duration
	Debug      bool
	HTTPClient *http.Client
	service    *service.Service
	sessions   []*session.Session
}

// New creates the web driver service/client.
func New(urlT string, commandT []string) *WebDriver {
	return &WebDriver{
		Timeout: session.DefaultWebdriverTimeout,
		Debug:   false,
		HTTPClient: &http.Client{
			Timeout: session.DefaultSessionClientTimeout,
		},
		service:  service.New(urlT, commandT),
		sessions: nil,
	}
}

// URL returns the url of the web driver service.
func (w *WebDriver) URL() string {
	return w.service.URL()
}

// Open returns the session to the web driver service.
func (w *WebDriver) Open(desiredCapabilities map[string]any) (*session.Session, error) {
	url := w.service.URL()
	if url == "" {
		return nil, fmt.Errorf("service not started")
	}
	session, err := session.OpenWithClient(w.HTTPClient, url, desiredCapabilities)
	if err != nil {
		return nil, err
	}
	w.sessions = append(w.sessions, session)
	return session, nil
}

// Start starts the web driver service.
func (w *WebDriver) Start(ctx context.Context) error {
	if err := w.service.Start(ctx, w.Debug); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	if err := w.service.WaitForBoot(w.Timeout); err != nil {
		_ = w.service.Stop()
		return err
	}
	return nil
}

// Stop stops the web driver service.
func (w *WebDriver) Stop() error {
	for _, session := range w.sessions {
		_ = session.Delete()
	}
	if err := w.service.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	return nil
}
