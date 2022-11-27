package session

import (
	"encoding/base64"
	"errors"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ikawaha/navi/event"
)

const (
	DefaultWebdriverTimeout     = 10 * time.Second
	DefaultSessionClientTimeout = 30 * time.Second
)

type Method = string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Delete Method = "DELETE"
)

type Session struct {
	*Connection
}

func OpenWithClient(client *http.Client, url string, capabilities map[string]any) (*Session, error) {
	c, err := newConnection(client, url, capabilities)
	if err != nil {
		return nil, err
	}
	return &Session{Connection: c}, nil
}

func (s *Session) Delete() error {
	return s.Send(Delete, "", nil, nil)
}

type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

func (s *Session) GetElement(selector Selector) (*Element, error) {
	var result elementResult
	if err := s.Send(Post, "element", selector, &result); err != nil {
		return nil, err
	}
	return &Element{
		ID:      result.ID(),
		Session: s,
	}, nil
}

func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	var results []elementResult
	if err := s.Send(Post, "elements", selector, &results); err != nil {
		return nil, err
	}

	var elements []*Element
	for _, result := range results {
		elements = append(elements, &Element{ID: result.ID(), Session: s})
	}

	return elements, nil
}

func (s *Session) GetActiveElement() (*Element, error) {
	var result elementResult
	if err := s.Send(Post, "element/active", nil, &result); err != nil {
		return nil, err
	}
	return &Element{ID: result.ID(), Session: s}, nil
}

