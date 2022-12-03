package navigator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ikawaha/navigator/event"
	"github.com/ikawaha/navigator/webdriver/session"
)

const aboutBlankURL = "about:blank"

// A Page represents an open browser session. Pages may be created using the
// *WebDriver.Page() method.
type Page struct {
	Selectable
	logs map[string][]Log
}

func newPage(session *session.Session) *Page {
	return &Page{
		Selectable: Selectable{
			session: session,
		},
	}
}

// String returns a string representation of the Page. Currently: "page"
func (p *Page) String() string {
	return "page"
}

// Session returns a *webdriver.Session that can be used to send direct commands
// to the WebDriver. See: https://code.google.com/p/selenium/wiki/JsonWireProtocol
func (p *Page) Session() *session.Session {
	return p.session
}

// Destroy closes any open browsers by ending the session.
func (p *Page) Destroy() error {
	if err := p.session.Delete(); err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}
	return nil
}

// Reset deletes all cookies set for the current domain and navigates to a blank page.
// Unlike Destroy, Reset will permit the page to be re-used after it is called.
// Reset is faster than Destroy, but any cookies from domains outside the current
// domain will remain after a page is reset.
func (p *Page) Reset() error {
	_ = p.ConfirmPopup()
	url, err := p.URL()
	if err != nil {
		return err
	}
	if url == aboutBlankURL {
		return nil
	}
	if err := p.DeleteCookies(); err != nil {
		return err
	}
	if err := p.session.DeleteLocalStorage(); err != nil {
		if err := p.RunScript("localStorage.clear();", nil, nil); err != nil {
			return err
		}
	}
	if err := p.session.DeleteSessionStorage(); err != nil {
		if err := p.RunScript("sessionStorage.clear();", nil, nil); err != nil {
			return err
		}
	}
	return p.Navigate(aboutBlankURL)
}

// Navigate navigates to the provided URL.
func (p *Page) Navigate(url string) error {
	if err := p.session.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}
	return nil
}

// GetCookies returns all cookies on the page.
func (p *Page) GetCookies() ([]*http.Cookie, error) {
	cookies, err := p.session.GetCookies()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}
	var ret []*http.Cookie
	for _, c := range cookies {
		expSeconds := int64(c.Expiry)
		expNano := int64(c.Expiry-float64(expSeconds)) * 1000000000
		ret = append(ret, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Secure:   c.Secure,
			HttpOnly: c.HTTPOnly,
			Expires:  time.Unix(expSeconds, expNano),
		})
	}
	return ret, nil
}

// SetCookie sets a cookie on the page.
func (p *Page) SetCookie(cookie *http.Cookie) error {
	if cookie == nil {
		return nil
	}
	var expiry int64
	if !cookie.Expires.IsZero() {
		expiry = cookie.Expires.Unix()
	}
	if err := p.session.SetCookie(&session.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		Expiry:   float64(expiry),
	}); err != nil {
		return fmt.Errorf("failed to set cookie: %w", err)
	}
	return nil
}

// DeleteCookie deletes a cookie on the page by name.
func (p *Page) DeleteCookie(name string) error {
	if err := p.session.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %w", name, err)
	}
	return nil
}

// DeleteCookies deletes all cookies on the page.
func (p *Page) DeleteCookies() error {
	if err := p.session.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %w", err)
	}
	return nil
}

// URL returns the current page URL.
func (p *Page) URL() (string, error) {
	url, err := p.session.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %w", err)
	}
	return url, nil
}

// Size sets the current page size in pixels.
func (p *Page) Size(width, height int) error {
	window, err := p.session.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %w", err)
	}
	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %w", err)
	}
	return nil
}

// Screenshot takes a screenshot and saves it to the provided filename.
// The provided filename may be an absolute or relative path.
func (p *Page) Screenshot(filename string) error {
	path, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to find absolute path for filename: %w", err)
	}
	screenshot, err := p.session.GetScreenshot()
	if err != nil {
		return fmt.Errorf("failed to retrieve screenshot: %w", err)
	}
	if err := os.WriteFile(path, screenshot, 0o644); err != nil {
		return fmt.Errorf("failed to save screenshot: %w", err)
	}
	return nil
}

// Title returns the page title.
func (p *Page) Title() (string, error) {
	return p.TitleWithContext(context.Background())
}

// TitleWithContext returns the page title.
func (p *Page) TitleWithContext(ctx context.Context) (string, error) {
	title, err := p.session.GetTitleWithContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %w", err)
	}
	return title, nil
}

// HTML returns the current contents of the DOM for the entire page.
func (p *Page) HTML() (string, error) {
	return p.HTMLWithContext(context.Background())
}

// HTMLWithContext returns the current contents of the DOM for the entire page.
func (p *Page) HTMLWithContext(ctx context.Context) (string, error) {
	html, err := p.session.GetSourceWithContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page HTML: %w", err)
	}
	return html, nil
}

