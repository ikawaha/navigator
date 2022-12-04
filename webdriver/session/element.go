package session

import (
	"context"
	"errors"
	"path"
	"strings"
)

// Element represents a web element.
type Element struct {
	ID      string
	Session *Session
}

// Send sends a message to the web driver service.
func (e *Element) Send(ctx context.Context, method, pathname string, body, result any) error {
	if e.ID != "" {
		pathname = path.Join("element", e.ID, pathname)
	}
	return e.Session.Send(ctx, method, pathname, body, result)
}

// GetElementWithContext gets an element by the selector.
func (e *Element) GetElementWithContext(ctx context.Context, selector Selector) (*Element, error) {
	var result elementResult
	if err := e.Send(ctx, Post, "element", selector, &result); err != nil {
		return nil, err
	}
	return &Element{ID: result.ID(), Session: e.Session}, nil
}

// GetElementsWithContext gets elements by the selector.
func (e *Element) GetElementsWithContext(ctx context.Context, selector Selector) ([]*Element, error) {
	var results []elementResult
	if err := e.Send(ctx, Post, "elements", selector, &results); err != nil {
		return nil, err
	}
	elements := make([]*Element, len(results))
	for i, result := range results {
		elements[i] = &Element{ID: result.ID(), Session: e.Session}
	}
	return elements, nil
}

// GetTextWithContext gets a text of the element.
func (e *Element) GetTextWithContext(ctx context.Context) (string, error) {
	var text string
	if err := e.Send(ctx, Get, "text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

// GetNameWithContext gets a name of the element.
func (e *Element) GetNameWithContext(ctx context.Context) (string, error) {
	var name string
	if err := e.Send(ctx, Get, "name", nil, &name); err != nil {
		return "", err
	}
	return name, nil
}

// GetAttributeWithContext gets an attribute of the element.
func (e *Element) GetAttributeWithContext(ctx context.Context, attribute string) (string, error) {
	var value string
	if err := e.Send(ctx, Get, path.Join("attribute", attribute), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// GetCSSWithContext gets a CSS property of the element.
func (e *Element) GetCSSWithContext(ctx context.Context, property string) (string, error) {
	var value string
	if err := e.Send(ctx, Get, path.Join("css", property), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// ClickWithContext clicks the element.
func (e *Element) ClickWithContext(ctx context.Context) error {
	return e.Send(ctx, Post, "click", nil, nil)
}

// ClearWithContext clears the element.
func (e *Element) ClearWithContext(ctx context.Context) error {
	return e.Send(ctx, Post, "clear", nil, nil)
}

// ValueWithContext sends keys corresponding to the text.
func (e *Element) ValueWithContext(ctx context.Context, text string) error {
	vec := strings.Split(text, "")
	req := struct {
		Value []string `json:"value"`
	}{
		Value: vec,
	}
	return e.Send(ctx, Post, "value", req, nil)
}

// IsSelectedWithContext returns true if the element is selected.
func (e *Element) IsSelectedWithContext(ctx context.Context) (bool, error) {
	var selected bool
	if err := e.Send(ctx, Get, "selected", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

// IsDisplayedWithContext returns true if the element is displayed.
func (e *Element) IsDisplayedWithContext(ctx context.Context) (bool, error) {
	var displayed bool
	if err := e.Send(ctx, Get, "displayed", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

// IsEnabledWithContext returns true if the element is enabled.
func (e *Element) IsEnabledWithContext(ctx context.Context) (bool, error) {
	var enabled bool
	if err := e.Send(ctx, Get, "enabled", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

// SubmitWithContext submits the element.
func (e *Element) SubmitWithContext(ctx context.Context) error {
	return e.Send(ctx, Post, "submit", nil, nil)
}

// IsEqualToWithContext returns true if the elements is equal to the other.
func (e *Element) IsEqualToWithContext(ctx context.Context, other *Element) (bool, error) {
	if other == nil {
		return false, errors.New("nil element is invalid")
	}
	var equal bool
	if err := e.Send(ctx, Get, path.Join("equals", other.ID), nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}

// GetLocationWithContext gets a location of the element.
func (e *Element) GetLocationWithContext(ctx context.Context) (x, y int, err error) {
	var location struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := e.Send(ctx, Get, "location", nil, &location); err != nil {
		return 0, 0, err
	}
	return round(location.X), round(location.Y), nil
}

// GetSizeWithContext gets a size of the element.
func (e *Element) GetSizeWithContext(ctx context.Context) (width, height int, err error) {
	var size struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
	if err := e.Send(ctx, Get, "size", nil, &size); err != nil {
		return 0, 0, err
	}
	return round(size.Width), round(size.Height), nil
}

func round(number float64) int {
	return int(number + 0.5)
}