func (s *Session) GetWindow() (*Window, error) {
	var windowID string
	if err := s.Send(Get, "window_handle", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{ID: windowID, Session: s}, nil
}

func (s *Session) GetWindows() ([]*Window, error) {
	var windowsID []string
	if err := s.Send(Get, "window_handles", nil, &windowsID); err != nil {
		return nil, err
	}

	var windows []*Window
	for _, windowID := range windowsID {
		windows = append(windows, &Window{windowID, s})
	}
	return windows, nil
}

type NameRequest struct {
	Name string `json:"name"`
}

func (s *Session) SetWindow(window *Window) error {
	if window == nil {
		return errors.New("nil window is invalid")
	}
	return s.Send(Post, "window", NameRequest{
		Name: window.ID,
	}, nil)
}

func (s *Session) SetWindowByName(name string) error {
	return s.Send(Post, "window", NameRequest{
		Name: name,
	}, nil)
}

func (s *Session) DeleteWindow() error {
	if err := s.Send(Delete, "window", nil, nil); err != nil {
		return err
	}
	return nil
}

func (s *Session) GetCookies() ([]*Cookie, error) {
	var cookies []*Cookie
	if err := s.Send(Get, "cookie", nil, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

type CookieRequest struct {
	Cookie *Cookie `json:"cookie"`
}

func (s *Session) SetCookie(cookie *Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}
	return s.Send(Post, "cookie", CookieRequest{
		Cookie: cookie,
	}, nil)
}

func (s *Session) DeleteCookie(cookieName string) error {
	return s.Send(Delete, path.Join("cookie", cookieName), nil, nil)
}

func (s *Session) DeleteCookies() error {
	return s.Send(Delete, "cookie", nil, nil)
}

func (s *Session) GetScreenshot() ([]byte, error) {
	var base64Image string
	if err := s.Send(Get, "screenshot", nil, &base64Image); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(base64Image)
}

func (s *Session) GetURL() (string, error) {
	var url string
	if err := s.Send(Get, "url", nil, &url); err != nil {
		return "", err
	}
	return url, nil
}

type URLRequest struct {
	URL string `json:"url"`
}

func (s *Session) SetURL(url string) error {
	return s.Send(Post, "url", URLRequest{
		URL: url,
	}, nil)
}

func (s *Session) GetTitle() (string, error) {
	var title string
	if err := s.Send(Get, "title", nil, &title); err != nil {
		return "", err
	}
	return title, nil
}

func (s *Session) GetSource() (string, error) {
	var source string
	if err := s.Send(Get, "source", nil, &source); err != nil {
		return "", err
	}
	return source, nil
}

func (s *Session) MoveTo(region *Element, offset Offset) error {
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
	return s.Send(Post, "moveto", req, nil)
}

func (s *Session) Frame(frame *Element) error {
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
	return s.Send(Post, "frame", req, nil)
}

func (s *Session) FrameParent() error {
	return s.Send(Post, "frame/parent", nil, nil)
}

type ScriptRequest struct {
	Script string `json:"script"`
	Args   []any  `json:"args"`
}

func (s *Session) Execute(body string, arguments []any, result any) error {
	if arguments == nil {
		arguments = []any{}
	}
	return s.Send(Post, "execute", ScriptRequest{
		Script: body,
		Args:   arguments,
	}, result)
}

func (s *Session) Forward() error {
	return s.Send(Post, "forward", nil, nil)
}

func (s *Session) Back() error {
	return s.Send(Post, "back", nil, nil)
}

func (s *Session) Refresh() error {
	return s.Send(Post, "refresh", nil, nil)
}

func (s *Session) GetAlertText() (string, error) {
	var text string
	if err := s.Send(Get, "alert_text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

type TextRequest struct {
	Text string `json:"text"`
}

func (s *Session) SetAlertText(text string) error {
	return s.Send(Post, "alert_text", TextRequest{
		Text: text,
	}, nil)
}

func (s *Session) AcceptAlert() error {
	return s.Send(Post, "accept_alert", nil, nil)
}

func (s *Session) DismissAlert() error {
	return s.Send(Post, "dismiss_alert", nil, nil)
}

type TypeRequest struct {
	Type string `json:"type"`
}

func (s *Session) NewLogs(logType string) ([]Log, error) {
	var logs []Log
	if err := s.Send(Post, "log", TypeRequest{
		Type: logType,
	}, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Session) GetLogTypes() ([]string, error) {
	var types []string
	if err := s.Send(Get, "log/types", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}

type ButtonRequest struct {
	Button event.Button `json:"button"`
}

func (s *Session) DoubleClick() error {
	return s.Send(Post, "doubleclick", nil, nil)
}

func (s *Session) Click(button event.Button) error {
	return s.Send(Post, "click", ButtonRequest{Button: button}, nil)
}

func (s *Session) ButtonDown(button event.Button) error {
	return s.Send(Post, "buttondown", ButtonRequest{Button: button}, nil)
}

func (s *Session) ButtonUp(button event.Button) error {
	return s.Send(Post, "buttonup", ButtonRequest{Button: button}, nil)
}

type XYRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s *Session) TouchDown(x, y int) error {
	return s.Send(Post, "touch/down", XYRequest{
		X: x,
		Y: y,
	}, nil)
}

func (s *Session) TouchUp(x, y int) error {
	return s.Send(Post, "touch/up", XYRequest{
		X: x,
		Y: y,
	}, nil)
}

func (s *Session) TouchMove(x, y int) error {
	return s.Send(Post, "touch/move", XYRequest{
		X: x,
		Y: y,
	}, nil)
}

type ElementRequest struct {
	Element string `json:"element"`
}

func (s *Session) TouchClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(Post, "touch/click", ElementRequest{
		Element: element.ID,
	}, nil)
}

func (s *Session) TouchDoubleClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(Post, "touch/doubleclick", ElementRequest{
		Element: element.ID,
	}, nil)
}

func (s *Session) TouchLongClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}
	return s.Send(Post, "touch/longclick", ElementRequest{
		Element: element.ID,
	}, nil)
}

type XYSpeedRequest struct {
	XSpeed int `json:"xspeed"`
	YSpeed int `json:"yspeed"`
}

type TouchFlickRequest struct {
	Element string `json:"element"`
	XOffset int    `json:"xoffset"`
	YOffset int    `json:"yoffset"`
	Speed   uint   `json:"speed"`
}

func (s *Session) TouchFlick(element *Element, offset Offset, speed Speed) error {
	if speed == nil {
		return errors.New("nil speed is invalid")
	}
	if (element == nil) != (offset == nil) {
		return errors.New("element must be provided if offset is provided and vice versa")
	}
	if element == nil {
		xSpeed, ySpeed := speed.vector()
		return s.Send(Post, "touch/flick", XYSpeedRequest{
			XSpeed: xSpeed,
			YSpeed: ySpeed,
		}, nil)
	}
	xOffset, yOffset := offset.position()
	return s.Send(Post, "touch/flick", TouchFlickRequest{
		Element: element.ID,
		XOffset: xOffset,
		YOffset: yOffset,
		Speed:   speed.scalar(),
	}, nil)
}

type TouchScrollRequest struct {
	Element string `json:"element,omitempty"`
	XOffset int    `json:"xoffset"`
	YOffset int    `json:"yoffset"`
}

func (s *Session) TouchScroll(element *Element, offset Offset) error {
	if element == nil {
		element = &Element{}
	}
	if offset == nil {
		return errors.New("nil offset is invalid")
	}
	xOffset, yOffset := offset.position()
	return s.Send(Post, "touch/scroll", TouchScrollRequest{
		Element: element.ID,
		XOffset: xOffset,
		YOffset: yOffset,
	}, nil)
}

type ValueSliceRequest struct {
	Value []string `json:"value"`
}

func (s *Session) Keys(text string) error {
	return s.Send(Post, "keys", ValueSliceRequest{
		Value: strings.Split(text, ""),
	}, nil)
}

func (s *Session) DeleteLocalStorage() error {
	return s.Send(Delete, "local_storage", nil, nil)
}

func (s *Session) DeleteSessionStorage() error {
	return s.Send(Delete, "session_storage", nil, nil)
}

type MSRequest struct {
	MS   int    `json:"ms"`
	Type string `json:"type,omitempty"`
}

func (s *Session) SetImplicitWait(timeout int) error {
	return s.Send(Post, "timeouts/implicit_wait", MSRequest{
		MS: timeout,
	}, nil)
}

func (s *Session) SetPageLoad(timeout int) error {
	return s.Send(Post, "timeouts", MSRequest{
		MS:   timeout,
		Type: "page load",
	}, nil)
}

func (s *Session) SetScriptTimeout(timeout int) error {
	return s.Send(Post, "timeouts/async_script", MSRequest{
		MS: timeout,
	}, nil)
}
