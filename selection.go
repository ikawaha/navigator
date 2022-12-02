package navigator

import (
	"fmt"

	"github.com/ikawaha/navigator/webdriver/session"
)

// Selection instances refer to a selection of elements.
// All Selection methods are also MultiSelection methods.
//
// Methods that take selectors apply their selectors to each element in the
// selection they are called on. If the selection they are called on refers to multiple
// elements, the resulting selection will refer to at least that many elements.
//
// Examples:
//
//	selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
//
// Checks the first checkbox in the third row of the only table.
//
//	selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
//
// Checks all checkboxes in the first-and-only cell of each row in the only table.
type Selection struct {
	Selectable
}

func newSelection(session *session.Session, selectors selectors) *Selection {
	return &Selection{
		Selectable: Selectable{
			session:   session,
			selectors: selectors,
		},
	}
}

// String returns a string representation of the selection, ex.
//
//	selection 'CSS: .some-class | XPath: //table [3] | Link "click me" [single]'
func (s *Selection) String() string {
	return fmt.Sprintf("selection '%s'", s.selectors)
}

// Elements returns a []*webdriver.Element that can be used to send direct commands
// to WebDriver elements. See: https://code.google.com/p/selenium/wiki/JsonWireProtocol
func (s *Selection) Elements() ([]*session.Element, error) {
	elements, err := s.getElements()
	if err != nil {
		return nil, err
	}
	var apiElements []*session.Element
	for _, selectedElement := range elements {
		apiElements = append(apiElements, selectedElement)
	}
	return apiElements, nil
}

// Count returns the number of elements that the selection refers to.
func (s *Selection) Count() (int, error) {
	elements, err := s.getElements()
	if err != nil {
		return 0, fmt.Errorf("failed to select elements from %s: %w", s, err)
	}

	return len(elements), nil
}

// EqualsElement returns whether two selections of exactly
// one element refer to the same element.
func (s *Selection) EqualsElement(other any) (bool, error) {
	otherSelection, ok := other.(*Selection)
	if !ok {
		multiSelection, ok := other.(*MultiSelection)
		if !ok {
			return false, fmt.Errorf("must be *Selection or *MultiSelection")
		}
		otherSelection = &multiSelection.Selection
	}

	selectedElement, err := s.getElementExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %w", s, err)
	}

	otherElement, err := otherSelection.getElementExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %w", other, err)
	}

	equal, err := selectedElement.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare %s to %s: %w", s, other, err)
	}

	return equal, nil
}

// MouseToElement moves the mouse over exactly one element in the selection.
func (s *Selection) MouseToElement() error {
	selectedElement, err := s.getElementExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	if err := s.session.MoveTo(selectedElement, nil); err != nil {
		return fmt.Errorf("failed to move mouse to element for %s: %w", s, err)
	}

	return nil
}