// RunScript runs the JavaScript provided in the body. Any keys present in
// the arguments map will be available as variables in the body.
// Values provided in arguments are converted into javascript objects.
// If the body returns a value, it will be unmarshalled into the result argument.
// Simple example:
//
//	var number int
//	page.RunScript("return test;", map[string]any{"test": 100}, &number)
//	fmt.Println(number)
//
// -> 100
func (p *Page) RunScript(body string, arguments map[string]any, result any) error {
	return p.RunScriptWithContext(context.Background(), body, arguments, result)
}

// RunScriptWithContext runs the JavaScript provided in the body. Any keys present in
// the arguments map will be available as variables in the body.
// Values provided in arguments are converted into javascript objects.
// If the body returns a value, it will be unmarshalled into the result argument.
// Simple example:
//
//	var number int
//	page.RunScript("return test;", map[string]any{"test": 100}, &number)
//	fmt.Println(number)
//
// -> 100
func (p *Page) RunScriptWithContext(ctx context.Context, body string, arguments map[string]any, result any) error {
	var keys []string
	var values []any
	for key, value := range arguments {
		keys = append(keys, key)
		values = append(values, value)
	}
	argumentList := strings.Join(keys, ", ")
	cleanBody := fmt.Sprintf("return (function(%s) { %s; }).apply(this, arguments);", argumentList, body)
	if err := p.session.ExecuteWithContext(ctx, cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}
	return nil
}

// PopupText returns the current alert, confirm, or prompt popup text.
func (p *Page) PopupText() (string, error) {
	return p.PopupTextWithContext(context.Background())
}

// PopupTextWithContext returns the current alert, confirm, or prompt popup text.
func (p *Page) PopupTextWithContext(ctx context.Context) (string, error) {
	text, err := p.session.GetAlertTextWithContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %w", err)
	}
	return text, nil
}

// EnterPopupText enters text into an open prompt popup.
func (p *Page) EnterPopupText(text string) error {
	return p.EnterPopupTextWithContext(context.Background(), text)
}

// EnterPopupTextWithContext enters text into an open prompt popup.
func (p *Page) EnterPopupTextWithContext(ctx context.Context, text string) error {
	if err := p.session.SetAlertTextWithContext(ctx, text); err != nil {
		return fmt.Errorf("failed to enter popup text: %w", err)
	}
	return nil
}

// ConfirmPopup confirms an alert, confirm, or prompt popup.
func (p *Page) ConfirmPopup() error {
	return p.ConfirmPopupWithContext(context.Background())
}

// ConfirmPopupWithContext confirms an alert, confirm, or prompt popup.
func (p *Page) ConfirmPopupWithContext(ctx context.Context) error {
	if err := p.session.AcceptAlertWithContext(ctx); err != nil {
		return fmt.Errorf("failed to confirm popup: %w", err)
	}
	return nil
}

// CancelPopup cancels an alert, confirm, or prompt popup.
func (p *Page) CancelPopup() error {
	return p.CancelPopupWithContext(context.Background())
}

// CancelPopupWithContext cancels an alert, confirm, or prompt popup.
func (p *Page) CancelPopupWithContext(ctx context.Context) error {
	if err := p.session.DismissAlertWithContext(ctx); err != nil {
		return fmt.Errorf("failed to cancel popup: %w", err)
	}
	return nil
}

// Forward navigates forward in history.
func (p *Page) Forward() error {
	return p.ForwardWithContext(context.Background())
}

// ForwardWithContext navigates forward in history.
func (p *Page) ForwardWithContext(ctx context.Context) error {
	if err := p.session.ForwardWithContext(ctx); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %w", err)
	}
	return nil
}

// Back navigates backwards in history.
func (p *Page) Back() error {
	return p.BackWithContext(context.Background())
}

// BackWithContext navigates backwards in history.
func (p *Page) BackWithContext(ctx context.Context) error {
	if err := p.session.BackWithContext(ctx); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %w", err)
	}
	return nil
}

// Refresh refreshes the page.
func (p *Page) Refresh() error {
	return p.RefreshWithContext(context.Background())
}

// RefreshWithContext refreshes the page.
func (p *Page) RefreshWithContext(ctx context.Context) error {
	if err := p.session.RefreshWithContext(ctx); err != nil {
		return fmt.Errorf("failed to refresh page: %w", err)
	}
	return nil
}

// SwitchToParentFrame focuses on the immediate parent frame of a frame selected
// by Selection.Frame. After switching, all new and existing selections will refer
// to the parent frame. All further Page methods will apply to this frame as well.
//
// This method is not supported by PhantomJS. Please use SwitchToRootFrame instead.
func (p *Page) SwitchToParentFrame() error {
	return p.SwitchToParentFrameWithContext(context.Background())
}

