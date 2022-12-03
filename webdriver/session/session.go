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
func OpenWithClient(client *http.Client, url string, capabilities map[string]any) (*Session, error) {
	c, err := newConnection(client, url, capabilities)
	if err != nil {
		return nil, err
	}
	return &Session{Connection: c}, nil
}

// Delete sends to delete message to terminate the session.
func (s *Session) Delete() error {
	return s.DeleteWithContext(context.Background())
}

// DeleteWithContext sends to delete message to terminate the session.
func (s *Session) DeleteWithContext(ctx context.Context) error {
	return s.Send(ctx, Delete, "", nil, nil)
}

// Selector represents a selector for elements.
type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

// GetElement retrieves the element matching the selector query of the session.
func (s *Session) GetElement(selector Selector) (*Element, error) {
	return s.GetElementWithContext(context.Background(), selector)
}

// GetElementWithContext retrieves the element matching the selector query of the session.
func (s *Session) GetElementWithContext(ctx context.Context, selector Selector) (*Element, error) {
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
func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	return s.GetElementsWithContext(context.Background(), selector)
}

// GetElementsWithContext retrieves elements matching the selector query of the session.
func (s *Session) GetElementsWithContext(ctx context.Context, selector Selector) ([]*Element, error) {
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
func (s *Session) GetActiveElement() (*Element, error) {
	return s.GetActiveElementWithContext(context.Background())
}

// GetActiveElementWithContext returns the active element of the session.
func (s *Session) GetActiveElementWithContext(ctx context.Context) (*Element, error) {
	var result elementResult
	if err := s.Send(ctx, Post, "element/active", nil, &result); err != nil {
		return nil, err
	}
	return &Element{ID: result.ID(), Session: s}, nil
}

// GetWindow returns the window handler of the session.
func (s *Session) GetWindow() (*Window, error) {
	return s.GetWindowWithContext(context.Background())
}

// GetWindowWithContext returns the window handler of the session.
func (s *Session) GetWindowWithContext(ctx context.Context) (*Window, error) {
	var windowID string
	if err := s.Send(ctx, Get, "window_handle", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{ID: windowID, Session: s}, nil
}

// GetWindows returns window handlers of the session.
func (s *Session) GetWindows() ([]*Window, error) {
	return s.GetWindowsWithContext(context.Background())
}

// GetWindowsWithContext returns window handlers of the session.
func (s *Session) GetWindowsWithContext(ctx context.Context) ([]*Window, error) {
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
func (s *Session) SetWindow(window *Window) error {
	return s.SetWindowWithContext(context.Background(), window)
}

// SetWindowWithContext sets the window to the browser.
func (s *Session) SetWindowWithContext(ctx context.Context, window *Window) error {
	if window == nil {
		return errors.New("nil window is invalid")
	}
	return s.Send(ctx, Post, "window", nameRequest{
		Name: window.ID,
	}, nil)
}

// SetWindowByName sets the window to the browser by name.
func (s *Session) SetWindowByName(name string) error {
	return s.SetWindowByNameWithContext(context.Background(), name)
}

// SetWindowByNameWithContext sets the window to the browser by name.
func (s *Session) SetWindowByNameWithContext(ctx context.Context, name string) error {
	return s.Send(ctx, Post, "window", nameRequest{
		Name: name,
	}, nil)
}

// DeleteWindow deletes the window of the session.
func (s *Session) DeleteWindow() error {
	return s.DeleteWindowWithContext(context.Background())
}

// DeleteWindowWithContext deletes the window of the session.
func (s *Session) DeleteWindowWithContext(ctx context.Context) error {
	if err := s.Send(ctx, Delete, "window", nil, nil); err != nil {
		return err
	}
	return nil
}

// GetCookies gets cookies of the session.
func (s *Session) GetCookies() ([]*Cookie, error) {
	return s.GetCookiesWithContext(context.Background())
}

// GetCookiesWithContext gets cookies of the session.
func (s *Session) GetCookiesWithContext(ctx context.Context) ([]*Cookie, error) {
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
func (s *Session) SetCookie(cookie *Cookie) error {
	return s.SetCookieWithContext(context.Background(), cookie)
}

// SetCookieWithContext sets a cookie to the browser.
func (s *Session) SetCookieWithContext(ctx context.Context, cookie *Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}
	return s.Send(ctx, Post, "cookie", cookieRequest{
		Cookie: cookie,
	}, nil)
}

// DeleteCookie deletes a cookie of the session.
func (s *Session) DeleteCookie(cookieName string) error {
	return s.DeleteCookieWithContext(context.Background(), cookieName)
}

// DeleteCookieWithContext deletes a cookie of the session.
func (s *Session) DeleteCookieWithContext(ctx context.Context, cookieName string) error {
	return s.Send(ctx, Delete, path.Join("cookie", cookieName), nil, nil)
}

// DeleteCookies deletes cookies of the session.
func (s *Session) DeleteCookies() error {
	return s.DeleteCookiesWithContext(context.Background())
}

// DeleteCookiesWithContext deletes cookies of the session.
func (s *Session) DeleteCookiesWithContext(ctx context.Context) error {
	return s.Send(ctx, Delete, "cookie", nil, nil)
}

// GetScreenshot gets a screenshot.
func (s *Session) GetScreenshot() ([]byte, error) {
	return s.GetScreenshotWithContext(context.Background())
}

// GetScreenshotWithContext gets a screenshot.
func (s *Session) GetScreenshotWithContext(ctx context.Context) ([]byte, error) {
	var base64Image string
	if err := s.Send(ctx, Get, "screenshot", nil, &base64Image); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(base64Image)
}

// GetURL gets the url of the session.
func (s *Session) GetURL() (string, error) {
	return s.GetURLWithContext(context.Background())
}

// GetURLWithContext gets the url of the session.
func (s *Session) GetURLWithContext(ctx context.Context) (string, error) {
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
func (s *Session) SetURL(url string) error {
	return s.SetURLWithContext(context.Background(), url)
}

// SetURLWithContext sets the url to the browser.
func (s *Session) SetURLWithContext(ctx context.Context, url string) error {
	return s.Send(ctx, Post, "url", urlRequest{
		URL: url,
	}, nil)
}

// GetTitle gets the title of the session.
func (s *Session) GetTitle() (string, error) {
	return s.GetTitleWithContext(context.Background())
}

// GetTitleWithContext gets the title of the session.
func (s *Session) GetTitleWithContext(ctx context.Context) (string, error) {
	var title string
	if err := s.Send(ctx, Get, "title", nil, &title); err != nil {
		return "", err
	}
	return title, nil
}

// GetSource gets the page source of the session.
func (s *Session) GetSource() (string, error) {
	return s.GetSourceWithContext(context.Background())
}

// GetSourceWithContext gets the page source of the session.
func (s *Session) GetSourceWithContext(ctx context.Context) (string, error) {
	var source string
	if err := s.Send(ctx, Get, "source", nil, &source); err != nil {
		return "", err
	}
	return source, nil
}

// MoveTo moves the element to the offset position.
func (s *Session) MoveTo(region *Element, offset Offset) error {
	return s.MoveToWithContext(context.Background(), region, offset)
}

// MoveToWithContext moves the element to the offset position.
func (s *Session) MoveToWithContext(ctx context.Context, region *Element, offset Offset) error {
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
func (s *Session) Frame(frame *Element) error {
	return s.FrameWithContext(context.Background(), frame)
}

// FrameWithContext sets the frame to the browser.
func (s *Session) FrameWithContext(ctx context.Context, frame *Element) error {
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
func (s *Session) FrameParent() error {
	return s.FrameParentWithContext(context.Background())
}

// FrameParentWithContext sets the parent frame to the browser.
func (s *Session) FrameParentWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "frame/parent", nil, nil)
}

type scriptRequest struct {
	Script string `json:"script"`
	Args   []any  `json:"args"`
}

// Execute executes the script.
func (s *Session) Execute(body string, arguments []any, result any) error {
	return s.ExecuteWithContext(context.Background(), body, arguments, result)
}

// ExecuteWithContext executes the script.
func (s *Session) ExecuteWithContext(ctx context.Context, body string, arguments []any, result any) error {
	if arguments == nil {
		arguments = []any{}
	}
	return s.Send(ctx, Post, "execute", scriptRequest{
		Script: body,
		Args:   arguments,
	}, result)
}

// Forward forwards the browser.
func (s *Session) Forward() error {
	return s.ForwardWithContext(context.Background())
}

// ForwardWithContext forwards the browser.
func (s *Session) ForwardWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "forward", nil, nil)
}

// Back backs the browser.
func (s *Session) Back() error {
	return s.BackWithContext(context.Background())
}

// BackWithContext backs the browser.
func (s *Session) BackWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "back", nil, nil)
}

