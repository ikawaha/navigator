package session

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Connection is a bus to the webdriver service.
type Connection struct {
	sessionURL string
	httpClient *http.Client
	debug      bool
}

func newConnection(ctx context.Context, client *http.Client, serviceURL string, capabilities map[string]any, debug bool) (*Connection, error) {
	req, err := capabilitiesToJSONRequest(capabilities)
	if err != nil {
		return nil, err
	}
	sessionID, err := openSession(ctx, client, serviceURL, req)
	if err != nil {
		return nil, err
	}
	return &Connection{
		sessionURL: serviceURL + "/session/" + sessionID,
		httpClient: client,
		debug:      debug,
	}, nil
}

type desiredCapabilities struct {
	DesiredCapabilities map[string]any `json:"desiredCapabilities"`
}

func capabilitiesToJSONRequest(capabilities map[string]any) (io.Reader, error) {
	if capabilities == nil {
		capabilities = map[string]any{}
	}
	capabilitiesJSON, err := json.Marshal(desiredCapabilities{
		DesiredCapabilities: capabilities,
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(capabilitiesJSON), err
}

func openSession(ctx context.Context, client *http.Client, serviceURL string, body io.Reader) (sessionID string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serviceURL+"/session", body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sessionResponse struct {
		SessionID string
		// fallback for GeckoDriver
		Value struct {
			SessionID string
		}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(b, &sessionResponse); err != nil {
		return "", err
	}

	if sessionResponse.SessionID != "" {
		return sessionResponse.SessionID, nil
	}

	// fallback for GeckoDriver
	if sessionResponse.Value.SessionID != "" {
		return sessionResponse.Value.SessionID, nil
	}
	return "", errors.New("failed to retrieve a session ID")
}

// Send sends the message to the browser.
func (c *Connection) Send(ctx context.Context, method string, pathname string, body, result any) error {
	req, err := bodyToJSON(body)
	if err != nil {
		return err
	}
	path := strings.TrimSuffix(c.sessionURL+"/"+pathname, "/")
	if c.debug {
		log.Printf("%s %s", path, string(req))
	}
	resp, err := c.doRequest(ctx, method, path, req)
	if err != nil {
		return err
	}
	if err := responseToValue(resp, result); err != nil {
		return err
	}
	return nil
}

func bodyToJSON(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	return bodyJSON, nil
}

func responseToValue(src []byte, dst any) error {
	if dst == nil {
		return nil
	}
	v := struct{ Value any }{Value: dst}
	if err := json.Unmarshal(src, &v); err != nil {
		return fmt.Errorf("unexpected response: %s", src)
	}
	return nil
}

func (c *Connection) doRequest(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, toResponseError(b)
	}
	return b, nil
}

func toResponseError(body []byte) error {
	var errBody struct{ Value struct{ Message string } }
	if err := json.Unmarshal(body, &errBody); err != nil {
		return fmt.Errorf("request unsuccessful: %s", body)
	}

	var errMessage struct{ ErrorMessage string }
	if err := json.Unmarshal([]byte(errBody.Value.Message), &errMessage); err != nil {
		return fmt.Errorf("request unsuccessful: %s", errBody.Value.Message)
	}

	return fmt.Errorf("request unsuccessful: %s", errMessage.ErrorMessage)
}