// SwitchToParentFrameWithContext focuses on the immediate parent frame of a frame selected
// by Selection.Frame. After switching, all new and existing selections will refer
// to the parent frame. All further Page methods will apply to this frame as well.
//
// This method is not supported by PhantomJS. Please use SwitchToRootFrame instead.
func (p *Page) SwitchToParentFrameWithContext(ctx context.Context) error {
	if err := p.session.FrameParentWithContext(ctx); err != nil {
		return fmt.Errorf("failed to switch to parent frame: %w", err)
	}
	return nil
}

// SwitchToRootFrame focuses on the original, default page frame before any calls
// to Selection.Frame were made. After switching, all new and existing selections
// will refer to the root frame. All further Page methods will apply to this frame
// as well.
func (p *Page) SwitchToRootFrame() error {
	return p.SwitchToParentFrameWithContext(context.Background())
}

// SwitchToRootFrameWithContext focuses on the original, default page frame before any calls
// to Selection.Frame were made. After switching, all new and existing selections
// will refer to the root frame. All further Page methods will apply to this frame
// as well.
func (p *Page) SwitchToRootFrameWithContext(ctx context.Context) error {
	if err := p.session.FrameWithContext(ctx, nil); err != nil {
		return fmt.Errorf("failed to switch to original page frame: %w", err)
	}
	return nil
}

// SwitchToWindow switches to the first available window with the provided name
// (JavaScript `window.name` attribute).
func (p *Page) SwitchToWindow(name string) error {
	return p.SwitchToWindowWithContext(context.Background(), name)
}

// SwitchToWindowWithContext switches to the first available window with the provided name
// (JavaScript `window.name` attribute).
func (p *Page) SwitchToWindowWithContext(ctx context.Context, name string) error {
	if err := p.session.SetWindowByNameWithContext(ctx, name); err != nil {
		return fmt.Errorf("failed to switch to named window: %w", err)
	}
	return nil
}

// NextWindow switches to the next available window.
func (p *Page) NextWindow() error {
	return p.NextWindowWithContext(context.Background())
}

// NextWindowWithContext switches to the next available window.
func (p *Page) NextWindowWithContext(ctx context.Context) error {
	windows, err := p.session.GetWindowsWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to find available windows: %w", err)
	}
	var windowIDs []string
	for _, v := range windows {
		windowIDs = append(windowIDs, v.ID)
	}

	// order not defined according to W3 spec
	sort.Strings(windowIDs)

	activeWindow, err := p.session.GetWindowWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to find active window: %w", err)
	}
	for position, windowID := range windowIDs {
		if windowID == activeWindow.ID {
			activeWindow.ID = windowIDs[(position+1)%len(windowIDs)]
			break
		}
	}
	if err := p.session.SetWindowWithContext(ctx, activeWindow); err != nil {
		return fmt.Errorf("failed to change active window: %w", err)
	}
	return nil
}

// CloseWindow closes the active window.
func (p *Page) CloseWindow() error {
	return p.CloseWindowWithContext(context.Background())
}

// CloseWindowWithContext closes the active window.
func (p *Page) CloseWindowWithContext(ctx context.Context) error {
	if err := p.session.DeleteWindowWithContext(ctx); err != nil {
		return fmt.Errorf("failed to close active window: %w", err)
	}
	return nil
}

// WindowCount returns the number of available windows.
func (p *Page) WindowCount() (int, error) {
	return p.WindowCountWithContext(context.Background())
}

// WindowCountWithContext returns the number of available windows.
func (p *Page) WindowCountWithContext(ctx context.Context) (int, error) {
	windows, err := p.session.GetWindowsWithContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to find available windows: %w", err)
	}
	return len(windows), nil
}

// LogTypes returns all the valid log types that may be used with a LogReader.
func (p *Page) LogTypes() ([]string, error) {
	return p.LogTypesWithContext(context.Background())
}

// LogTypesWithContext returns all the valid log types that may be used with a LogReader.
func (p *Page) LogTypesWithContext(ctx context.Context) ([]string, error) {
	types, err := p.session.GetLogTypesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve log types: %w", err)
	}
	return types, nil
}

var messageMatcher = regexp.MustCompile(`^(?s:(.+))\s\(([^)]*:\w*)\)$`)

// ReadNewLogs returns new log messages of the provided log type. For example,
// page.ReadNewLogs("browser") returns browser console logs, such as JavaScript
// logs and errors. Only logs since the last call to ReadNewLogs are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadNewLogs(logType string) ([]Log, error) {
	return p.ReadNewLogsWithContext(context.Background(), logType)
}

