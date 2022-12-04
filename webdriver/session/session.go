package session

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ikawaha/navigator/event"
)

const (
	// DefaultWebdriverTimeout is the waiting time limit for the web driver service to start.
	DefaultWebdriverTimeout = 10 * time.Second

	// DefaultSessionClientTimeout is the default time limit for requests
	// to the web driver service of the session client.
	DefaultSessionClientTimeout = 30 * time.Second
)

// Method is the (HTTP) method to access to the web driver service.
type Method = string

const (
	// Get is the GET method.
	Get Method = http.MethodGet

	// Post is the POST method.
	Post Method = http.MethodPost

	// Delete is the DELETE method.
	Delete Method = http.MethodDelete
)

// Session represents a session to the web driver service.
type Session struct {
	*Connection
}

// OpenWithClient returns a session to the web driver service.
func OpenWithClient(ctx context.Context, client *http.Client, url string, capabilities map[string]any, debug bool) (*Session, error) {
	c, err := newConnection(ctx, client, url, capabilities, debug)
	if err != nil {
		return nil, err
	}
	return &Session{Connection: c}, nil
}

// Delete sends to delete message to terminate the session.
func (s *Session) Delete(ctx context.Context) error {
	return s.Send(ctx, Delete, "", nil, nil)
}

// Selector represents a selector for elements.
type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

// GetElement retrieves the element matching the selector query of the session.
func (s *Session) GetElement(ctx context.Context, selector Selector) (*Element, error) {
	var result elementResult
	if err := s.Send(ctx, Post, "element", selector, &result); err != nil {
		return nil, err
	}
	return &Element{
		ID:      result.ID(),
		Session: s,
	}, nil
}

// GetElements retrieves elements matching the selector query of the session.
func (s *Session) GetElements(ctx context.Context, selector Selector) ([]*Element, error) {
	var results []elementResult
	if err := s.Send(ctx, Post, "elements", selector, &results); err != nil {
		return nil, err
	}
	var elements []*Element
	for _, result := range results {
		elements = append(elements, &Element{ID: result.ID(), Session: s})
	}
	return elements, nil
}

// GetActiveElement returns the active element of the session.
func (s *Session) GetActiveElement(ctx context.Context) (*Element, error) {
	var result elementResult
	if err := s.Send(ctx, Post, "element/active", nil, &result); err != nil {
		return nil, err
	}
	return &Element{ID: result.ID(), Session: s}, nil
}

