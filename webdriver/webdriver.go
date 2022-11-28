package webdriver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ikawaha/navigator/webdriver/service"
	"github.com/ikawaha/navigator/webdriver/session"
)

type WebDriver struct {
	Timeout    time.Duration
	Debug      bool
	HTTPClient *http.Client
	service    *service.Service
	sessions   []*session.Session
}

func New(urlT string, commandT []string) *WebDriver {
	return &WebDriver{
		Timeout: session.DefaultWebdriverTimeout,
		Debug:   false,
		HTTPClient: &http.Client{
			Timeout: session.DefaultSessionClientTimeout,
		},
		service: &service.Service{
			URLT:     urlT,
			CommandT: commandT,
		},
		sessions: nil,
	}
}

func (w *WebDriver) URL() string {
	return w.service.URL()
}

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

func (w *WebDriver) Start() error {
	if err := w.service.Start(w.Debug); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}
	if err := w.service.WaitForBoot(w.Timeout); err != nil {
		w.service.Stop()
		return err
	}
	return nil
}

func (w *WebDriver) Stop() error {
	for _, session := range w.sessions {
		session.Delete()
	}
	if err := w.service.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %s", err)
	}
	return nil
}
