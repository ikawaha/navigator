package navigator

import (
	"strings"
)

type selectors []selector

func (s selectors) canMergeType(selectorType selectorType) bool {
	if len(s) == 0 {
		return false
	}
	last := s[len(s)-1]
	bothCSS := selectorType == cssType && last.Type == cssType
	return bothCSS && !last.Indexed && !last.Single
}

// XXX to private ?
func (s selectors) Append(selectorType selectorType, value string) selectors {
	selector := selector{
		Type:  selectorType,
		Value: value,
	}
	if s.canMergeType(selectorType) {
		lastIndex := len(s) - 1
		selector.Value = s[lastIndex].Value + " " + selector.Value
		return s[:lastIndex].append(selector)
	}
	return s.append(selector)
}

func (s selectors) Single() selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}
	selector := s[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return s[:lastIndex].append(selector)
}

func (s selectors) At(index int) selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}
	selector := s[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return s[:lastIndex].append(selector)
}

func (s selectors) String() string {
	var tags []string
	for _, selector := range s {
		tags = append(tags, selector.String())
	}
	return strings.Join(tags, " | ")
}

func (s selectors) append(selector selector) selectors {
	selectorsCopy := append(selectors(nil), s...)
	return append(selectorsCopy, selector)
}