// GetWindow returns the window handler of the session.
func (s *Session) GetWindow(ctx context.Context) (*Window, error) {
	var windowID string
	if err := s.Send(ctx, Get, "window_handle", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{ID: windowID, Session: s}, nil
}

// GetWindows returns window handlers of the session.
func (s *Session) GetWindows(ctx context.Context) ([]*Window, error) {
	var windowsID []string
	if err := s.Send(ctx, Get, "window_handles", nil, &windowsID); err != nil {
		return nil, err
	}

	var windows []*Window
	for _, windowID := range windowsID {
		windows = append(windows, &Window{windowID, s})
	}
	return windows, nil
}

type nameRequest struct {
	Name string `json:"name"`
}

// SetWindow sets the window to the browser.
func (s *Session) SetWindow(ctx context.Context, window *Window) error {
	if window == nil {
		return errors.New("nil window is invalid")
	}
	return s.Send(ctx, Post, "window", nameRequest{
		Name: window.ID,
	}, nil)
}

// SetWindowByName sets the window to the browser by name.
func (s *Session) SetWindowByName(ctx context.Context, name string) error {
	return s.Send(ctx, Post, "window", nameRequest{
		Name: name,
	}, nil)
}

// DeleteWindow deletes the window of the session.
func (s *Session) DeleteWindow(ctx context.Context) error {
	if err := s.Send(ctx, Delete, "window", nil, nil); err != nil {
		return err
	}
	return nil
}

// GetCookies gets cookies of the session.
func (s *Session) GetCookies(ctx context.Context) ([]*Cookie, error) {
	var cookies []*Cookie
	if err := s.Send(ctx, Get, "cookie", nil, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

type cookieRequest struct {
	Cookie *Cookie `json:"cookie"`
}

// SetCookie sets a cookie to the browser.
func (s *Session) SetCookie(ctx context.Context, cookie *Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}
	return s.Send(ctx, Post, "cookie", cookieRequest{
		Cookie: cookie,
	}, nil)
}

// DeleteCookie deletes a cookie of the session.
func (s *Session) DeleteCookie(ctx context.Context, cookieName string) error {
	return s.Send(ctx, Delete, path.Join("cookie", cookieName), nil, nil)
}

// DeleteCookies deletes cookies of the session.
func (s *Session) DeleteCookies(ctx context.Context) error {
	return s.Send(ctx, Delete, "cookie", nil, nil)
}

// GetScreenshot gets a screenshot.
func (s *Session) GetScreenshot(ctx context.Context) ([]byte, error) {
	var base64Image string
	if err := s.Send(ctx, Get, "screenshot", nil, &base64Image); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(base64Image)
}

// GetURL gets the url of the session.
func (s *Session) GetURL(ctx context.Context) (string, error) {
	var url string
	if err := s.Send(ctx, Get, "url", nil, &url); err != nil {
		return "", err
	}
	return url, nil
}

type urlRequest struct {
	URL string `json:"url"`
}

// SetURL sets the url to the browser.
func (s *Session) SetURL(ctx context.Context, url string) error {
	return s.Send(ctx, Post, "url", urlRequest{
		URL: url,
	}, nil)
}

// GetTitle gets the title of the session.
func (s *Session) GetTitle(ctx context.Context) (string, error) {
	var title string
	if err := s.Send(ctx, Get, "title", nil, &title); err != nil {
		return "", err
	}
	return title, nil
}

// GetSource gets the page source of the session.
func (s *Session) GetSource(ctx context.Context) (string, error) {
	var source string
	if err := s.Send(ctx, Get, "source", nil, &source); err != nil {
		return "", err
	}
	return source, nil
}

// MoveTo moves the element to the offset position.
func (s *Session) MoveTo(ctx context.Context, region *Element, offset Offset) error {
	req := map[string]any{}
	if region != nil {
		req["element"] = region.ID
	}
	if offset != nil {
		if xoffset, present := offset.x(); present {
			req["xoffset"] = xoffset
		}
		if yoffset, present := offset.y(); present {
			req["yoffset"] = yoffset
		}
	}
	return s.Send(ctx, Post, "moveto", req, nil)
}

// Frame sets the frame to the browser.
func (s *Session) Frame(ctx context.Context, frame *Element) error {
	var elementID any
	if frame != nil {
		elementID = struct {
			Element string `json:"ELEMENT"`
		}{
			Element: frame.ID,
		}
	}
	req := struct {
		ID any `json:"id"`
	}{
		ID: elementID,
	}
	return s.Send(ctx, Post, "frame", req, nil)
}

// FrameParent sets the parent frame to the browser.
func (s *Session) FrameParent(ctx context.Context) error {
	return s.Send(ctx, Post, "frame/parent", nil, nil)
}

type scriptRequest struct {
	Script string `json:"script"`
	Args   []any  `json:"args"`
}

// Execute executes the script.
func (s *Session) Execute(ctx context.Context, body string, arguments []any, result any) error {
	if arguments == nil {
		arguments = []any{}
	}
	return s.Send(ctx, Post, "execute", scriptRequest{
		Script: body,
		Args:   arguments,
	}, result)
}

// Forward forwards the browser.
func (s *Session) Forward(ctx context.Context) error {
	return s.Send(ctx, Post, "forward", nil, nil)
}

// Back backs the browser.
func (s *Session) Back(ctx context.Context) error {
	return s.Send(ctx, Post, "back", nil, nil)
}

// Refresh refreshes the browser.
func (s *Session) Refresh(ctx context.Context) error {
	return s.Send(ctx, Post, "refresh", nil, nil)
}

// GetAlertText gets the alert text of the browser.
func (s *Session) GetAlertText(ctx context.Context) (string, error) {
	var text string
	if err := s.Send(ctx, Get, "alert_text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

type textRequest struct {
	Text string `json:"text"`
}

// SetAlertText sets the text to the browser.
func (s *Session) SetAlertText(ctx context.Context, text string) error {
	return s.Send(ctx, Post, "alert_text", textRequest{
		Text: text,
	}, nil)
}

// AcceptAlert accepts the alert of the browser.
func (s *Session) AcceptAlert(ctx context.Context) error {
	return s.Send(ctx, Post, "accept_alert", nil, nil)
}

// DismissAlert dismisses the alert of the browser.
func (s *Session) DismissAlert(ctx context.Context) error {
	return s.Send(ctx, Post, "dismiss_alert", nil, nil)
}

type typeRequest struct {
	Type string `json:"type"`
}

// NewLogs gets logs of the browser.
func (s *Session) NewLogs(ctx context.Context, logType string) ([]Log, error) {
	var logs []Log
	if err := s.Send(ctx, Post, "log", typeRequest{
		Type: logType,
	}, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLogTypes gets log types.
func (s *Session) GetLogTypes(ctx context.Context) ([]string, error) {
	var types []string
	if err := s.Send(ctx, Get, "log/types", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}

type buttonRequest struct {
	Button event.Button `json:"button"`
}

// DoubleClick sends the double click event to the browser.
func (s *Session) DoubleClick(ctx context.Context) error {
	return s.Send(ctx, Post, "doubleclick", nil, nil)
}

// Click sends the click event to the browser.
func (s *Session) Click(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "click", buttonRequest{Button: button}, nil)
}

// ButtonDown sends the button down event to the browser.
func (s *Session) ButtonDown(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "buttondown", buttonRequest{Button: button}, nil)
}

// ButtonUp sends the button up event to the browser.
func (s *Session) ButtonUp(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "buttonup", buttonRequest{Button: button}, nil)
}

type xyRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// TouchDown sends the touch-down event to the browser.
func (s *Session) TouchDown(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/down", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

// TouchUp sends the touch-up event to the browser.
func (s *Session) TouchUp(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/up", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

// TouchMove sends the touch-move event to the browser.
func (s *Session) TouchMove(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/move", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

type elementRequest struct {
	Element string `json:"element"`
}

// TouchClick sends touch-click event to the browser.
func (s *Session) TouchClick(ctx context.Context, element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(ctx, Post, "touch/click", elementRequest{
		Element: element.ID,
	}, nil)
}

// TouchDoubleClick sends the touch-double-click event to the browser.
func (s *Session) TouchDoubleClick(ctx context.Context, element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(ctx, Post, "touch/doubleclick", elementRequest{
		Element: element.ID,
	}, nil)
}

// TouchLongClick sends the touch-long-click event to the browser.
func (s *Session) TouchLongClick(ctx context.Context, element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(ctx, Post, "touch/longclick", elementRequest{
		Element: element.ID,
	}, nil)
}

type xySpeedRequest struct {
	XSpeed int `json:"xspeed"`
	YSpeed int `json:"yspeed"`
}

type touchFlickRequest struct {
	Element string `json:"element"`
	XOffset int    `json:"xoffset"`
	YOffset int    `json:"yoffset"`
	Speed   uint   `json:"speed"`
}

// TouchFlick sends the touch-flick event to the browser.
func (s *Session) TouchFlick(ctx context.Context, element *Element, offset Offset, speed Speed) error {
	if speed == nil {
		return errors.New("nil speed is invalid")
	}
	if (element == nil) != (offset == nil) {
		return errors.New("element must be provided if offset is provided and vice versa")
	}
	if element == nil {
		xSpeed, ySpeed := speed.vector()
		return s.Send(ctx, Post, "touch/flick", xySpeedRequest{
			XSpeed: xSpeed,
			YSpeed: ySpeed,
		}, nil)
	}
	xOffset, yOffset := offset.position()
	return s.Send(ctx, Post, "touch/flick", touchFlickRequest{
		Element: element.ID,
		XOffset: xOffset,
		YOffset: yOffset,
		Speed:   speed.scalar(),
	}, nil)
}

type touchScrollRequest struct {
	Element string `json:"element,omitempty"`
	XOffset int    `json:"xoffset"`
	YOffset int    `json:"yoffset"`
}

// TouchScroll sends the touch-scroll event to the browser.
func (s *Session) TouchScroll(ctx context.Context, element *Element, offset Offset) error {
	if element == nil {
		element = &Element{}
	}
	if offset == nil {
		return errors.New("nil offset is invalid")
	}
	xOffset, yOffset := offset.position()
	return s.Send(ctx, Post, "touch/scroll", touchScrollRequest{
		Element: element.ID,
		XOffset: xOffset,
		YOffset: yOffset,
	}, nil)
}

type valueSliceRequest struct {
	Value []string `json:"value"`
}

// Keys sends key events of the text to the browser.
func (s *Session) Keys(ctx context.Context, text string) error {
	return s.Send(ctx, Post, "keys", valueSliceRequest{
		Value: strings.Split(text, ""),
	}, nil)
}

// DeleteLocalStorage deletes the local storage of the browser.
func (s *Session) DeleteLocalStorage(ctx context.Context) error {
	return s.Send(ctx, Delete, "local_storage", nil, nil)
}

// DeleteSessionStorage deletes the session storage of the browser.
func (s *Session) DeleteSessionStorage(ctx context.Context) error {
	return s.Send(ctx, Delete, "session_storage", nil, nil)
}

type msRequest struct {
	MS   int    `json:"ms"`
	Type string `json:"type,omitempty"`
}

// SetImplicitWait sets the implicit wait to the browser.
func (s *Session) SetImplicitWait(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts/implicit_wait", msRequest{
		MS: timeout,
	}, nil)
}

// SetPageLoad sets the timeout to the page load of the browser.
func (s *Session) SetPageLoad(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts", msRequest{
		MS:   timeout,
		Type: "page load",
	}, nil)
}

// SetScriptTimeout sets the timeout to the asynchronous scripts execution.
func (s *Session) SetScriptTimeout(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts/async_script", msRequest{
		MS: timeout,
	}, nil)
}
