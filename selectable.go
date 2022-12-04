package navigator

import (
	"context"
	"errors"
	"fmt"

	"github.com/ikawaha/navigator/webdriver/session"
)

// Selectable represents a set of selectable elements.
type Selectable struct {
	session   *session.Session
	selectors selectors
}

// Find finds exactly one element by CSS selector.
func (s *Selectable) Find(css string) *Selection {
	return newSelection(s.session, s.selectors.Append(cssType, css).Single())
}

// FindByXPath finds exactly one element by XPath selector.
func (s *Selectable) FindByXPath(xpath string) *Selection {
	return newSelection(s.session, s.selectors.Append(xPathType, xpath).Single())
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
func (s *Selectable) FindByClass(class string) *Selection {
	return newSelection(s.session, s.selectors.Append(classType, class).Single())
}

// FindByID finds exactly one element that has the given ID.
func (s *Selectable) FindByID(id string) *Selection {
	return newSelection(s.session, s.selectors.Append(idType, id).Single())
}

// First finds the first element by CSS selector.
func (s *Selectable) First(css string) *Selection {
	return newSelection(s.session, s.selectors.Append(cssType, css).At(0))
}

// FirstByXPath finds the first element by XPath selector.
func (s *Selectable) FirstByXPath(xpath string) *Selection {
	return newSelection(s.session, s.selectors.Append(xPathType, xpath).At(0))
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
func (s *Selectable) FirstByClass(class string) *Selection {
	return newSelection(s.session, s.selectors.Append(classType, class).At(0))
}

// All finds zero or more elements by CSS selector.
func (s *Selectable) All(css string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(cssType, css))
}

// AllByXPath finds zero or more elements by XPath selector.
func (s *Selectable) AllByXPath(xpath string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(xPathType, xpath))
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
func (s *Selectable) AllByClass(class string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(classType, class))
}

// AllByID finds zero or more elements with a given ID.
func (s *Selectable) AllByID(id string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.Append(idType, id))
}

func (s *Selectable) String() string {
	ss := make([]string, len(s.selectors))
	for i, v := range s.selectors {
		ss[i] = v.String()
	}
	return fmt.Sprintf("%+v", ss)
}

func (s *Selectable) getElementsAtLeastOne(ctx context.Context) ([]*session.Element, error) {
	elements, err := s.getElements(ctx)
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, errors.New("no elements found")
	}
	return elements, nil
}

func (s *Selectable) getElementExactlyOne(ctx context.Context) (*session.Element, error) {
	elements, err := s.getElementsAtLeastOne(ctx)
	if err != nil {
		return nil, err
	}
	if len(elements) > 1 {
		return nil, fmt.Errorf("method does not support multiple elements (%d)", len(elements))
	}
	return elements[0], nil
}

func (s *Selectable) getElements(ctx context.Context) ([]*session.Element, error) {
	if len(s.selectors) == 0 {
		return nil, errors.New("empty selection")
	}
	ret := []*session.Element{{Session: s.session}} // initial dummy element
	for _, sl := range s.selectors {
		var next []*session.Element
		for _, el := range ret {
			els, err := retrieveElements(ctx, el, sl)
			if err != nil {
				return nil, err
			}
			next = append(next, els...)
		}
		ret = next
	}
	return ret, nil
}

func retrieveElements(ctx context.Context, element *session.Element, selector selector) ([]*session.Element, error) {
	switch {
	case selector.Single:
		els, err := element.GetElements(ctx, selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		if len(els) == 0 {
			return nil, errors.New("element not found")
		} else if len(els) > 1 {
			return nil, errors.New("ambiguous find")
		}
		return els[:1], nil
	case selector.Indexed && selector.Index == 0:
		el, err := element.GetElement(ctx, selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		return []*session.Element{el}, nil
	case selector.Indexed && selector.Index > 0:
		els, err := element.GetElements(ctx, selector.SessionSelector())
		if err != nil {
			return nil, err
		}
		if selector.Index < 0 || selector.Index >= len(els) {
			return nil, errors.New("element index out of range")
		}
		return []*session.Element{els[selector.Index]}, nil
	}
	return element.GetElements(ctx, selector.SessionSelector())
}