// Refresh refreshes the browser.
func (s *Session) Refresh() error {
	return s.RefreshWithContext(context.Background())
}

// RefreshWithContext refreshes the browser.
func (s *Session) RefreshWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "refresh", nil, nil)
}

// GetAlertText gets the alert text of the browser.
func (s *Session) GetAlertText() (string, error) {
	return s.GetAlertTextWithContext(context.Background())
}

// GetAlertTextWithContext gets the alert text of the browser.
func (s *Session) GetAlertTextWithContext(ctx context.Context) (string, error) {
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
func (s *Session) SetAlertText(text string) error {
	return s.SetAlertTextWithContext(context.Background(), text)
}

// SetAlertTextWithContext sets the text to the browser.
func (s *Session) SetAlertTextWithContext(ctx context.Context, text string) error {
	return s.Send(ctx, Post, "alert_text", textRequest{
		Text: text,
	}, nil)
}

// AcceptAlert accepts the alert of the browser.
func (s *Session) AcceptAlert() error {
	return s.AcceptAlertWithContext(context.Background())
}

// AcceptAlertWithContext accepts the alert of the browser.
func (s *Session) AcceptAlertWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "accept_alert", nil, nil)
}

// DismissAlert dismisses the alert of the browser.
func (s *Session) DismissAlert() error {
	return s.DismissAlertWithContext(context.Background())
}

