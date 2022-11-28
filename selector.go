package navigator

import (
	"fmt"

	"github.com/ikawaha/navi/webdriver/session"
)

type selectorType string

const (
	cssType             selectorType = "CSS: %s"
	xPathType           selectorType = "XPath: %s"
	linkType            selectorType = `Link: "%s"`
	labelType           selectorType = `Label: "%s"`
	buttonType          selectorType = `Button: "%s"`
	nameType            selectorType = `Name: "%s"`
	accessibilityIDType selectorType = "Accessibility ID: %s"
	androidAutType      selectorType = "Android UIAut.: %s"
	iosAutType          selectorType = "iOS UIAut.: %s"
	classType           selectorType = "Class: %s"
	idType              selectorType = "ID: %s"
)

func (t selectorType) format(value string) string {
	return fmt.Sprintf(string(t), value)
}

type selector struct {
	Type    selectorType
	Value   string
	Index   int
	Indexed bool
	Single  bool
}

func (s selector) String() string {
	var suffix string
	if s.Single {
		suffix = " [single]"
	} else if s.Indexed {
		suffix = fmt.Sprintf(" [%d]", s.Index)
	}
	return s.Type.format(s.Value) + suffix
}

func (s selector) SessionSelector() session.Selector {
	return session.Selector{
		Using: s.selectorType(),
		Value: s.value(),
	}
}

func (s selector) selectorType() string {
	switch s.Type {
	case cssType:
		return "css selector"
	case classType:
		return "class name"
	case idType:
		return "id"
	case linkType:
		return "link text"
	case nameType:
		return "name"
	case accessibilityIDType:
		return "accessibility id"
	case androidAutType:
		return "-android uiautomator"
	case iosAutType:
		return "-ios uiautomation"
	}
	return "xpath"
}

const (
	labelXPath  = `//input[@id=(//label[normalize-space()="%s"]/@for)] | //label[normalize-space()="%[1]s"]/input`
	buttonXPath = `//input[@type="submit" or @type="button"][normalize-space(@value)="%s"] | //button[normalize-space()="%[1]s"]`
)

func (s selector) value() string {
	switch s.Type {
	case labelType:
		return fmt.Sprintf(labelXPath, s.Value)
	case buttonType:
		return fmt.Sprintf(buttonXPath, s.Value)
	}
	return s.Value
}
