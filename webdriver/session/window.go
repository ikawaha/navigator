package session

import (
	"context"
	"path"
)

// Window represents the window handler of the browser.
type Window struct {
	ID      string
	Session *Session
}

// Send sends the message to the window.
func (w *Window) Send(method, pathname string, body, result any) error {
	return w.SendWithContext(context.Background(), method, pathname, body, result)
}

// SendWithContext sends the message to the window.
func (w *Window) SendWithContext(ctx context.Context, method, pathname string, body, result any) error {
	return w.Session.Send(ctx, method, path.Join("window", w.ID, pathname), body, result)
}

type widthHeightRequest struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// SetSize sets the size of the window of the browser.
func (w *Window) SetSize(width, height int) error {
	return w.Send(Post, "size", widthHeightRequest{
		Width:  width,
		Height: height,
	}, nil)
}