// DismissAlertWithContext dismisses the alert of the browser.
func (s *Session) DismissAlertWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "dismiss_alert", nil, nil)
}

type typeRequest struct {
	Type string `json:"type"`
}

// NewLogs gets logs of the browser.
func (s *Session) NewLogs(logType string) ([]Log, error) {
	return s.NewLogsWithContext(context.Background(), logType)
}

// NewLogsWithContext gets logs of the browser.
func (s *Session) NewLogsWithContext(ctx context.Context, logType string) ([]Log, error) {
	var logs []Log
	if err := s.Send(ctx, Post, "log", typeRequest{
		Type: logType,
	}, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLogTypes gets log types.
func (s *Session) GetLogTypes() ([]string, error) {
	return s.GetLogTypesWithContext(context.Background())
}

// GetLogTypesWithContext gets log types.
func (s *Session) GetLogTypesWithContext(ctx context.Context) ([]string, error) {
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
func (s *Session) DoubleClick() error {
	return s.DoubleClickWithContext(context.Background())
}

// DoubleClickWithContext sends the double click event to the browser.
func (s *Session) DoubleClickWithContext(ctx context.Context) error {
	return s.Send(ctx, Post, "doubleclick", nil, nil)
}

// Click sends the click event to the browser.
func (s *Session) Click(button event.Button) error {
	return s.ClickWithContext(context.Background(), button)
}

// ClickWithContext sends the click event to the browser.
func (s *Session) ClickWithContext(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "click", buttonRequest{Button: button}, nil)
}

// ButtonDown sends the button down event to the browser.
func (s *Session) ButtonDown(button event.Button) error {
	return s.ButtonDownWithContext(context.Background(), button)
}

// ButtonDownWithContext sends the button down event to the browser.
func (s *Session) ButtonDownWithContext(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "buttondown", buttonRequest{Button: button}, nil)
}

// ButtonUp sends the button up event to the browser.
func (s *Session) ButtonUp(button event.Button) error {
	return s.ButtonUpWithContext(context.Background(), button)
}

// ButtonUpWithContext sends the button up event to the browser.
func (s *Session) ButtonUpWithContext(ctx context.Context, button event.Button) error {
	return s.Send(ctx, Post, "buttonup", buttonRequest{Button: button}, nil)
}

type xyRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// TouchDown sends the touch-down event to the browser.
func (s *Session) TouchDown(x, y int) error {
	return s.TouchDownWithContext(context.Background(), x, y)
}

// TouchDownWithContext sends the touch-down event to the browser.
func (s *Session) TouchDownWithContext(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/down", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

// TouchUp sends the touch-up event to the browser.
func (s *Session) TouchUp(x, y int) error {
	return s.TouchUpWithContext(context.Background(), x, y)
}

// TouchUpWithContext sends the touch-up event to the browser.
func (s *Session) TouchUpWithContext(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/up", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

// TouchMove sends the touch-move event to the browser.
func (s *Session) TouchMove(x, y int) error {
	return s.TouchMoveWithContext(context.Background(), x, y)
}

// TouchMoveWithContext sends the touch-move event to the browser.
func (s *Session) TouchMoveWithContext(ctx context.Context, x, y int) error {
	return s.Send(ctx, Post, "touch/move", xyRequest{
		X: x,
		Y: y,
	}, nil)
}

type elementRequest struct {
	Element string `json:"element"`
}

// TouchClick sends touch-click event to the browser.
func (s *Session) TouchClick(element *Element) error {
	return s.TouchClickWithContext(context.Background(), element)
}

// TouchClickWithContext sends touch-click event to the browser.
func (s *Session) TouchClickWithContext(ctx context.Context, element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(ctx, Post, "touch/click", elementRequest{
		Element: element.ID,
	}, nil)
}

// TouchDoubleClick sends the touch-double-click event to the browser.
func (s *Session) TouchDoubleClick(element *Element) error {
	return s.TouchDoubleClickWithContext(context.Background(), element)
}

// TouchDoubleClickWithContext sends the touch-double-click event to the browser.
func (s *Session) TouchDoubleClickWithContext(ctx context.Context, element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(ctx, Post, "touch/doubleclick", elementRequest{
		Element: element.ID,
	}, nil)
}

// TouchLongClick sends the touch-long-click event to the browser.
func (s *Session) TouchLongClick(element *Element) error {
	return s.TouchLongClickWithContext(context.Background(), element)
}

// TouchLongClickWithContext sends the touch-long-click event to the browser.
func (s *Session) TouchLongClickWithContext(ctx context.Context, element *Element) error {
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
func (s *Session) TouchFlick(element *Element, offset Offset, speed Speed) error {
	return s.TouchFlickWithContext(context.Background(), element, offset, speed)
}

// TouchFlickWithContext sends the touch-flick event to the browser.
func (s *Session) TouchFlickWithContext(ctx context.Context, element *Element, offset Offset, speed Speed) error {
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
func (s *Session) TouchScroll(element *Element, offset Offset) error {
	return s.TouchScrollWithContext(context.Background(), element, offset)
}

// TouchScrollWithContext sends the touch-scroll event to the browser.
func (s *Session) TouchScrollWithContext(ctx context.Context, element *Element, offset Offset) error {
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
func (s *Session) Keys(text string) error {
	return s.KeysWithContext(context.Background(), text)
}

// KeysWithContext sends key events of the text to the browser.
func (s *Session) KeysWithContext(ctx context.Context, text string) error {
	return s.Send(ctx, Post, "keys", valueSliceRequest{
		Value: strings.Split(text, ""),
	}, nil)
}

// DeleteLocalStorage deletes the local storage of the browser.
func (s *Session) DeleteLocalStorage() error {
	return s.DeleteLocalStorageWithContext(context.Background())
}

// DeleteLocalStorageWithContext deletes the local storage of the browser.
func (s *Session) DeleteLocalStorageWithContext(ctx context.Context) error {
	return s.Send(ctx, Delete, "local_storage", nil, nil)
}

// DeleteSessionStorage deletes the session storage of the browser.
func (s *Session) DeleteSessionStorage() error {
	return s.DeleteSessionStorageWithContext(context.Background())
}

// DeleteSessionStorageWithContext deletes the session storage of the browser.
func (s *Session) DeleteSessionStorageWithContext(ctx context.Context) error {
	return s.Send(ctx, Delete, "session_storage", nil, nil)
}

type msRequest struct {
	MS   int    `json:"ms"`
	Type string `json:"type,omitempty"`
}

// SetImplicitWait sets the implicit wait to the browser.
func (s *Session) SetImplicitWait(timeout int) error {
	return s.SetImplicitWaitWithContext(context.Background(), timeout)
}

// SetImplicitWaitWithContext sets the implicit wait to the browser.
func (s *Session) SetImplicitWaitWithContext(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts/implicit_wait", msRequest{
		MS: timeout,
	}, nil)
}

// SetPageLoad sets the timeout to the page load of the browser.
func (s *Session) SetPageLoad(timeout int) error {
	return s.SetPageLoadWithContext(context.Background(), timeout)
}

// SetPageLoadWithContext sets the timeout to the page load of the browser.
func (s *Session) SetPageLoadWithContext(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts", msRequest{
		MS:   timeout,
		Type: "page load",
	}, nil)
}

// SetScriptTimeout sets the timeout to the asynchronous scripts execution.
func (s *Session) SetScriptTimeout(timeout int) error {
	return s.SetScriptTimeoutWithContext(context.Background(), timeout)
}

// SetScriptTimeoutWithContext sets the timeout to the asynchronous scripts execution.
func (s *Session) SetScriptTimeoutWithContext(ctx context.Context, timeout int) error {
	return s.Send(ctx, Post, "timeouts/async_script", msRequest{
		MS: timeout,
	}, nil)
}
