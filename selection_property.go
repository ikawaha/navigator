package navigator

import (
	"context"
	"fmt"

	"github.com/ikawaha/navigator/webdriver/session"
)

// Text returns the entirety of the text content for exactly one element.
func (s *Selection) Text() (string, error) {
	return s.TextWithContext(context.Background())
}

// TextWithContext returns the entirety of the text content for exactly one element.
func (s *Selection) TextWithContext(ctx context.Context) (string, error) {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	text, err := selectedElement.GetTextWithContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for %s: %w", s, err)
	}
	return text, nil
}

// Active returns true if the single element that the selection refers to is active.
func (s *Selection) Active() (bool, error) {
	return s.ActiveWithContext(context.Background())
}

// ActiveWithContext returns true if the single element that the selection refers to is active.
func (s *Selection) ActiveWithContext(ctx context.Context) (bool, error) {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	activeElement, err := s.session.GetActiveElementWithContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve active element: %w", err)
	}
	equal, err := selectedElement.IsEqualToWithContext(ctx, activeElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare selection to active element: %w", err)
	}
	return equal, nil
}

type propertyMethod func(element *session.Element, ctx context.Context, property string) (string, error)

func (s *Selection) hasProperty(ctx context.Context, method propertyMethod, property, name string) (string, error) {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	value, err := method(selectedElement, ctx, property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s value for %s: %w", name, s, err)
	}
	return value, nil
}

// Attribute returns an attribute value for exactly one element.
func (s *Selection) Attribute(attribute string) (string, error) {
	return s.AttributeWithContext(context.Background(), attribute)
}

// AttributeWithContext returns an attribute value for exactly one element.
func (s *Selection) AttributeWithContext(ctx context.Context, attribute string) (string, error) {
	return s.hasProperty(ctx, (*session.Element).GetAttributeWithContext, attribute, "attribute")
}

// CSS returns a CSS style property value for exactly one element.
func (s *Selection) CSS(property string) (string, error) {
	return s.CSSWithContext(context.Background(), property)
}

// CSSWithContext returns a CSS style property value for exactly one element.
func (s *Selection) CSSWithContext(ctx context.Context, property string) (string, error) {
	return s.hasProperty(ctx, (*session.Element).GetCSSWithContext, property, "CSS property")
}

type stateMethod func(element *session.Element, ctx context.Context) (bool, error)

func (s *Selection) hasState(ctx context.Context, method stateMethod, name string) (bool, error) {
	elements, err := s.getElementsAtLeastOne(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to select elements from %s: %w", s, err)
	}
	for _, selectedElement := range elements {
		pass, err := method(selectedElement, ctx)
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
func (s *Selection) Selected() (bool, error) {
	return s.SelectedWithContext(context.Background())
}

// SelectedWithContext returns true if all the elements that the selection refers to are selected.
func (s *Selection) SelectedWithContext(ctx context.Context) (bool, error) {
	return s.hasState(ctx, (*session.Element).IsSelectedWithContext, "selected")
}

// Visible returns true if all the elements that the selection refers to are visible.
func (s *Selection) Visible() (bool, error) {
	return s.VisibleWithContext(context.Background())
}

// VisibleWithContext returns true if all the elements that the selection refers to are visible.
func (s *Selection) VisibleWithContext(ctx context.Context) (bool, error) {
	return s.hasState(ctx, (*session.Element).IsDisplayedWithContext, "visible")
}

// Enabled returns true if all the elements that the selection refers to are enabled.
func (s *Selection) Enabled() (bool, error) {
	return s.EnabledWithContext(context.Background())
}

// EnabledWithContext returns true if all the elements that the selection refers to are enabled.
func (s *Selection) EnabledWithContext(ctx context.Context) (bool, error) {
	return s.hasState(ctx, (*session.Element).IsEnabledWithContext, "enabled")
}
