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

// GetElement gets an element by the selector.
func (e *Element) GetElement(ctx context.Context, selector Selector) (*Element, error) {
	var result elementResult
	if err := e.Send(ctx, Post, "element", selector, &result); err != nil {
		return nil, err
	}
	return &Element{ID: result.ID(), Session: e.Session}, nil
}

// GetElements gets elements by the selector.
func (e *Element) GetElements(ctx context.Context, selector Selector) ([]*Element, error) {
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

// GetText gets a text of the element.
func (e *Element) GetText(ctx context.Context) (string, error) {
	var text string
	if err := e.Send(ctx, Get, "text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

// GetName gets a name of the element.
func (e *Element) GetName(ctx context.Context) (string, error) {
	var name string
	if err := e.Send(ctx, Get, "name", nil, &name); err != nil {
		return "", err
	}
	return name, nil
}

// GetAttribute gets an attribute of the element.
func (e *Element) GetAttribute(ctx context.Context, attribute string) (string, error) {
	var value string
	if err := e.Send(ctx, Get, path.Join("attribute", attribute), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// GetCSS gets a CSS property of the element.
func (e *Element) GetCSS(ctx context.Context, property string) (string, error) {
	var value string
	if err := e.Send(ctx, Get, path.Join("css", property), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// Click clicks the element.
func (e *Element) Click(ctx context.Context) error {
	return e.Send(ctx, Post, "click", nil, nil)
}

// Clear clears the element.
func (e *Element) Clear(ctx context.Context) error {
	return e.Send(ctx, Post, "clear", nil, nil)
}

// Value sends keys corresponding to the text.
func (e *Element) Value(ctx context.Context, text string) error {
	vec := strings.Split(text, "")
	req := struct {
		Value []string `json:"value"`
	}{
		Value: vec,
	}
	return e.Send(ctx, Post, "value", req, nil)
}

// IsSelected returns true if the element is selected.
func (e *Element) IsSelected(ctx context.Context) (bool, error) {
	var selected bool
	if err := e.Send(ctx, Get, "selected", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

// IsDisplayed returns true if the element is displayed.
func (e *Element) IsDisplayed(ctx context.Context) (bool, error) {
	var displayed bool
	if err := e.Send(ctx, Get, "displayed", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

// IsEnabled returns true if the element is enabled.
func (e *Element) IsEnabled(ctx context.Context) (bool, error) {
	var enabled bool
	if err := e.Send(ctx, Get, "enabled", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

// Submit submits the element.
func (e *Element) Submit(ctx context.Context) error {
	return e.Send(ctx, Post, "submit", nil, nil)
}

// IsEqualTo returns true if the elements is equal to the other.
func (e *Element) IsEqualTo(ctx context.Context, other *Element) (bool, error) {
	if other == nil {
		return false, errors.New("nil element is invalid")
	}
	var equal bool
	if err := e.Send(ctx, Get, path.Join("equals", other.ID), nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}

// GetLocation gets a location of the element.
func (e *Element) GetLocation(ctx context.Context) (x, y int, err error) {
	var location struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := e.Send(ctx, Get, "location", nil, &location); err != nil {
		return 0, 0, err
	}
	return round(location.X), round(location.Y), nil
}

// GetSize gets a size of the element.
func (e *Element) GetSize(ctx context.Context) (width, height int, err error) {
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
