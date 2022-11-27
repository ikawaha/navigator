package navigator

import (
	"fmt"

	"github.com/ikawaha/navi/webdriver/session"
)

type SelectorType string

const (
	cssType             SelectorType = "CSS: %s"
	xPathType           SelectorType = "XPath: %s"
	linkType            SelectorType = `Link: "%s"`
	labelType           SelectorType = `Label: "%s"`
	buttonType          SelectorType = `Button: "%s"`
	nameType            SelectorType = `Name: "%s"`
	accessibilityIDType SelectorType = "Accessibility ID: %s"
	androidAutType      SelectorType = "Android UIAut.: %s"
	iosAutType          SelectorType = "iOS UIAut.: %s"
	classType           SelectorType = "Class: %s"
	idType              SelectorType = "ID: %s"
)

func (t SelectorType) format(value string) string {
	return fmt.Sprintf(string(t), value)
}

type Selector struct {
	Type    SelectorType
	Value   string
	Index   int
	Indexed bool
	Single  bool
}

func (s Selector) String() string {
	var suffix string
	if s.Single {
		suffix = " [single]"
	} else if s.Indexed {
		suffix = fmt.Sprintf(" [%d]", s.Index)
	}
	return s.Type.format(s.Value) + suffix
}

func (s Selector) SessionSelector() session.Selector {
	return session.Selector{
		Using: s.selectorType(),
		Value: s.value(),
	}
}

func (s Selector) selectorType() string {
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

func (s Selector) value() string {
	switch s.Type {
	case labelType:
		return fmt.Sprintf(labelXPath, s.Value)
	case buttonType:
		return fmt.Sprintf(buttonXPath, s.Value)
	}
	return s.Value
}
