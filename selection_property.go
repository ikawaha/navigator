package navigator

import (
	"fmt"

	"github.com/ikawaha/navigator/webdriver/session"
)

// Text returns the entirety of the text content for exactly one element.
func (s *Selection) Text() (string, error) {
	selectedElement, err := s.getElementExactlyOne()
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	text, err := selectedElement.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for %s: %w", s, err)
	}
	return text, nil
}

// Active returns true if the single element that the selection refers to is active.
func (s *Selection) Active() (bool, error) {
	selectedElement, err := s.getElementExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	activeElement, err := s.session.GetActiveElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve active element: %w", err)
	}
	equal, err := selectedElement.IsEqualTo(activeElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare selection to active element: %w", err)
	}
	return equal, nil
}

type propertyMethod func(element *session.Element, property string) (string, error)

func (s *Selection) hasProperty(method propertyMethod, property, name string) (string, error) {
	selectedElement, err := s.getElementExactlyOne()
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	value, err := method(selectedElement, property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s value for %s: %w", name, s, err)
	}
	return value, nil
}

// Attribute returns an attribute value for exactly one element.
// XXX refactoring ???
func (s *Selection) Attribute(attribute string) (string, error) {
	return s.hasProperty((*session.Element).GetAttribute, attribute, "attribute")
}

// CSS returns a CSS style property value for exactly one element.
// XXX refactoring ???
func (s *Selection) CSS(property string) (string, error) {
	return s.hasProperty((*session.Element).GetCSS, property, "CSS property")
}

type stateMethod func(element *session.Element) (bool, error)

func (s *Selection) hasState(method stateMethod, name string) (bool, error) {
	elements, err := s.getElementsAtLeastOne()
	if err != nil {
		return false, fmt.Errorf("failed to select elements from %s: %w", s, err)
	}
	for _, selectedElement := range elements {
		pass, err := method(selectedElement)
		if err != nil {
			return false, fmt.Errorf("failed to determine whether %s is %s: %w", s, name, err)
		}
		if !pass {
			return false, nil
		}
	}
	return true, nil
}

// Selected returns true if all the elements that the selection refers to are selected.
// XXX refactoring ???
func (s *Selection) Selected() (bool, error) {
	return s.hasState((*session.Element).IsSelected, "selected")
}

// Visible returns true if all the elements that the selection refers to are visible.
// XXX refactoring ???
func (s *Selection) Visible() (bool, error) {
	return s.hasState((*session.Element).IsDisplayed, "visible")
}

// Enabled returns true if all the elements that the selection refers to are enabled.
// XXX refactoring ???
func (s *Selection) Enabled() (bool, error) {
	return s.hasState((*session.Element).IsEnabled, "enabled")
}