// ReadNewLogsWithContext returns new log messages of the provided log type. For example,
// page.ReadNewLogs("browser") returns browser console logs, such as JavaScript
// logs and errors. Only logs since the last call to ReadNewLogs are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadNewLogsWithContext(ctx context.Context, logType string) ([]Log, error) {
	if p.logs == nil {
		p.logs = map[string][]Log{}
	}
	clientLogs, err := p.session.NewLogsWithContext(ctx, logType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %w", err)
	}
	var logs []Log
	for _, v := range clientLogs {
		matches := messageMatcher.FindStringSubmatch(v.Message)
		message, location := v.Message, ""
		if len(matches) > 2 {
			message, location = matches[1], matches[2]
		}
		log := Log{
			Message:  message,
			Location: location,
			Level:    v.Level,
			Time:     msToTime(v.Timestamp),
		}
		logs = append(logs, log)
		p.logs[logType] = append(p.logs[logType], log)
	}
	return logs, nil
}

// ReadAllLogs returns all log messages of the provided log type. For example,
// page.ReadAllLogs("browser") returns browser console logs, such as JavaScript logs
// and errors. All logs since the session was created are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadAllLogs(logType string) ([]Log, error) {
	return p.ReadAllLogsWithContext(context.Background(), logType)
}

// ReadAllLogsWithContext returns all log messages of the provided log type. For example,
// page.ReadAllLogs("browser") returns browser console logs, such as JavaScript logs
// and errors. All logs since the session was created are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadAllLogsWithContext(ctx context.Context, logType string) ([]Log, error) {
	if _, err := p.ReadNewLogsWithContext(ctx, logType); err != nil {
		return nil, err
	}
	ret := make([]Log, len(p.logs[logType]))
	copy(ret, p.logs[logType])
	return ret, nil
}

func msToTime(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}

// MoveMouseBy moves the mouse by the provided offset.
func (p *Page) MoveMouseBy(xOffset, yOffset int) error {
	return p.MoveMouseByWithContext(context.Background(), xOffset, yOffset)
}

// MoveMouseByWithContext moves the mouse by the provided offset.
func (p *Page) MoveMouseByWithContext(ctx context.Context, xOffset, yOffset int) error {
	if err := p.session.MoveTo(nil, session.XYOffset{
		X: xOffset,
		Y: yOffset,
	}); err != nil {
		return fmt.Errorf("failed to move mouse: %w", err)
	}
	return nil
}

// DoubleClick double-clicks the left mouse button at the current mouse position.
func (p *Page) DoubleClick() error {
	return p.DoubleClickWithContext(context.Background())
}

// DoubleClickWithContext double-clicks the left mouse button at the current mouse position.
func (p *Page) DoubleClickWithContext(ctx context.Context) error {
	if err := p.session.DoubleClickWithContext(ctx); err != nil {
		return fmt.Errorf("failed to double click: %w", err)
	}
	return nil
}

// Click performs the provided Click event using the provided Button at the
// current mouse position.
func (p *Page) Click(click event.Click, button event.Button) error {
	return p.ClickWithContext(context.Background(), click, button)
}

// ClickWithContext performs the provided Click event using the provided Button at the
// current mouse position.
func (p *Page) ClickWithContext(ctx context.Context, click event.Click, button event.Button) error {
	var err error
	switch click {
	case event.SingleClick:
		err = p.session.ClickWithContext(ctx, button)
	case event.HoldClick:
		err = p.session.ButtonDownWithContext(ctx, button)
	case event.ReleaseClick:
		err = p.session.ButtonUpWithContext(ctx, button)
	default:
		err = errors.New("invalid touch event")
	}
	if err != nil {
		return fmt.Errorf("failed to %s %s: %w", click, button, err)
	}
	return nil
}

// SetImplicitWait sets the implicit wait timeout (in ms)
func (p *Page) SetImplicitWait(timeout int) error {
	return p.SetImplicitWaitWithContext(context.Background(), timeout)
}

// SetImplicitWaitWithContext sets the implicit wait timeout (in ms)
func (p *Page) SetImplicitWaitWithContext(ctx context.Context, timeout int) error {
	return p.session.SetImplicitWaitWithContext(ctx, timeout)
}

// SetPageLoad sets the page load timeout (in ms)
func (p *Page) SetPageLoad(timeout int) error {
	return p.SetPageLoadWithContext(context.Background(), timeout)
}

// SetPageLoadWithContext sets the page load timeout (in ms)
func (p *Page) SetPageLoadWithContext(ctx context.Context, timeout int) error {
	return p.session.SetPageLoadWithContext(ctx, timeout)
}

// SetScriptTimeout sets the script timeout (in ms)
func (p *Page) SetScriptTimeout(timeout int) error {
	return p.SetScriptTimeoutWithContext(context.Background(), timeout)
}

// SetScriptTimeoutWithContext sets the script timeout (in ms)
func (p *Page) SetScriptTimeoutWithContext(ctx context.Context, timeout int) error {
	return p.session.SetScriptTimeoutWithContext(ctx, timeout)
}
