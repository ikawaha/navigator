package session

import (
	"path"
)

type Window struct {
	ID      string
	Session *Session
}

func (w *Window) Send(method, pathname string, body, result any) error {
	return w.Session.Send(method, path.Join("window", w.ID, pathname), body, result)
}

type WidthHeightRequest struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (w *Window) SetSize(width, height int) error {
	return w.Send(Post, "size", WidthHeightRequest{
		Width:  width,
		Height: height,
	}, nil)
}
