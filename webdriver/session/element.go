package session

import (
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
func (e *Element) Send(method, pathname string, body, result any) error {
	return e.Session.Send(method, path.Join("element", e.ID, pathname), body, result)
}

// GetText gets a text of the element.
func (e *Element) GetText() (string, error) {
	var text string
	if err := e.Send(Get, "text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

// GetName gets a name of the element.
func (e *Element) GetName() (string, error) {
	var name string
	if err := e.Send(Get, "name", nil, &name); err != nil {
		return "", err
	}
	return name, nil
}

// GetAttribute gets an attribute of the element.
func (e *Element) GetAttribute(attribute string) (string, error) {
	var value string
	if err := e.Send(Get, path.Join("attribute", attribute), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// GetCSS gets a CSS property of the element.
func (e *Element) GetCSS(property string) (string, error) {
	var value string
	if err := e.Send(Get, path.Join("css", property), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

// Click clicks the element.
func (e *Element) Click() error {
	return e.Send(Post, "click", nil, nil)
}

// Clear clears the element.
func (e *Element) Clear() error {
	return e.Send(Post, "clear", nil, nil)
}

// Value gets a value of the element.
func (e *Element) Value(text string) error {
	vec := strings.Split(text, "")
	req := struct {
		Value []string `json:"value"`
	}{
		Value: vec,
	}
	return e.Send(Post, "value", req, nil)
}

// IsSelected returns true if the element is selected.
func (e *Element) IsSelected() (bool, error) {
	var selected bool
	if err := e.Send(Get, "selected", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

// IsDisplayed returns true if the element is displayed.
func (e *Element) IsDisplayed() (bool, error) {
	var displayed bool
	if err := e.Send(Get, "displayed", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

// IsEnabled returns true if the element is enabled.
func (e *Element) IsEnabled() (bool, error) {
	var enabled bool
	if err := e.Send(Get, "enabled", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

// Submit submits the element.
func (e *Element) Submit() error {
	return e.Send(Post, "submit", nil, nil)
}

// IsEqualTo returns true if the elements is equal to the other.
func (e *Element) IsEqualTo(other *Element) (bool, error) {
	if other == nil {
		return false, errors.New("nil element is invalid")
	}
	var equal bool
	if err := e.Send(Get, path.Join("equals", other.ID), nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}

// GetLocation gets a location of the element.
func (e *Element) GetLocation() (x, y int, err error) {
	var location struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := e.Send(Get, "location", nil, &location); err != nil {
		return 0, 0, err
	}
	return round(location.X), round(location.Y), nil
}

// GetSize gets a size of the element.
func (e *Element) GetSize() (width, height int, err error) {
	var size struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
	if err := e.Send(Get, "size", nil, &size); err != nil {
		return 0, 0, err
	}
	return round(size.Width), round(size.Height), nil
}

func round(number float64) int {
	return int(number + 0.5)
}
