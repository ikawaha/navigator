package navigator

import (
	"errors"
	"fmt"

	"github.com/ikawaha/navigator/webdriver/session"
)

type Selectable struct {
	session   *session.Session
	selectors selectors
}

// Find finds exactly one element by CSS selector.
func (s *Selectable) Find(selector string) *Selection {
	return newSelection(s.session, s.selectors.Append(cssType, selector).Single())
}

// FindByXPath finds exactly one element by XPath selector.
func (s *Selectable) FindByXPath(selector string) *Selection {
	return newSelection(s.session, s.selectors.Append(xPathType, selector).Single())
}

// FindByLink finds exactly one anchor element by its text content.
func (s *Selectable) FindByLink(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(linkType, text).Single())
}

// FindByLabel finds exactly one element by associated label text.
func (s *Selectable) FindByLabel(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(labelType, text).Single())
}

// FindByButton finds exactly one button element with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *Selectable) FindByButton(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(buttonType, text).Single())
}

// FindByName finds exactly element with the provided name attribute.
func (s *Selectable) FindByName(name string) *Selection {
	return newSelection(s.session, s.selectors.Append(nameType, name).Single())
}

// FindByClass finds exactly one element with a given CSS class.
func (s *Selectable) FindByClass(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(classType, text).Single())
}

// FindByID finds exactly one element that has the given ID.
func (s *Selectable) FindByID(id string) *Selection {
	return newSelection(s.session, s.selectors.Append(idType, id).Single())
}

// First finds the first element by CSS selector.
func (s *Selectable) First(selector string) *Selection {
	return newSelection(s.session, s.selectors.Append(cssType, selector).At(0))
}

// FirstByXPath finds the first element by XPath selector.
func (s *Selectable) FirstByXPath(selector string) *Selection {
	return newSelection(s.session, s.selectors.Append(xPathType, selector).At(0))
}

// FirstByLink finds the first anchor element by its text content.
func (s *Selectable) FirstByLink(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(linkType, text).At(0))
}

// FirstByLabel finds the first element by associated label text.
func (s *Selectable) FirstByLabel(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(labelType, text).At(0))
}

// FirstByButton finds the first button element with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *Selectable) FirstByButton(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(buttonType, text).At(0))
}

// FirstByName finds the first element with the provided name attribute.
func (s *Selectable) FirstByName(name string) *Selection {
	return newSelection(s.session, s.selectors.Append(nameType, name).At(0))
}

// FirstByClass finds the first element with a given CSS class.
func (s *Selectable) FirstByClass(text string) *Selection {
	return newSelection(s.session, s.selectors.Append(classType, text).At(0))
}

// All finds zero or more elements by CSS selector.
func (s *Selectable) All(selector string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(cssType, selector))
}

// AllByXPath finds zero or more elements by XPath selector.
func (s *Selectable) AllByXPath(selector string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(xPathType, selector))
}

// AllByLink finds zero or more anchor elements by their text content.
func (s *Selectable) AllByLink(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(linkType, text))
}

// AllByLabel finds zero or more elements by associated label text.
func (s *Selectable) AllByLabel(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(labelType, text))
}

// AllByButton finds zero or more button elements with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *Selectable) AllByButton(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(buttonType, text))
}

// AllByName finds zero or more elements with the provided name attribute.
func (s *Selectable) AllByName(name string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(nameType, name))
}

// AllByClass finds zero or more elements with a given CSS class.
func (s *Selectable) AllByClass(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(classType, text))
}

// AllByID finds zero or more elements with a given ID.
func (s *Selectable) AllByID(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(idType, text))
}

func (s *Selectable) String() string {
	ss := make([]string, len(s.selectors))
	for i, v := range s.selectors {
		ss[i] = v.String()
	}
	return fmt.Sprintf("%+v", ss)
}

func (s *Selectable) getElementsAtLeastOne() ([]*session.Element, error) {
	elements, err := s.getElements()
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, errors.New("no elements found")
	}
	return elements, nil
}

func (s *Selectable) getElementExactlyOne() (*session.Element, error) {
	elements, err := s.getElementsAtLeastOne()
	if err != nil {
		return nil, err
	}
	if len(elements) > 1 {
		return nil, fmt.Errorf("method does not support multiple elements (%d)", len(elements))
	}
	return elements[0], nil
}

// XXX ???
func (s *Selectable) getElements() ([]*session.Element, error) {
	if len(s.selectors) == 0 {
		return nil, errors.New("empty selection")
	}
	lastElements, err := retrieveElements(s.session, s.selectors[0])
	if err != nil {
		return nil, err
	}
	for _, selector := range s.selectors[1:] {
		var elements []*session.Element
		for _, element := range lastElements {
			subElements, err := retrieveElements(element.Session, selector)
			if err != nil {
				return nil, err
			}
			elements = append(elements, subElements...)
		}
		lastElements = elements
	}
	return lastElements, nil
}

// XXX refactoring ???
func retrieveElements(client *session.Session, selector selector) ([]*session.Element, error) {
	if selector.Single {
		elements, err := client.GetElements(selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		if len(elements) == 0 {
			return nil, errors.New("element not found")
		} else if len(elements) > 1 {
			return nil, errors.New("ambiguous find")
		}
		return []*session.Element{elements[0]}, nil
	}

	if selector.Indexed && selector.Index > 0 {
		elements, err := client.GetElements(selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		if selector.Index >= len(elements) {
			return nil, errors.New("element index out of range")
		}

		return []*session.Element{elements[selector.Index]}, nil
	}

	if selector.Indexed && selector.Index == 0 {
		element, err := client.GetElement(selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		return []*session.Element{element}, nil
	}

	elements, err := client.GetElements(selector.SessionSelector())
	if err != nil {
		return nil, err
	}
	return elements, nil
}
