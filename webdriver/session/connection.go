package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Connection struct {
	SessionURL string
	httpClient *http.Client
}

func newConnection(client *http.Client, serviceURL string, capabilities map[string]any) (*Connection, error) {
	req, err := capabilitiesToJSONRequest(capabilities)
	if err != nil {
		return nil, err
	}
	sessionID, err := openSession(client, serviceURL, req)
	if err != nil {
		return nil, err
	}
	return &Connection{
		SessionURL: serviceURL + "/session/" + sessionID,
		httpClient: client,
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

func openSession(client *http.Client, serviceURL string, body io.Reader) (sessionID string, err error) {
	req, err := http.NewRequest(http.MethodPost, serviceURL+"/session", body)
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

func (c *Connection) Send(method string, pathname string, body, result any) error {
	req, err := bodyToJSON(body)
	if err != nil {
		return err
	}
	path := strings.TrimSuffix(c.SessionURL+"/"+pathname, "/")

	log.Println(path) //XXX

	resp, err := c.doRequest(method, path, req)
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

func (c *Connection) doRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, toResponseError(resp)
	}

	return resp, nil
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
